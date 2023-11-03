package lumen

import (
	"context"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strings"
	"terraform-provider-lumen/lumen/client/model/bare_metal"
	"terraform-provider-lumen/lumen/validation"
	"time"
)

var pendingServerStatuses = []string{"provisioning", "network_provisioned", "allocated", "configured", "unknown"}
var targetServerStatuses = []string{"provisioned", "failed", "error"}
var possibleServerStatus = append(pendingServerStatuses, targetServerStatuses...)

func ResourceBareMetalServer() *schema.Resource {
	return &schema.Resource{
		Description: "Provision lumen bare metal server",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
		},
		CreateContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			log.Printf("[INFO] Create Executed")
			client := i.(*Clients).BareMetal

			provisionRequest := bare_metal.ServerProvisionRequest{
				Name:          data.Get("name").(string),
				LocationID:    data.Get("location_id").(string),
				Configuration: data.Get("configuration_name").(string),
				OSImage:       data.Get("os_image_name").(string),
				Credentials: bare_metal.Credentials{
					Username:  data.Get("username").(string),
					Password:  data.Get("password").(string),
					PublicKey: data.Get("ssh_public_key").(string),
				},
			}

			networkId := data.Get("network_id").(string)
			if len(networkId) != 0 {
				provisionRequest.NetworkID = networkId
			} else {
				provisionRequest.NetworkRequest = &bare_metal.NetworkProvisionRequest{
					Name:          data.Get("network_name").(string),
					LocationID:    provisionRequest.LocationID,
					NetworkSizeID: data.Get("network_size_id").(string),
					NetworkType:   "INTERNET",
				}
			}

			server, err := client.ProvisionServer(provisionRequest)
			if err != nil {
				return diag.FromErr(err)
			}

			stateChangeConf := &resource.StateChangeConf{
				Pending: pendingServerStatuses,
				Target:  targetServerStatuses,
				Refresh: func() (interface{}, string, error) {
					s, err := client.GetServer(server.ID)
					if err != nil {
						return nil, "", err
					}

					found := false
					for _, status := range possibleServerStatus {
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
				Timeout:      90 * time.Minute,
				Delay:        4 * time.Minute,
				PollInterval: 30 * time.Second,
			}
			refreshResult, err := stateChangeConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.FromErr(err)
			}

			serv := refreshResult.(bare_metal.Server)
			data.SetId(serv.ID)
			populateServerSchema(data, serv)
			return nil
		},
		ReadContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			log.Printf("[INFO] Read Executed")
			client := i.(*Clients).BareMetal

			serverId := data.Id()
			server, err := client.GetServer(serverId)
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
			log.Printf("[INFO] Update Executed")
			return nil
		},
		DeleteContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			log.Printf("[INFO] Destroy Executed")
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
			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{
					"network_name",
					"network_size_id",
				},
			},
			"network_name": {
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{
					"network_id",
				},
			},
			"network_size_id": {
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{
					"network_id",
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
					"public_key",
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
					"public_key",
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
			"disks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"boot": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"disk_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
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

func populateServerSchema(d *schema.ResourceData, server bare_metal.Server) {
	d.Set("name", server.Name)
	d.Set("location_id", server.LocationID)
	d.Set("configuration_name", server.Configuration.Name)
	d.Set("os_image", server.OSImage)
	d.Set("machine_id", server.MachineID)
	d.Set("machine_name", server.MachineName)
	d.Set("location", server.Location)
	d.Set("configuration_cores", server.Configuration.Cores)
	d.Set("configuration_memory", server.Configuration.Memory)
	d.Set("configuration_storage", server.Configuration.Storage)
	d.Set("configuration_disks", server.Configuration.Disks)
	d.Set("configuration_nics", server.Configuration.NICs)
	d.Set("configuration_processors", server.Configuration.Processors)
	networks := make([]map[string]interface{}, len(server.Networks))
	for i, network := range server.Networks {
		networks[i] = map[string]interface{}{
			"id":             network.ID,
			"network_id":     network.NetworkID,
			"network_name":   network.NetworkName,
			"network_type":   network.NetworkType,
			"status":         network.Status,
			"status_message": network.StatusMessage,
			"ip":             network.IP,
			"vlan":           network.VLAN,
		}
	}
	d.Set("networks", networks)
	d.Set("status", server.Status)
	d.Set("status_message", server.StatusMessage)
	disks := make([]map[string]interface{}, len(server.Disks))
	for i, disk := range server.Disks {
		disks[i] = map[string]interface{}{
			"boot":      disk.Boot,
			"disk_type": disk.DiskType,
			"path":      disk.Path,
			"size":      disk.Size,
		}
	}
	d.Set("disks", disks)
	d.Set("boot_disk", server.BootDisk)
	d.Set("service_id", server.ServiceID)
	prices := make([]map[string]interface{}, len(server.Prices))
	for i, price := range server.Prices {
		prices[i] = map[string]interface{}{
			"type":  price.Type,
			"price": price.Price.String(),
		}
	}
	d.Set("prices", prices)
	d.Set("account_id", server.AccountID)
	d.Set("created", server.Created)
	d.Set("updated", server.Updated)
}
