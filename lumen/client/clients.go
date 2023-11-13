package client

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type Clients struct {
	BareMetal *BareMetalClient
}

type Config struct {
	ApigeeBaseURL string
	Username      string
	Password      string
	AccountNumber string
	clients       *Clients
}

func (c *Config) LumenClients() (*Clients, diag.Diagnostics) {
	if c.clients == nil {
		c.clients = &Clients{
			BareMetal: NewBareMetalClient(c.ApigeeBaseURL, c.Username, c.Password, c.AccountNumber),
		}
	}
	return c.clients, nil
}
