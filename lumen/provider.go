package lumen

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	/* User authentication schema */
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Lumen API endpoint URL where requests will be directed",
				DefaultFunc: schema.EnvDefaultFunc("LUMEN_API_URL", nil),
			},
			"access_token": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				Description:   "Access Token of Lumen API user, instead of authenticating with username and password",
				DefaultFunc:   schema.EnvDefaultFunc("LUMEN_API_TOKEN", nil),
				ConflictsWith: []string{"username", "password"},
			},
			"username": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Username of Lumen API user for authentication",
				DefaultFunc:   schema.EnvDefaultFunc("LUMEN_USERNAME", nil),
				ConflictsWith: []string{"access_token"},
			},
			"password": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				Description:   "Password of Lumen API user for authentication",
				DefaultFunc:   schema.EnvDefaultFunc("LUMEN_PASSWORD", nil),
				ConflictsWith: []string{"access_token"},
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			/*
				resource_bare_metal_instance :
					- create instance
					- read created instance
					- delete instance
					- update instance
			*/
			"lumen_bare_metal_instance": ResourceBareMetalInstance(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			/*
				dataSourceBareMetalAllInstances : lists all instances currently with tenant.
				dataSourceBareMetalInstanceId : lists details for a particular instance based on instance id.
				dataSourceBareMetalInstanceName : lists details for a particular instance based on instance name.
			*/
			"lumen_bare_metal_instances":     DataSourceBareMetalAllInstances(),
			"lumen_bare_metal_instance_id":   DataSourceBareMetalInstanceId(),
			"lumen_bare_metal_instance_name": DataSourceBareMetalInstanceName(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	// Lumen url and user auth based on tenant
	config := Config{
		Url:         d.Get("url").(string),
		AccessToken: d.Get("access_token").(string),
		Username:    d.Get("username").(string),
		Password:    d.Get("password").(string),
	}
	return config.Client()
}
