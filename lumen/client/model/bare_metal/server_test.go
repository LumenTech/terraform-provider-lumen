package bare_metal

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestServerProvisionRequest_NetworkRequestAndPublicKeyOmitEmpty(t *testing.T) {
	request := ServerProvisionRequest{
		Name:          "name",
		LocationID:    "location",
		Configuration: "large",
		OSImage:       "Ubuntu 22.04",
		NetworkID:     "network-id",
		Credentials: Credentials{
			Username: "test-user",
			Password: "test-password",
		},
	}

	data, err := json.Marshal(request)
	assert.Nil(t, err)
	jsonString := string(data)
	assert.True(t, !strings.Contains(jsonString, "networkRequest"))
	assert.True(t, !strings.Contains(jsonString, "publicKey"))
	assert.JSONEq(t, `
{
	"name":"name",
	"locationId":"location",
	"configuration":"large",
	"osImage":"Ubuntu 22.04",
	"networkId":"network-id",
	"credentials":{
		"username":"test-user",
		"password":"test-password"
	}
}
`, jsonString)
}

func TestServerProvisionRequest_NetworkIdAndPasswordOmitEmpty(t *testing.T) {
	request := ServerProvisionRequest{
		Name:          "name",
		LocationID:    "location",
		Configuration: "large",
		OSImage:       "Ubuntu 22.04",
		NetworkRequest: &NetworkProvisionRequest{
			Name:          "test-net",
			LocationID:    "location",
			NetworkSizeID: "some-uuid",
			NetworkType:   "DUAL_STACK_INTERNET",
		},
		Credentials: Credentials{
			Username:  "test-user",
			PublicKey: "public-key",
		},
		AssignIPV6Address: true,
	}

	data, err := json.Marshal(request)
	assert.Nil(t, err)
	jsonString := string(data)
	assert.True(t, !strings.Contains(jsonString, "networkId"))
	assert.True(t, !strings.Contains(jsonString, "password"))
	assert.JSONEq(t, `
{
	"name":"name",
	"locationId":"location",
	"configuration":"large",
	"osImage":"Ubuntu 22.04",
	"networkRequest": {
		"name": "test-net", 
		"locationId": "location", 
		"networkSizeId":"some-uuid",
		"networkType": "DUAL_STACK_INTERNET"
	},
	"credentials":{
		"username":"test-user",
		"publicKey":"public-key"
	},
	"assignIpv6Address":true
}
`, jsonString)
}
