package resource_bare_metal_server

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"terraform-provider-lumen/lumen/client"
	"time"
)

var readTimeout = schema.DefaultTimeout(5 * time.Minute)

func readContext(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	bmClient := i.(*client.Clients).BareMetal

	serverId := data.Id()
	server, err := bmClient.GetServer(serverId)
	if err != nil {
		return diag.FromErr(err)
	}

	if server == nil {
		log.Printf("[DEBUG] Bare metal server %s was not found - removing from state!", serverId)
		data.SetId("")
		return nil
	}

	populateServerSchema(data, *server)
	return nil
}
