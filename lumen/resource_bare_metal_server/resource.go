package resource_bare_metal_server

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strings"
	"terraform-provider-lumen/lumen/client"
	"terraform-provider-lumen/lumen/client/model/bare_metal"
	"terraform-provider-lumen/lumen/validation"
	"time"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Description: "Provision lumen bare metal server",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
		},
		CreateContext: createContext,
		ReadContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
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
		},
		UpdateContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			// Currently update can only be used for changing the server name
			if data.HasChange("name") {
				serverId := data.Id()
				updateRequest := bare_metal.ServerUpdateRequest{
					Name: data.Get("name").(string),
				}

				bmClient := i.(*client.Clients).BareMetal
				server, err := bmClient.UpdateServer(serverId, updateRequest)
				if err != nil {
					return diag.FromErr(err)
				}

				populateServerSchema(data, *server)
			}

			return nil
		},
		DeleteContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			bmClient := i.(*client.Clients).BareMetal
			serverId := data.Id()
			server, err := bmClient.DeleteServer(serverId)
			if err != nil {
				return diag.FromErr(err)
			} else if server == nil || strings.ToLower(server.Status) == "released" {
				data.SetId("")
				return nil
			}

			stateChangeConf := &resource.StateChangeConf{
				Pending: []string{"releasing", "networking_removed"},
				Target:  []string{"released"},
				Refresh: func() (interface{}, string, error) {
					s, e := bmClient.GetServer(serverId)
					if e != nil {
						return nil, "", err
					} else if s == nil {
						s = server
						s.Status = "released"
					}
					return *s, strings.ToLower(s.Status), nil
				},
				Timeout:      90 * time.Minute,
				Delay:        2 * time.Minute,
				PollInterval: 30 * time.Second,
			}
			refreshResult, err := stateChangeConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.FromErr(err)
			}

			fmt.Printf(
				"[DEBUG] Deleted server (%s) final status (%s)",
				serverId,
				refreshResult.(bare_metal.Server).Status,
			)
			data.SetId("")
			return nil
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					if validationError := validation.ValidateBareMetalServerName(i.(string)); validationError != nil {
						return diag.FromErr(validationError)
					}
					return nil
				},
			},
			"location_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"configuration_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"os_image_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_ids": {
				Description: `List of existing networks to attach to the server being provisioned.  
If providing multiple values it will require you to make server configuration changes for change to take effect.`,
				Type:     schema.TypeList,
				Optional: true,
				ConflictsWith: []string{
					"network_name",
					"network_size_id",
				},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"network_name": {
				Description: "The name of the new network to create, this is only used on initial creation.",
				Type:        schema.TypeString,
				Optional:    true,
				ConflictsWith: []string{
					"network_ids",
				},
			},
			"network_size_id": {
				Description: "The id of the network size being used for the new network, this is only used on initial creation.",
				Type:        schema.TypeString,
				Optional:    true,
				ConflictsWith: []string{
					"network_ids",
				},
			},
			"username": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					if err := validation.ValidateBareMetalUsername(i.(string)); err != nil {
						return diag.FromErr(err)
					}
					return nil
				},
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				AtLeastOneOf: []string{
					"password",
					"ssh_public_key",
				},
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					if err := validation.ValidateBareMetalPassword(i.(string)); err != nil {
						return diag.FromErr(err)
					}
					return nil
				},
			},
			"ssh_public_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				AtLeastOneOf: []string{
					"password",
					"ssh_public_key",
				},
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"machine_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"machine_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"configuration_cores": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"configuration_memory": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"configuration_storage": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"configuration_disks": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"configuration_nics": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"configuration_processors": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"networks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vlan": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"boot_disk": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_id": {
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
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
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
