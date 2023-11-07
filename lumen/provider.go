package lumen

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const CustomerDeprecationNotice = "Resource/Argument is deprecated once the customer has migrated to the new version of Lumen Edge Bare Metal API."

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
				Optional:    true,
				Sensitive:   true,
				Description: "Lumen API access token",
				DefaultFunc: schema.EnvDefaultFunc("LUMEN_API_ACCESS_TOKEN", nil),
				Deprecated:  CustomerDeprecationNotice,
			},
			"api_refresh_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Lumen API refresh token",
				DefaultFunc: schema.EnvDefaultFunc("LUMEN_API_REFRESH_TOKEN", nil),
				Deprecated:  CustomerDeprecationNotice,
			},
			"account_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Lumen customer account number",
				DefaultFunc: schema.EnvDefaultFunc("LUMEN_ACCOUNT_NUMBER", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"lumen_bare_metal_server":  ResourceBareMetalServer(),
			"lumen_bare_metal_network": ResourceBareMetalNetwork(),
			/*
				ResourceBareMetalInstance :
					- create bare metal instance
					- read created instance
					- delete bare metal instance
					- update bare instance
				ResourceNetworkInstance:
					- create network instance
					- delete network instance
			*/
			"lumen_bare_metal_instance": ResourceBareMetalInstance(),
			"lumen_network_instance":    ResourceNetworkInstance(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"lumen_bare_metal_configurations": DataSourceBareMetalConfigurations(),
			"lumen_bare_metal_locations":      DataSourceBareMetalLocations(),
			"lumen_bare_metal_network_sizes":  DataSourceBareMetalNetworkSizes(),
			"lumen_bare_metal_os_images":      DataSourceBareMetalOsImages(),
			/*
				Deprecated Resources
				DataSourceBareMetalAllInstances : lists all instances currently with tenant.
				DataSourceBareMetalInstanceId : lists details for a particular instance based on instance id.
				DataSourceBareMetalInstanceName : lists details for a particular instance based on instance name.
				DataSourceNetworkAllInstances : lists all network instances currently under a tenant.
				DataSourceNetworkInstanceId : lists details for network instance based on instance id.
				DataSourceNetworkInstanceName : lists details for network instance based on instance name.

			*/
			"lumen_bare_metal_instances":     DataSourceBareMetalAllInstances(),
			"lumen_bare_metal_instance_id":   DataSourceBareMetalInstanceId(),
			"lumen_bare_metal_instance_name": DataSourceBareMetalInstanceName(),
			"lumen_network_instances":        DataSourceNetworkAllInstances(),
			"lumen_network_instance_id":      DataSourceNetworkInstanceId(),
			"lumen_network_instance_name":    DataSourceNetworkInstanceName(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apigeeBaseURL := "https://api.lumen.com"
	// Lumen API url
	morpheusBareMetalApiUrl := fmt.Sprintf("%s/EdgeServices/v1/Compute/", apigeeBaseURL)
	// Lumen Auth url
	authUrl := fmt.Sprintf("%s/oauth/v1/token", apigeeBaseURL)
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
	//if apiAccessToken == "" {
	//	return nil, diag.FromErr(fmt.Errorf("Lumen api access token cannot be empty"))
	//}
	// Lumen API refresh token
	apiRefreshToken := d.Get("api_refresh_token").(string)
	/*
		if apiRefreshToken == "" {
			return nil, diag.FromErr(fmt.Errorf("Lumen api refresh token cannot be empty"))
		}*/
	accountNumber := d.Get("account_number").(string)

	// Populating clients config
	config := Config{
		ApigeeBaseURL:           apigeeBaseURL,
		MorpheusBareMetalApiUrl: morpheusBareMetalApiUrl,
		AuthUrl:                 authUrl,
		Username:                username,
		Password:                password,
		AccountNumber:           accountNumber,
		ApiAccessToken:          apiAccessToken,
		ApiRefreshToken:         apiRefreshToken,
	}
	return config.LumenClients()
}
