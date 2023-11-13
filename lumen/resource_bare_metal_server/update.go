package resource_bare_metal_server

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-lumen/lumen/client"
	"terraform-provider-lumen/lumen/client/model/bare_metal"
	"time"
)

var updateTimeout = schema.DefaultTimeout(5 * time.Minute)

func updateContext(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	// Currently update can only be used for changing the server name
	if data.HasChange("name") {
		serverId := data.Id()
		updateRequest := bare_metal.ServerUpdateRequest{
			Name: data.Get("name").(string),
		}

		bmClient := i.(*client.Clients).BareMetal
		server, err := bmClient.UpdateServer(serverId, updateRequest)
		if err != nil {
			return diag.FromErr(err)
		}

		populateServerSchema(data, *server)
	}

	return nil
}
