package lumen

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-lumen/lumen/client"
	"terraform-provider-lumen/lumen/client/model/bare_metal"
)

func DataSourceBareMetalLocations() *schema.Resource {
	return &schema.Resource{
		Description: "Provides the list of bare metal locations",
		ReadContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			bmClient := i.(*client.Clients).BareMetal
			locations, diagnostics := bmClient.GetLocations()
			if diagnostics.HasError() {
				return diagnostics
			}

			if err := data.Set("locations", bare_metal.ConvertToListMap(*locations)); err != nil {
				return append(diagnostics, diag.Diagnostic{
					Severity: diag.Error,
					Detail:   fmt.Sprintf("failed to set locations: %s", err.Error()),
				})
			}
			data.SetId("locations")
			return diagnostics
		},
		Schema: map[string]*schema.Schema{
			"locations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "The location id",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "The name of the location",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"status": {
							Description: "The status of the location",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"region": {
							Description: "The region the location is in",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
