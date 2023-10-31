package lumen

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceBareMetalNetworkSizes() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a list of network sizes at a specific location",
		ReadContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			bmClient := i.(*Clients).BareMetal
			networkSizes, err := bmClient.GetNetworkSizes(data.Get("location_id").(string))
			if err != nil {
				return diag.FromErr(err)
			}

			if err := data.Set("network_sizes", networkSizes.ToMapList()); err != nil {
				return diag.FromErr(err)
			}
			data.SetId("network_sizes")
			return nil
		},
		Schema: map[string]*schema.Schema{
			"location_id": {
				Description: "The id of a location",
				Type:        schema.TypeString,
				Required:    true,
			},
			"network_sizes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "The id of a network size",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "The name of this network size",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"cidr": {
							Description: "The CIDR for this network size",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"available_ips": {
							Description: "The number of available IPs for this network size",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"price": {
							Description: "The price for this network size",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
