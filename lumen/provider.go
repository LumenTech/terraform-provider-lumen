package lumen

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	/* User authentication schema */
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Lumen username",
				DefaultFunc: schema.EnvDefaultFunc("LUMEN_USERNAME", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Lumen password",
				DefaultFunc: schema.EnvDefaultFunc("LUMEN_PASSWORD", nil),
			},
			"api_access_token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Lumen API access token",
				DefaultFunc: schema.EnvDefaultFunc("LUMEN_API_ACCESS_TOKEN", nil),
			},
			"api_refresh_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Lumen API refresh token",
				DefaultFunc: schema.EnvDefaultFunc("LUMEN_API_REFRESH_TOKEN", nil),
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
	// Lumen API url
	apiUrl := "https://api.lumen.com/EdgeServices/v1/Compute/"
	// Lumen Auth url
	authUrl := "https://api.lumen.com/oauth/v1/token"
	// Lumen username
	username := d.Get("username").(string)
	if username == "" {
		return nil, diag.FromErr(fmt.Errorf("Lumen username"))
	}
	// Lumen password
	password := d.Get("password").(string)
	if password == "" {
		return nil, diag.FromErr(fmt.Errorf("Lumen password"))
	}
	// Lumen API access token
	apiAccessToken := d.Get("api_access_token").(string)
	if apiAccessToken == "" {
		return nil, diag.FromErr(fmt.Errorf("Lumen api access token cannot be empty"))
	}
	// Lumen API refresh token
	apiRefreshToken := d.Get("api_refresh_token").(string)
	/*
		if apiRefreshToken == "" {
			return nil, diag.FromErr(fmt.Errorf("Lumen api refresh token cannot be empty"))
		}*/

	// Populating client config
	config := Config{
		ApiUrl:          apiUrl,
		AuthUrl:         authUrl,
		Username:        username,
		Password:        password,
		ApiAccessToken:  apiAccessToken,
		ApiRefreshToken: apiRefreshToken,
	}
	return config.LumenClient()
}
