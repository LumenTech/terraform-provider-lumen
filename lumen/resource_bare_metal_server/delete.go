package resource_bare_metal_server

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
	"terraform-provider-lumen/lumen/client"
	"terraform-provider-lumen/lumen/helper"
	"time"
)

var deleteTimeout = schema.DefaultTimeout(30 * time.Minute)

func deleteContext(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	bmClient := i.(*client.Clients).BareMetal
	serverId := data.Id()
	server, diagnostics := bmClient.DeleteServer(serverId)
	if diagnostics.HasError() {
		return diagnostics
	}

	if server != nil && strings.ToLower(server.Status) != "released" {
		_, waitError := waitForServerDeletion(ctx, bmClient, serverId)
		if waitError != nil {
			return append(diagnostics, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  waitError.Error(),
			})
		}
	}
	data.SetId("")
	return diagnostics
}

var deleteServerPendingStatus = []string{"releasing", "networking_removed", "unknown"}
var deleteServerTargetStatus = []string{"released"}
var deleteServerStatuses = append(append(deleteServerPendingStatus, deleteServerTargetStatus...), "failed", "error")

func waitForServerDeletion(ctx context.Context, bmClient *client.BareMetalClient, serverId string) (interface{}, error) {
	stateChangeConf := &resource.StateChangeConf{
		Pending: deleteServerPendingStatus,
		Target:  deleteServerTargetStatus,
		Refresh: func() (interface{}, string, error) {
			server, getDiagnostics := bmClient.GetServer(serverId)
			if err := helper.ExtractDiagnosticErrorIfPresent(getDiagnostics); err != nil {
				return nil, "", fmt.Errorf(err.Summary)
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
