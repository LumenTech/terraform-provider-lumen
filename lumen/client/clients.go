package client

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type Clients struct {
	BareMetal *BareMetalClient
}

type Config struct {
	ApigeeBaseURL  string
	ConsumerKey    string
	ConsumerSecret string
	AccountNumber  string
	clients        *Clients
}

func (c *Config) LumenClients() (*Clients, diag.Diagnostics) {
	if c.clients == nil {
		c.clients = &Clients{
			BareMetal: NewBareMetalClient(c.ApigeeBaseURL, c.ConsumerKey, c.ConsumerSecret, c.AccountNumber),
		}
	}
	return c.clients, nil
}
