package lumen

import (
	"context"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-lumen/lumen/validation"
)

func ResourceBareMetalServer() *schema.Resource {
	return &schema.Resource{
		Description: "Provision lumen bare metal server",
		CreateContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			return nil
		},
		ReadContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			return nil
		},
		UpdateContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			return nil
		},
		DeleteContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
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
			"os_image": {
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
			"public_key": {
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
			"configuration": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cores": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disks": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"nics": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"processors": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"networks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: schema.Resource{
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
			"status":         {},
			"status_message": {},
			"disks":          {},
			"boot_disk":      {},
			"service_id":     {},
			"prices":         {},
			"account_id":     {},
			"created":        {},
			"updated":        {},
		},
	}
}
