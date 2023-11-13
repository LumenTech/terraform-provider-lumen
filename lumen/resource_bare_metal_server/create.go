package resource_bare_metal_server

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
	"terraform-provider-lumen/lumen/client"
	"terraform-provider-lumen/lumen/client/model/bare_metal"
	"time"
)

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

	networkIds := data.Get("network_ids").([]interface{})
	if len(networkIds) != 0 {
		provisionRequest.NetworkID = networkIds[0].(string)
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
	serv := *server
	data.SetId(serv.ID)
	populateServerSchema(data, serv)

	// Attach Additional Network
	networkWarnings := diag.Diagnostics{}
	if len(networkIds) > 1 {
		for _, networkId := range networkIds[1:] {
			_, warningDiagnostic := attachNetwork(bmClient, server.ID, networkId.(string))
			if warningDiagnostic != nil {
				networkWarnings = append(networkWarnings, *warningDiagnostic)
			}
		}
	}
	// TODO: Poll until all networks have finished
	return networkWarnings
}

var pendingServerStatuses = []string{"provisioning", "network_provisioned", "allocated", "configured", "unknown"}
var targetServerStatuses = []string{"provisioned", "failed", "error"}
var possibleServerStatus = append(pendingServerStatuses, targetServerStatuses...)

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

func attachNetwork(bmClient *client.BareMetalClient, serverId, networkId string) (*bare_metal.Server, *diag.Diagnostic) {
	server, err := bmClient.AttachNetwork(serverId, networkId)
	if err != nil {
		return nil, &diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  fmt.Sprintf("Error attaching network %s", networkId),
			Detail:   fmt.Sprintf("Network %s errored on attachment reason - %s", networkId, err),
		}
	}
	return server, nil
}
