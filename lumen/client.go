package lumen

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"terraform-provider-lumen/lumen/client"
)

type Client struct {
	Morpheus  *client.MorpheusClient
	BareMetal *client.BareMetalClient
}

type Config struct {
	ApiUrl          string
	AuthUrl         string
	Username        string
	Password        string
	AccountAlias    string
	ApiAccessToken  string
	ApiRefreshToken string

	client *Client
}

// Configuring Lumen Provider client
func (c *Config) LumenClient() (*Client, diag.Diagnostics) {
	if c.client == nil {
		mClient := client.NewMorpheusClient(c.ApiUrl, c.AuthUrl)
		mClient.SetCredsAndTokens(c.Username, c.Password, c.ApiAccessToken, c.ApiRefreshToken)

		c.client = &Client{
			Morpheus:  mClient,
			BareMetal: client.NewBareMetalClient(c.Username, c.Password, c.AccountAlias),
		}
	}
	return c.client, nil
}
