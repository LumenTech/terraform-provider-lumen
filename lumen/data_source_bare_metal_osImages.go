package lumen

import (
	"context"
	"fmt"
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
			osImages, diagnostics := bmClient.GetOsImages(data.Get("location_id").(string))
			if diagnostics.HasError() {
				return diagnostics
			}

			if err := data.Set("os_images", bare_metal.ConvertOSImagesToListMap(*osImages)); err != nil {
				return append(diagnostics, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("failed to set os_images: %s", err.Error()),
				})
			}
			data.SetId("os_images")
			return diagnostics
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
						"tier": {
							Description: "OS image tier that is used to match to server configuration tier for pricing",
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
