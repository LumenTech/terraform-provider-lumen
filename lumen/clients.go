package lumen

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"terraform-provider-lumen/lumen/client"
)

type Clients struct {
	Morpheus  *client.MorpheusClient
	BareMetal *client.BareMetalClient
}

type Config struct {
	ApigeeBaseURL           string
	MorpheusBareMetalApiUrl string
	AuthUrl                 string
	Username                string
	Password                string
	AccountNumber           string
	ApiAccessToken          string
	ApiRefreshToken         string

	clients *Clients
}

// Configuring Lumen Provider clients
func (c *Config) LumenClients() (*Clients, diag.Diagnostics) {
	if c.clients == nil {
		mClient := client.NewMorpheusClient(c.MorpheusBareMetalApiUrl, c.AuthUrl)
		mClient.SetCredsAndTokens(c.Username, c.Password, c.ApiAccessToken, c.ApiRefreshToken)

		c.clients = &Clients{
			Morpheus:  mClient,
			BareMetal: client.NewBareMetalClient(c.ApigeeBaseURL, c.Username, c.Password, c.AccountNumber),
		}
	}
	return c.clients, nil
}
