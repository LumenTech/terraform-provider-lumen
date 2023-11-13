package lumen

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-lumen/lumen/client"
	"terraform-provider-lumen/lumen/client/model/bare_metal"
)

func DataSourceBareMetalOsImages() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a list of available os images at a specific location",
		ReadContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			bmClient := i.(*client.Clients).BareMetal
			osImages, err := bmClient.GetOsImages(data.Get("location_id").(string))
			if err != nil {
				return diag.FromErr(err)
			}

			if err := data.Set("os_images", bare_metal.ConvertToListMap(*osImages)); err != nil {
				return diag.FromErr(err)
			}
			data.SetId("os_images")
			return nil
		},
		Schema: map[string]*schema.Schema{
			"location_id": {
				Description: "The id of a location",
				Type:        schema.TypeString,
				Required:    true,
			},
			"os_images": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "The name of this OS image",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"price": {
							Description: "The price for using this OS image",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
