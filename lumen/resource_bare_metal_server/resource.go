package resource_bare_metal_server

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-lumen/lumen/validation"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Description: "Provision lumen bare metal server",
		Timeouts: &schema.ResourceTimeout{
			Create: createTimeout,
			Read:   readTimeout,
			Delete: deleteTimeout,
			Update: updateTimeout,
		},
		CreateContext: createContext,
		ReadContext:   readContext,
		UpdateContext: updateContext,
		DeleteContext: deleteContext,
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
			"attach_networks": {
				Description: `List of existing networks to attach to the server being provisioned.  
If providing multiple values it will require you to make server configuration changes for change to take effect.`,
				Type:     schema.TypeList,
				Optional: true,
				ConflictsWith: []string{
					"network_name",
					"network_size_id",
				},
				ExactlyOneOf: []string{
					"network_name",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"assign_ipv6_address": {
							Description: `A boolean (true/false) value indicating whether or not to assign an IPv6 address 
for this server if using a dual stack network. Defaults to false if not set.`,
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"network_name": {
				Description: "The name of the new network to create, this is only used on initial creation.",
				Type:        schema.TypeString,
				Optional:    true,
				ConflictsWith: []string{
					"attach_networks",
				},
				RequiredWith: []string{
					"network_size_id",
				},
				ExactlyOneOf: []string{
					"attach_networks",
				},
			},
			"network_size_id": {
				Description: "The id of the network size being used for the new network, this is only used on initial creation.",
				Type:        schema.TypeString,
				Optional:    true,
				ConflictsWith: []string{
					"attach_networks",
				},
				RequiredWith: []string{
					"network_name",
				},
			},
			"network_type": {
				Description: "The type of network being provisioned for this server.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "INTERNET",
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
			"assign_ipv6_address": {
				Description: `A boolean (true/false) value indicating whether or not to assign an IPv6 address 
for this server if using a dual stack network. Defaults to false if not set.`,
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
						"ipv6": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vlan": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"assign_ipv6_address": {
							Type:     schema.TypeBool,
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
