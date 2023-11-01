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

type HttpResponses []HttpResponse

type HttpResponse struct {
	StatusCode int
	Body       interface{}
}

func setupTestServer(t *testing.T, apigeeResponses HttpResponses, apigeeCallCount *int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.RequestURI == "/oauth/token" {
			*apigeeCallCount++
			assert.Equal(t, "POST", req.Method)
			assert.Contains(t, req.Header.Get("Authorization"), "Basic")
			assert.Equal(t, "application/json", req.Header.Get("Accept"))
			assert.Equal(t, "application/x-www-form-urlencoded", req.Header.Get("Content-Type"))
			assert.Equal(t, "lumen-terraform-plugin v0.5.3", req.Header.Get("User-Agent"))

			response := apigeeResponses[0]
			apigeeResponses = apigeeResponses[1:]
			if response.Body != nil {
				_ = json.NewEncoder(w).Encode(response.Body)
			}
			w.WriteHeader(response.StatusCode)
		}
	}))
}

func TestRefreshToken_NotExpired(t *testing.T) {
	apigeeCallCount := 0

	fakeApigeeToken, _ := uuid.GenerateUUID()
	issuedAt := time.Now().UnixMilli()
	expiresIn := int64(1800)

	testServer := setupTestServer(t, HttpResponses{
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"access_token": fakeApigeeToken,
				"issued_at":    fmt.Sprintf("%d", issuedAt),
				"expires_in":   fmt.Sprintf("%d", expiresIn),
			},
		},
	}, &apigeeCallCount)
	defer testServer.Close()

	client := NewBareMetalClient(testServer.URL, "test_user", "test_password", "test_account")

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

	client := NewBareMetalClient(testServer.URL, "test_user", "test_password", "test_account")

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

func TestRefreshToken_RetryableClient(t *testing.T) {
	apigeeCallCount := 0
	responses := make(HttpResponses, retryCount)
	for i := 0; i < retryCount; i++ {
		responses[i] = HttpResponse{
			StatusCode: 500,
		}
	}
	testServer := setupTestServer(t, responses, &apigeeCallCount)
	defer testServer.Close()

	retryWaitTime = 1 * time.Second
	retryMaxWaitTime = 1 * time.Second
	client := NewBareMetalClient(testServer.URL, "test_user", "test_password", "test_account")

	err := client.refreshApigeeToken()
	assert.NotNil(t, err)
	assert.Equal(t, apigeeCallCount, 5)
}
