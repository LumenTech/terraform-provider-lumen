package resource_bare_metal_server

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-lumen/lumen/helper"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-lumen/lumen/client"
	"terraform-provider-lumen/lumen/client/model/bare_metal"
	"terraform-provider-lumen/lumen/validation"
)

var createTimeout = schema.DefaultTimeout(90 * time.Minute)

func createContext(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	bmClient := i.(*client.Clients).BareMetal

	provisionRequest := bare_metal.ServerProvisionRequest{
		LocationID:    data.Get("location_id").(string),
		Configuration: data.Get("configuration_name").(string),
		OSImage:       data.Get("os_image_name").(string),
		Name:          data.Get("name").(string),
		Credentials: bare_metal.Credentials{
			Username:  data.Get("username").(string),
			Password:  data.Get("password").(string),
			PublicKey: data.Get("ssh_public_key").(string),
		},
		AssignIPV6Address: data.Get("assign_ipv6_address").(bool),
		Hyperthreading:    data.Get("enable_hyperthreading").(bool),
	}

	attachNetworks := convertDataToAttachedNetworks(data.Get("attach_networks").([]interface{}))
	if len(attachNetworks) != 0 {
		if err := validation.ValidateBareMetalNetworkIds(attachNetworks); err != nil {
			return diag.FromErr(err)
		}
		provisionRequest.NetworkID = attachNetworks[0].NetworkID
		provisionRequest.AssignIPV6Address = attachNetworks[0].AssignIPV6
	} else {
		provisionRequest.NetworkRequest = &bare_metal.NetworkProvisionRequest{
			Name:           data.Get("network_name").(string),
			LocationID:     provisionRequest.LocationID,
			NetworkSizeID:  data.Get("network_size_id").(string),
			NetworkType:    data.Get("network_type").(string),
			VRF:            data.Get("vrf").(string),
			VRFDescription: data.Get("vrf_description").(string),
		}
	}

	server, diagnostics := createServerAndWaitForCompletion(ctx, bmClient, provisionRequest)
	if diagnostics.HasError() {
		return diagnostics
	}
	data.SetId(server.ID)
	populateServerSchema(data, *server)

	if len(attachNetworks) > 1 {
		refreshServer, networkDiagnostics := attachNetworksAndWaitForCompletion(ctx, bmClient, server.ID, attachNetworks[1:])
		if refreshServer != nil {
			populateServerSchema(data, *refreshServer)
		}
		return append(diagnostics, networkDiagnostics...)
	}

	return diagnostics
}

var pendingServerStatuses = []string{"provisioning", "network_provisioned", "allocated", "configured", "unknown"}
var targetServerStatuses = []string{"provisioned"}
var possibleServerStatus = append(append(pendingServerStatuses, targetServerStatuses...), "failed", "error")

func createServerAndWaitForCompletion(ctx context.Context, bmClient *client.BareMetalClient, provisionRequest bare_metal.ServerProvisionRequest) (*bare_metal.Server, diag.Diagnostics) {
	server, diagnostics := bmClient.ProvisionServer(provisionRequest)
	if diagnostics.HasError() {
		return nil, diagnostics
	}

	stateChangeConf := &resource.StateChangeConf{
		Pending: pendingServerStatuses,
		Target:  targetServerStatuses,
		Refresh: func() (interface{}, string, error) {
			s, getDiagnostics := bmClient.GetServer(server.ID)
			if err := helper.ExtractDiagnosticErrorIfPresent(getDiagnostics); err != nil {
				return nil, "", fmt.Errorf(err.Summary)
			}

			found := false
			for _, status := range possibleServerStatus {
				// This is to avoid failure if a new status is added due to the logic
				// of the polling mechanism if the status isn't in the pending or target list
				// it is considered a failure and errors out immediate.
				if status == strings.ToLower(s.Status) {
					found = true
					break
				}
			}

			if !found {
				s.Status = "unknown"
			}
			return *s, strings.ToLower(s.Status), nil
		},
		Timeout:      90 * time.Minute,
		Delay:        4 * time.Minute,
		PollInterval: 30 * time.Second,
	}
	refreshResult, err := stateChangeConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, append(diagnostics, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("error waiting for server (%s) to be provisioned: (%s)", server.ID, err.Error()),
		})
	}

	serv := refreshResult.(bare_metal.Server)
	return &serv, diagnostics
}

func attachNetworksAndWaitForCompletion(ctx context.Context, bmClient *client.BareMetalClient, serverId string, networks []bare_metal.AttachNetwork) (*bare_metal.Server, diag.Diagnostics) {
	var networkDiagnostics diag.Diagnostics
	var addedNetworkIds []string
	for _, network := range networks {
		addNetworkRequest := bare_metal.AddNetworkRequest{
			NetworkId:         network.NetworkID,
			AssignIPV6Address: network.AssignIPV6,
		}
		_, attachDiagnostics := bmClient.AttachNetwork(serverId, addNetworkRequest)
		if err := helper.ExtractDiagnosticErrorIfPresent(attachDiagnostics); err != nil {
			networkDiagnostics = append(networkDiagnostics, diag.Diagnostic{
				Severity:      diag.Warning,
				Summary:       fmt.Sprintf("Error attaching network %s", network.NetworkID),
				Detail:        fmt.Sprintf("Network %s errored on attachment reason - %s", network.NetworkID, err.Summary),
				AttributePath: cty.GetAttrPath("network_ids"),
			})
		} else {
			addedNetworkIds = append(addedNetworkIds, network.NetworkID)
		}
	}

	var refreshServer *bare_metal.Server
	for _, networkId := range addedNetworkIds {
		stateChangeConf := &resource.StateChangeConf{
			Pending: []string{"provisioning"},
			Target:  []string{"provisioned"},
			Refresh: func() (interface{}, string, error) {
				s, getDiagnostics := bmClient.GetServer(serverId)
				if err := helper.ExtractDiagnosticErrorIfPresent(getDiagnostics); err != nil {
					return nil, "", fmt.Errorf(err.Summary)
				}

				refreshServer = s
				status := "provisioning"
				for _, n := range s.Networks {
					if n.NetworkID == networkId {
						status = n.Status
						break
					}
				}

				return s, status, nil
			},
			Timeout:      10 * time.Minute,
			Delay:        30 * time.Second,
			PollInterval: 30 * time.Second,
		}

		_, waitError := stateChangeConf.WaitForStateContext(ctx)
		if waitError != nil {
			networkDiagnostics = append(networkDiagnostics, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Error polling server %s for network %s", serverId, networkId),
				Detail:   waitError.Error(),
			})
		}
	}

	networkDiagnostics = append(networkDiagnostics, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Server Configuration Updates Required",
		Detail: `Automation configures Lumen networking infrastructure only, but not server configuration.
Adding a network to an existing server will require you to make configuration changes on your server.`,
		AttributePath: cty.GetAttrPath("network_ids"),
	})
	return refreshServer, networkDiagnostics
}
