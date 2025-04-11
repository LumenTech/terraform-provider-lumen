package lumen

import (
	"context"
	"fmt"
	"terraform-provider-lumen/lumen/client"
	"terraform-provider-lumen/lumen/client/model/bare_metal"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceBareMetalConfigurations() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a list of bare metal configurations at a specific location",
		ReadContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			bmClient := i.(*client.Clients).BareMetal
			configurations, diagnostics := bmClient.GetConfigurations(data.Get("location_id").(string))
			if diagnostics.HasError() {
				return diagnostics
			}

			if err := data.Set("configurations", bare_metal.ConvertToListMap(*configurations)); err != nil {
				return append(diagnostics, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("error setting configurations: %s", err.Error()),
				})
			}
			data.SetId("configurations")
			return diagnostics
		},
		Schema: map[string]*schema.Schema{
			"location_id": {
				Description: "The id of a location",
				Type:        schema.TypeString,
				Required:    true,
			},
			"configurations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "The type of configuration (ie small, medium, large)",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"display_name": {
							Description: "The display name of the configuration",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"cores": {
							Description: "The number of cores in this configuration",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"memory": {
							Description: "The memory available for this configuration",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"storage": {
							Description: "The storage available for this configuration",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"disks": {
							Description: "The number of disks in this configuration",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"nics": {
							Description: "The number of NICs in this configuration",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"processors": {
							Description: "The number of processors in this configuration",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"machine_count": {
							Description: "The number of machines in this configuration",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"price": {
							Description: "The price for this configuration",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"tier": {
							Description: "Tier tier of this configuration used for os pricing",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
