package client

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRefreshToken_NotExpired(t *testing.T) {
	apigeeCallCount := 0
	issuedAt := time.Now().UnixMilli()
	expiresIn := int64(1800)
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.RequestURI == "/oauth/token" {
			apigeeCallCount++
			assert.Equal(t, "POST", req.Method)
			assert.Contains(t, req.Header.Get("Authorization"), "Basic")
			assert.Equal(t, "application/json", req.Header.Get("Accept"))
			assert.Equal(t, "application/x-www-form-urlencoded", req.Header.Get("Content-Type"))
			assert.Equal(t, "lumen-terraform-plugin v0.5.3", req.Header.Get("User-Agent"))

			fakeApigeeToken, _ := uuid.GenerateUUID()
			if err := json.NewEncoder(w).Encode(map[string]string{
				"access_token": fakeApigeeToken,
				"issued_at":    fmt.Sprintf("%d", issuedAt),
				"expires_in":   fmt.Sprintf("%d", expiresIn),
			}); err == nil {
				w.WriteHeader(200)
			}
		}
	}))
	defer testServer.Close()

	client := NewBareMetalClient("test_user", "test_password", "test_account")
	client.BaseURL = testServer.URL + "/Infrastructure/v1/BMC/"
	client.ApigeeAuthEndpoint = testServer.URL + "/oauth/token"

	err := client.refreshApigeeToken()
	assert.Nil(t, err)
	assert.Equal(t, apigeeCallCount, 1)

	apigeeToken := client.ApigeeToken
	assert.NotEmpty(t, apigeeToken)
	expectedExpireTime := issuedAt + (expiresIn * 1000)
	assert.Equal(t, expectedExpireTime, client.ExpireTime)

	err = client.refreshApigeeToken()
	assert.Nil(t, err)
	assert.Equal(t, client.ApigeeToken, apigeeToken)
	assert.Equal(t, apigeeCallCount, 1)
}

func TestRefreshToken_ExpiredToken(t *testing.T) {
	apigeeCallCount := 0
	offset := time.Minute * 30
	issuedAt := time.Now().UnixMilli() - offset.Milliseconds()
	expiresIn := int64(1800)
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.RequestURI == "/oauth/token" {
			apigeeCallCount++
			assert.Equal(t, "POST", req.Method)
			assert.Contains(t, req.Header.Get("Authorization"), "Basic")
			assert.Equal(t, "application/json", req.Header.Get("Accept"))
			assert.Equal(t, "application/x-www-form-urlencoded", req.Header.Get("Content-Type"))
			assert.Equal(t, "lumen-terraform-plugin v0.5.3", req.Header.Get("User-Agent"))

			fakeApigeeToken, _ := uuid.GenerateUUID()
			if err := json.NewEncoder(w).Encode(map[string]string{
				"access_token": fakeApigeeToken,
				"issued_at":    fmt.Sprintf("%d", issuedAt),
				"expires_in":   fmt.Sprintf("%d", expiresIn),
			}); err == nil {
				w.WriteHeader(200)
			}
		}
	}))
	defer testServer.Close()

	client := NewBareMetalClient("test_user", "test_password", "test_account")
	client.BaseURL = testServer.URL + "/Infrastructure/v1/BMC/"
	client.ApigeeAuthEndpoint = testServer.URL + "/oauth/token"

	err := client.refreshApigeeToken()
	assert.Nil(t, err)
	assert.Equal(t, apigeeCallCount, 1)

	apigeeToken := client.ApigeeToken
	assert.NotEmpty(t, apigeeToken)
	expectedExpireTime := issuedAt + (expiresIn * 1000)
	assert.Equal(t, expectedExpireTime, client.ExpireTime)

	err = client.refreshApigeeToken()
	assert.Nil(t, err)
	assert.NotEqual(t, client.ApigeeToken, apigeeToken)
	assert.Equal(t, apigeeCallCount, 2)
}
