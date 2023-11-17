package resource_bare_metal_server

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
	"terraform-provider-lumen/lumen/client"
	"terraform-provider-lumen/lumen/client/model/bare_metal"
	"terraform-provider-lumen/lumen/validation"
	"time"
)

var createTimeout = schema.DefaultTimeout(90 * time.Minute)

func createContext(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	bmClient := i.(*client.Clients).BareMetal

	provisionRequest := bare_metal.ServerProvisionRequest{
		Name:          data.Get("name").(string),
		LocationID:    data.Get("location_id").(string),
		Configuration: data.Get("configuration_name").(string),
		OSImage:       data.Get("os_image_name").(string),
		Credentials: bare_metal.Credentials{
			Username:  data.Get("username").(string),
			Password:  data.Get("password").(string),
			PublicKey: data.Get("ssh_public_key").(string),
		},
	}

	networkIds := convertListOfInterfaceToListOfString(data.Get("network_ids").([]interface{}))
	if len(networkIds) != 0 {
		if err := validation.ValidateBareMetalNetworkIds(networkIds); err != nil {
			return diag.FromErr(err)
		}
		provisionRequest.NetworkID = networkIds[0]
	} else {
		provisionRequest.NetworkRequest = &bare_metal.NetworkProvisionRequest{
			Name:          data.Get("network_name").(string),
			LocationID:    provisionRequest.LocationID,
			NetworkSizeID: data.Get("network_size_id").(string),
			NetworkType:   "INTERNET",
		}
	}

	server, errorDiagnostics := createServerAndWaitForCompletion(ctx, bmClient, provisionRequest)
	if errorDiagnostics != nil {
		return errorDiagnostics
	}
	data.SetId(server.ID)
	populateServerSchema(data, *server)

	if len(networkIds) > 1 {
		refreshServer, networkDiagnostics := attachNetworksAndWaitForCompletion(ctx, bmClient, server.ID, networkIds[1:])
		if refreshServer != nil {
			populateServerSchema(data, *refreshServer)
		}
		return networkDiagnostics
	}

	return nil
}

var pendingServerStatuses = []string{"provisioning", "network_provisioned", "allocated", "configured", "unknown"}
var targetServerStatuses = []string{"provisioned"}
var possibleServerStatus = append(append(pendingServerStatuses, targetServerStatuses...), "failed", "error")

func createServerAndWaitForCompletion(ctx context.Context, bmClient *client.BareMetalClient, provisionRequest bare_metal.ServerProvisionRequest) (*bare_metal.Server, diag.Diagnostics) {
	server, err := bmClient.ProvisionServer(provisionRequest)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	stateChangeConf := &resource.StateChangeConf{
		Pending: pendingServerStatuses,
		Target:  targetServerStatuses,
		Refresh: func() (interface{}, string, error) {
			s, err := bmClient.GetServer(server.ID)
			if err != nil {
				return nil, "", err
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
		return nil, diag.FromErr(err)
	}

	serv := refreshResult.(bare_metal.Server)
	return &serv, nil
}

func attachNetworksAndWaitForCompletion(ctx context.Context, bmClient *client.BareMetalClient, serverId string, networkIds []string) (*bare_metal.Server, diag.Diagnostics) {
	var networkDiagnostics diag.Diagnostics
	var addedNetworkIds []string
	for _, networkId := range networkIds {
		_, e := bmClient.AttachNetwork(serverId, networkId)
		if e != nil {
			networkDiagnostics = append(networkDiagnostics, diag.Diagnostic{
				Severity:      diag.Warning,
				Summary:       fmt.Sprintf("Error attaching network %s", networkId),
				Detail:        fmt.Sprintf("Network %s errored on attachment reason - %s", networkId, e),
				AttributePath: cty.GetAttrPath("network_ids"),
			})
		} else {
			addedNetworkIds = append(addedNetworkIds, networkId)
		}
	}

	var refreshServer *bare_metal.Server
	for _, networkId := range addedNetworkIds {
		stateChangeConf := &resource.StateChangeConf{
			Pending: []string{"provisioning"},
			Target:  []string{"provisioned"},
			Refresh: func() (interface{}, string, error) {
				s, err := bmClient.GetServer(serverId)
				if err != nil {
					return nil, "", err
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
		Severity:      diag.Warning,
		Summary:       fmt.Sprintf("Server Configuration Updates Required"),
		Detail:        "You will need to make changes on your server for networking changes to take effect",
		AttributePath: cty.GetAttrPath("network_ids"),
	})
	return refreshServer, networkDiagnostics
}
