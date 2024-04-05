package lumen

import (
	"context"
	"fmt"
	"terraform-provider-lumen/lumen/client"
	"terraform-provider-lumen/lumen/resource_bare_metal_server"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	/* User authentication schema */
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"consumer_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Consumer key for Lumen API",
				DefaultFunc: schema.EnvDefaultFunc("LUMEN_CONSUMER_KEY", nil),
			},
			"consumer_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Consumer secret for Lumen API",
				DefaultFunc: schema.EnvDefaultFunc("LUMEN_CONSUMER_SECRET", nil),
			},
			"account_number": {
				Type:        schema.TypeString,
				Required:    true,
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
	consumerKey := d.Get("consumer_key").(string)
	if consumerKey == "" {
		return nil, diag.FromErr(fmt.Errorf("Consumer key not found"))
	}
	consumerSecret := d.Get("consumer_secret").(string)
	if consumerSecret == "" {
		return nil, diag.FromErr(fmt.Errorf("Consumer secret not found"))
	}
	accountNumber := d.Get("account_number").(string)

	// Populating clients config
	config := client.Config{
		ApigeeBaseURL:  apigeeBaseURL,
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		AccountNumber:  accountNumber,
	}
	return config.LumenClients()
}
