package client

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRefreshToken(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.RequestURI == "/oauth/token" {
			assert.Equal(t, "POST", req.Method)
			assert.Contains(t, req.Header.Get("Authorization"), "Basic")
			assert.Equal(t, "application/json", req.Header.Get("Accept"))
			assert.Equal(t, "application/x-www-form-urlencoded", req.Header.Get("Content-Type"))
			assert.Equal(t, "lumen-terraform-plugin v0.5.3", req.Header.Get("User-Agent"))
			if err := json.NewEncoder(w).Encode(map[string]string{
				"access_token": "token",
				"issued_at":    "1698412954953",
				"expires_in":   "1800",
			}); err == nil {
				w.WriteHeader(200)
			}
		}
	}))
	defer testServer.Close()

	client := NewBareMetalClient("test_user", "test_password")
	client.BaseURL = testServer.URL + "/Infrastructure/v1/BMC/"
	client.ApigeeAuthEndpoint = testServer.URL + "/oauth/token"

	err := client.refreshApigeeToken()
	assert.Nil(t, err)

	assert.NotEmpty(t, client.ApigeeToken)
	expectedExpireTime := int64(1698412954953 + (1800 * 1000))
	assert.Equal(t, expectedExpireTime, client.ExpireTime)
}
