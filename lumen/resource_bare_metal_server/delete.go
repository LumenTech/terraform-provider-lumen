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

var deleteTimeout = schema.DefaultTimeout(90 * time.Minute)

func deleteContext(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	bmClient := i.(*client.Clients).BareMetal
	serverId := data.Id()
	server, err := bmClient.DeleteServer(serverId)
	if err != nil {
		return diag.FromErr(err)
	} else if server == nil || strings.ToLower(server.Status) == "released" {
		data.SetId("")
		return nil
	}

	refreshResult, waitError := waitForServerDeletion(ctx, bmClient, serverId)
	if waitError != nil {
		return diag.FromErr(waitError)
	}

	fmt.Printf(
		"[DEBUG] Deleted server (%s) final status (%s)",
		serverId,
		refreshResult.(bare_metal.Server).Status,
	)
	data.SetId("")
	return nil
}

var deleteServerPendingStatus = []string{"releasing", "networking_removed", "unknown"}
var deleteServerTargetStatus = []string{"released"}
var deleteServerStatuses = append(deleteServerPendingStatus, deleteServerTargetStatus...)

func waitForServerDeletion(ctx context.Context, bmClient *client.BareMetalClient, serverId string) (interface{}, error) {
	stateChangeConf := &resource.StateChangeConf{
		Pending: deleteServerPendingStatus,
		Target:  deleteServerTargetStatus,
		Refresh: func() (interface{}, string, error) {
			s, e := bmClient.GetServer(serverId)
			if e != nil {
				return nil, "", e
			} else if s == nil {
				return nil, "released", nil
			}

			found := false
			for _, status := range deleteServerStatuses {
				if status == strings.ToLower(s.Status) || strings.ToLower(s.Status) == "failed" {
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
		Delay:        2 * time.Minute,
		PollInterval: 30 * time.Second,
	}
	return stateChangeConf.WaitForStateContext(ctx)
}
