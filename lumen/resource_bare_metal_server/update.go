package resource_bare_metal_server

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-lumen/lumen/client"
	"terraform-provider-lumen/lumen/client/model/bare_metal"
	"terraform-provider-lumen/lumen/validation"
	"time"
)

var updateTimeout = schema.DefaultTimeout(5 * time.Minute)

func updateContext(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	serverId := data.Id()
	bmClient := i.(*client.Clients).BareMetal
	if data.HasChange("name") {
		updateRequest := bare_metal.ServerUpdateRequest{
			Name: data.Get("name").(string),
		}

		server, diagnostics := bmClient.UpdateServer(serverId, updateRequest)
		if diagnostics.HasError() {
			return diagnostics
		}

		populateServerSchema(data, *server)
	}

	if data.HasChange("attach_networks") {
		newNetworks := convertDataToAttachedNetworks(data.Get("attach_networks").([]interface{}))
		if validationError := validation.ValidateBareMetalNetworkIds(newNetworks); validationError != nil {
			return diag.FromErr(validationError)
		}

		server, getDiagnostics := bmClient.GetServer(serverId)
		if getDiagnostics.HasError() {
			return getDiagnostics
		}
		oldNetworks := convertNetworksToListOfAttachNetworks(server.Networks)

		var networkDiag diag.Diagnostics
		attachNetworks := difference(newNetworks, oldNetworks)
		if len(attachNetworks) > 0 {
			refreshServer, attachNetworkDiagnostics := attachNetworksAndWaitForCompletion(ctx, bmClient, serverId, attachNetworks)
			if refreshServer != nil {
				populateServerSchema(data, *refreshServer)
			}
			networkDiag = append(networkDiag, attachNetworkDiagnostics...)
		}

		detachNetworks := difference(oldNetworks, newNetworks)
		if len(detachNetworks) > 0 {
			refreshServer2, detachNetworkDiagnostics := detachNetworksAndWaitForCompletion(ctx, bmClient, serverId, detachNetworks)
			if refreshServer2 != nil {
				populateServerSchema(data, *refreshServer2)
			}
			networkDiag = append(networkDiag, detachNetworkDiagnostics...)
		}
		return networkDiag
	}

	return nil
}

func detachNetworksAndWaitForCompletion(ctx context.Context, bmClient *client.BareMetalClient, serverId string, networks []bare_metal.AttachNetwork) (*bare_metal.Server, diag.Diagnostics) {
	var networkDiagnostics diag.Diagnostics
	for _, network := range networks {
		_, diagnostics := bmClient.RemoveNetwork(serverId, network.NetworkID)
		if diagnostics.HasError() {
			reason := ""
			for _, d := range diagnostics {
				if d.Severity == diag.Error {
					reason = d.Summary
				}
			}

			networkDiagnostics = append(networkDiagnostics, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Error detaching network %s", network.NetworkID),
				Detail:   fmt.Sprintf("Network %s errored on detachment reason - %s", network.NetworkID, reason),
			})
		}
	}

	var refreshServer *bare_metal.Server
	stateChangeConf := &resource.StateChangeConf{
		Pending: []string{"detaching"},
		Target:  []string{"detached"},
		Refresh: func() (interface{}, string, error) {
			s, getDiagnostics := bmClient.GetServer(serverId)
			if getDiagnostics.HasError() {
				return nil, "", fmt.Errorf("error getting server %s details", serverId)
			}

			refreshServer = s
			status := "detached"
			for _, n := range s.Networks {
				for _, network := range networks {
					if n.NetworkID == network.NetworkID {
						status = "detaching"
						break
					}
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
			Summary:  "Polling server for networking statuses failed",
			Detail:   waitError.Error(),
		})
	}
	networkDiagnostics = append(networkDiagnostics, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Server Configuration Updates Required",
		Detail: `Automation configures Lumen networking infrastructure only, but not server configuration.
Removing a network from an existing server will require you to make configuration changes on your server.`,
		AttributePath: cty.GetAttrPath("network_ids"),
	})
	return refreshServer, networkDiagnostics
}

// difference returns the elements in `a` that aren't in `b`.
func difference(a []bare_metal.AttachNetwork, b []bare_metal.AttachNetwork) []bare_metal.AttachNetwork {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x.NetworkID] = struct{}{}
	}
	var diff []bare_metal.AttachNetwork
	for _, x := range a {
		if _, found := mb[x.NetworkID]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
