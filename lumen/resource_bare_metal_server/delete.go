package resource_bare_metal_server

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
	"terraform-provider-lumen/lumen/client"
	"time"
)

var deleteTimeout = schema.DefaultTimeout(30 * time.Minute)

func deleteContext(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	bmClient := i.(*client.Clients).BareMetal
	serverId := data.Id()
	server, err := bmClient.DeleteServer(serverId)
	if err != nil {
		return diag.FromErr(err)
	}

	if server != nil && strings.ToLower(server.Status) != "released" {
		_, waitError := waitForServerDeletion(ctx, bmClient, serverId)
		if waitError != nil {
			return diag.FromErr(waitError)
		}
	}
	data.SetId("")
	return nil
}

var deleteServerPendingStatus = []string{"releasing", "networking_removed", "unknown"}
var deleteServerTargetStatus = []string{"released"}
var deleteServerStatuses = append(append(deleteServerPendingStatus, deleteServerTargetStatus...), "failed", "error")

func waitForServerDeletion(ctx context.Context, bmClient *client.BareMetalClient, serverId string) (interface{}, error) {
	stateChangeConf := &resource.StateChangeConf{
		Pending: deleteServerPendingStatus,
		Target:  deleteServerTargetStatus,
		Refresh: func() (interface{}, string, error) {
			server, err := bmClient.GetServer(serverId)
			if err != nil {
				return nil, "", err
			}

			status := "unknown"
			if server == nil {
				status = "released"
			} else {
				for _, deleteStatus := range deleteServerStatuses {
					if deleteStatus == strings.ToLower(server.Status) {
						status = server.Status
						break
					}
				}
			}

			return server, strings.ToLower(status), nil
		},
		Timeout:      *deleteTimeout,
		Delay:        2 * time.Minute,
		PollInterval: 30 * time.Second,
	}
	return stateChangeConf.WaitForStateContext(ctx)
}
