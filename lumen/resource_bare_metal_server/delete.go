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

	stateChangeConf := &resource.StateChangeConf{
		Pending: []string{"releasing", "networking_removed"},
		Target:  []string{"released"},
		Refresh: func() (interface{}, string, error) {
			s, e := bmClient.GetServer(serverId)
			if e != nil {
				return nil, "", err
			} else if s == nil {
				s = server
				s.Status = "released"
			}
			return *s, strings.ToLower(s.Status), nil
		},
		Timeout:      90 * time.Minute,
		Delay:        2 * time.Minute,
		PollInterval: 30 * time.Second,
	}
	refreshResult, err := stateChangeConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	fmt.Printf(
		"[DEBUG] Deleted server (%s) final status (%s)",
		serverId,
		refreshResult.(bare_metal.Server).Status,
	)
	data.SetId("")
	return nil
}
