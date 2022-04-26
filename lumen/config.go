package lumen

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// Provider Configuration
type Config struct {
	ApiUrl          string
	AuthUrl         string
	Username        string
	Password        string
	ApiAccessToken  string
	ApiRefreshToken string

	client *Client
}

// Configuring Lumen Provider client
func (c *Config) LumenClient() (*Client, diag.Diagnostics) {
	if c.client == nil {
		client := NewClient(c.ApiUrl, c.AuthUrl)
		client.SetCredsAndTokens(c.Username, c.Password, c.ApiAccessToken, c.ApiRefreshToken)
		c.client = client
	}
	return c.client, nil
}
