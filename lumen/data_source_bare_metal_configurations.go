package lumen

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceBareMetalConfigurations() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a list of bare metal configurations at a specific location",
		ReadContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			bmClient := i.(*Clients).BareMetal
			configurations, err := bmClient.GetConfigurations(data.Get("location_id").(string))
			if err != nil {
				return diag.FromErr(err)
			}

			if err := data.Set("configurations", configurations.ToMapList()); err != nil {
				return diag.FromErr(err)
			}
			data.SetId("configurations")
			return nil
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
						"machineCount": {
							Description: "The number of machines in this configuration",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"price": {
							Description: "The price for this configuration",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
