package lumen

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceBareMetalLocations() *schema.Resource {
	return &schema.Resource{
		Description: "Provides the list of bare metal locations",
		ReadContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			bmClient := i.(*Client).BareMetal
			locations, err := bmClient.GetLocations()
			if err != nil {
				return diag.FromErr(err)
			}

			if err := data.Set("locations", locations.ToMapList()); err != nil {
				return diag.FromErr(err)
			}
			data.SetId("locations")
			return nil
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
