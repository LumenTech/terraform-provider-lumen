package lumen

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-lumen/lumen/client"
	"terraform-provider-lumen/lumen/resource_bare_metal_server"
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
			"account_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Lumen customer account number (required for new versions of bare metal resources/data sources)",
				DefaultFunc: schema.EnvDefaultFunc("LUMEN_ACCOUNT_NUMBER", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"lumen_bare_metal_server":  resource_bare_metal_server.Resource(),
			"lumen_bare_metal_network": ResourceBareMetalNetwork(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"lumen_bare_metal_configurations": DataSourceBareMetalConfigurations(),
			"lumen_bare_metal_locations":      DataSourceBareMetalLocations(),
			"lumen_bare_metal_network_sizes":  DataSourceBareMetalNetworkSizes(),
			"lumen_bare_metal_os_images":      DataSourceBareMetalOsImages(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apigeeBaseURL := "https://api.lumen.com"
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
	accountNumber := d.Get("account_number").(string)

	// Populating clients config
	config := client.Config{
		ApigeeBaseURL: apigeeBaseURL,
		Username:      username,
		Password:      password,
		AccountNumber: accountNumber,
	}
	return config.LumenClients()
}
