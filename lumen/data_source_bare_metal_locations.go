package lumen

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceBareMetalLocations() *schema.Resource {
	return &schema.Resource{
		//ReadContext: GetLocations,
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
