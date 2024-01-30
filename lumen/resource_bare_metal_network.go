package lumen

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/net/context"
	"log"
	"strings"
	client2 "terraform-provider-lumen/lumen/client"
	"terraform-provider-lumen/lumen/client/model/bare_metal"
	"time"
)

var pendingNetworkStatuses = []string{"provisioning", "unknown"}
var targetNetworkStatuses = []string{"provisioned"}
var possibleNetworkStatus = append(append(pendingNetworkStatuses, targetNetworkStatuses...), "failed")

func ResourceBareMetalNetwork() *schema.Resource {
	return &schema.Resource{
		Description: "Provision lumen bare metal network",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		CreateContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			client := i.(*client2.Clients).BareMetal

			provisionRequest := bare_metal.NetworkProvisionRequest{
				Name:          data.Get("name").(string),
				LocationID:    data.Get("location_id").(string),
				NetworkSizeID: data.Get("network_size_id").(string),
				NetworkType:   data.Get("network_type").(string),
			}

			network, err := client.ProvisionNetwork(provisionRequest)
			if err != nil {
				return diag.FromErr(err)
			}

			stateChangeConf := &resource.StateChangeConf{
				Pending: pendingNetworkStatuses,
				Target:  targetNetworkStatuses,
				Refresh: func() (interface{}, string, error) {
					s, err := client.GetNetwork(network.ID)
					if err != nil {
						return nil, "", err
					}

					found := false
					for _, status := range possibleNetworkStatus {
						// This is to avoid failure if a new status is added due to the logic
						// of the polling mechanism if the status isn't in the pending or target list
						// it is considered a failure and errors out immediate.
						if status == strings.ToLower(s.Status) {
							found = true
							break
						}
					}

					if !found {
						s.Status = "unknown"
					}
					return *s, strings.ToLower(s.Status), nil
				},
				Timeout:      10 * time.Minute,
				Delay:        2 * time.Minute,
				PollInterval: 30 * time.Second,
			}
			refreshResult, err := stateChangeConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.FromErr(err)
			}

			retNetwork := refreshResult.(bare_metal.Network)
			data.SetId(retNetwork.ID)
			populateNetworkSchema(data, retNetwork)
			return nil
		},
		ReadContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			client := i.(*client2.Clients).BareMetal

			networkId := data.Id()
			network, err := client.GetNetwork(networkId)
			if err != nil {
				return diag.FromErr(err)
			}

			if network == nil {
				log.Printf("[DEBUG] Network %s was not found - removing from state!", networkId)
				data.SetId("")
				return nil
			}

			populateNetworkSchema(data, *network)
			return nil
		},
		DeleteContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			client := i.(*client2.Clients).BareMetal
			networkId := data.Id()
			network, err := client.DeleteNetwork(networkId)
			if err != nil {
				return diag.FromErr(err)
			} else if network == nil {
				data.SetId("")
				return nil
			}

			stateChangeConf := &resource.StateChangeConf{
				Pending: []string{"deleting"},
				Target:  []string{"deleted"},
				Refresh: func() (interface{}, string, error) {
					n, e := client.GetNetwork(networkId)
					if e != nil {
						return nil, "", err
					} else if n == nil {
						n = network
						n.Status = "deleted"
					}
					return *n, strings.ToLower(n.Status), nil
				},
				Timeout:      10 * time.Minute,
				Delay:        30 * time.Second,
				PollInterval: 30 * time.Second,
			}
			refreshResult, err := stateChangeConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.FromErr(err)
			}

			fmt.Printf(
				"[DEBUG] Deleted network (%s) final status (%s)",
				networkId,
				refreshResult.(bare_metal.Network).Status,
			)
			data.SetId("")
			return nil
		},
		UpdateContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			// Currently update can only be used for changing the network name
			if data.HasChange("name") {
				networkId := data.Id()
				updateRequest := bare_metal.NetworkUpdateRequest{
					Name: data.Get("name").(string),
				}

				client := i.(*client2.Clients).BareMetal
				network, err := client.UpdateNetwork(networkId, updateRequest)
				if err != nil {
					return diag.FromErr(err)
				}

				populateNetworkSchema(data, *network)
			}

			return nil
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"location_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_size_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_type": {
				Type:     schema.TypeString,
				Default:  "INTERNET",
				ForceNew: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_block": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_block": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"available_ips": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_ips": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"prices": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"price": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func populateNetworkSchema(d *schema.ResourceData, network bare_metal.Network) {
	d.Set("name", network.Name)
	d.Set("account_id", network.AccountID)
	d.Set("service_id", network.ServiceID)
	d.Set("location", network.Location)
	d.Set("location_id", network.LocationID)
	d.Set("ip_block", network.IPBlock)
	d.Set("ipv6_block", network.IPV6Block)
	d.Set("gateway", network.Gateway)
	d.Set("available_ips", network.AvailableIPs)
	d.Set("total_ips", network.TotalIPs)
	d.Set("type", network.Type)
	d.Set("status", network.Status)
	prices := make([]map[string]interface{}, len(network.Prices))
	for i, price := range network.Prices {
		prices[i] = map[string]interface{}{
			"type":  price.Type,
			"price": price.Price.String(),
		}
	}
	d.Set("prices", prices)
	d.Set("created", network.Created)
	d.Set("updated", network.Updated)
}
