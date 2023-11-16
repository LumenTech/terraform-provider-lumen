package resource_bare_metal_server

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-lumen/lumen/client"
	"terraform-provider-lumen/lumen/client/model/bare_metal"
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

		server, err := bmClient.UpdateServer(serverId, updateRequest)
		if err != nil {
			return diag.FromErr(err)
		}

		populateServerSchema(data, *server)
	}

	if data.HasChange("network_ids") {
		server, err := bmClient.GetServer(serverId)
		if err != nil {
			return diag.FromErr(err)
		}
		oldNetworks := convertNetworksToListOfNetworkIds(server.Networks)
		newNetworks := convertListOfInterfaceToListOfString(data.Get("network_ids").([]interface{}))

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

func detachNetworksAndWaitForCompletion(ctx context.Context, bmClient *client.BareMetalClient, serverId string, networkIds []string) (*bare_metal.Server, diag.Diagnostics) {
	var networkDiagnostics diag.Diagnostics
	for _, networkId := range networkIds {
		_, e := bmClient.RemoveNetwork(serverId, networkId)
		if e != nil {
			networkDiagnostics = append(networkDiagnostics, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Error detaching network %s", networkId),
				Detail:   fmt.Sprintf("Network %s errored on detachment reason - %s", networkId, e),
			})
		}
	}

	var refreshServer *bare_metal.Server
	stateChangeConf := &resource.StateChangeConf{
		Pending: []string{"detaching"},
		Target:  []string{"detached"},
		Refresh: func() (interface{}, string, error) {
			s, err := bmClient.GetServer(serverId)
			if err != nil {
				return nil, "", err
			}

			refreshServer = s
			status := "detached"
			for _, n := range s.Networks {
				for _, networkId := range networkIds {
					if n.NetworkID == networkId {
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
	return refreshServer, networkDiagnostics
}

// difference returns the elements in `a` that aren't in `b`.
func difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
