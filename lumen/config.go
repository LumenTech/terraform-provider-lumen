package lumen

import (
	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type Config struct {
	Url          string
	AccessToken  string
	RefreshToken string
	Username     string
	Password     string
	ClientId     string
	Insecure     bool
	client       *morpheus.Client
	userAgent    string
}

func (c *Config) Client() (*morpheus.Client, diag.Diagnostics) {
	if c.client == nil {
		client := morpheus.NewClient(c.Url)
		if c.AccessToken != "" {
			var expiresIn int64 = 86400 // not used
			client.SetAccessToken(c.AccessToken, c.RefreshToken, expiresIn, "write")
		} else {
			client.SetUsernameAndPassword(c.Username, c.Password)
		}
		c.client = client
	}
	return c.client, nil
}
