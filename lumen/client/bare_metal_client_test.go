package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type HttpResponses []HttpResponse

type HttpResponse struct {
	StatusCode int
	Body       interface{}
}

func setupTestServerWithDefaultApigeeResponse(t *testing.T, apiResponse HttpResponses) (*httptest.Server, *int) {
	apigeeCallCount := 0
	return setupTestServer(t, HttpResponses{
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"access_token": fmt.Sprintf("%d", time.Now().Nanosecond()),
				"issued_at":    fmt.Sprintf("%d", time.Now().UnixMilli()),
				"expires_in":   fmt.Sprintf("%d", int64(1800)),
			},
		},
	}, &apigeeCallCount, apiResponse), &apigeeCallCount
}

func setupTestServer(t *testing.T, apigeeResponses HttpResponses, apigeeCallCount *int, apiResponses HttpResponses) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var response HttpResponse
		if req.RequestURI == "/oauth/token" || req.RequestURI == "/oauth/v2/token" {
			*apigeeCallCount++
			assert.Equal(t, "POST", req.Method)
			assert.Contains(t, req.Header.Get("Authorization"), "Basic")
			assert.Equal(t, "application/json", req.Header.Get("Accept"))
			assert.Equal(t, "application/x-www-form-urlencoded", req.Header.Get("Content-Type"))
			assert.Equal(t, "lumen-terraform-plugin v2.5.0", req.Header.Get("User-Agent"))

			response = apigeeResponses[0]
			if len(apigeeResponses) > 1 {
				apigeeResponses = apigeeResponses[1:]
			}
		} else {
			// Handle other requests
			assert.Contains(t, req.Header.Get("Authorization"), "Bearer ")
			assert.Equal(t, "application/json", req.Header.Get("Accept"))
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
			assert.NotEmpty(t, req.Header.Get("x-billing-account-number"))

			response = apiResponses[0]
			if len(apiResponses) > 1 {
				apiResponses = apiResponses[1:]
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.StatusCode)
		if response.Body != nil {
			_ = json.NewEncoder(w).Encode(response.Body)
		}
	}))
}

func TestRefreshToken_NotExpired(t *testing.T) {
	testServer, apigeeCallCount := setupTestServerWithDefaultApigeeResponse(t, nil)
	defer testServer.Close()

	client := NewBareMetalClient(testServer.URL, "test_user", "test_password", "test_account")

	err := client.refreshApigeeToken()
	assert.Nil(t, err)
	assert.Equal(t, *apigeeCallCount, 1)

	apigeeToken := client.ApigeeToken
	expireTime := client.ExpireTime
	assert.NotEmpty(t, apigeeToken)
	assert.NotEmpty(t, expireTime)

	err = client.refreshApigeeToken()
	assert.Nil(t, err)
	assert.Equal(t, *apigeeCallCount, 1)
	assert.Equal(t, apigeeToken, client.ApigeeToken)
	assert.Equal(t, expireTime, client.ExpireTime)
}

func TestRefreshToken_ExpiredToken(t *testing.T) {
	apigeeCallCount := 0
	offset := time.Minute * 30
	issuedAt := time.Now().UnixMilli() - offset.Milliseconds()
	expiresIn := int64(1800)

	testServer := setupTestServer(t, HttpResponses{
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"access_token": "test-token-1",
				"issued_at":    fmt.Sprintf("%d", issuedAt),
				"expires_in":   fmt.Sprintf("%d", expiresIn),
			},
		},
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"access_token": "test-token-2",
				"issued_at":    fmt.Sprintf("%d", time.Now().UnixMilli()),
				"expires_in":   fmt.Sprintf("%d", expiresIn),
			},
		},
	}, &apigeeCallCount, nil)
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
	assert.NotEqual(t, apigeeToken, client.ApigeeToken)
	assert.Equal(t, apigeeCallCount, 2)
}

func TestRefreshToken_RetryableClient(t *testing.T) {
	apigeeCallCount := 0
	testServer := setupTestServer(t, HttpResponses{
		{
			StatusCode: 500,
		},
	}, &apigeeCallCount, nil)
	defer testServer.Close()

	retryWaitTime = 1 * time.Second
	retryMaxWaitTime = 1 * time.Second
	client := NewBareMetalClient(testServer.URL, "test_user", "test_password", "test_account")

	err := client.refreshApigeeToken()
	assert.NotNil(t, err)
	assert.Equal(t, apigeeCallCount, 10)
}

func TestGetLocations(t *testing.T) {
	responseBody := []map[string]interface{}{
		{
			"id":     "test-id",
			"name":   "Test Site",
			"status": "Test Status",
			"region": "NA",
		},
	}
	apiResponses := HttpResponses{
		{
			StatusCode: 200,
			Body:       responseBody,
		},
	}
	testServer, apigeeCallCount := setupTestServerWithDefaultApigeeResponse(t, apiResponses)
	defer testServer.Close()

	client := NewBareMetalClient(testServer.URL, "test_user", "test_password", "test_account")

	locations, err := client.GetLocations()
	deref := *locations
	assert.Nil(t, err)
	assert.Equal(t, 1, *apigeeCallCount)
	assert.Equal(t, 1, len(deref))

	location := deref[0]
	assert.Equal(t, responseBody[0]["id"], location.ID)
	assert.Equal(t, responseBody[0]["name"], location.Name)
	assert.Equal(t, responseBody[0]["status"], location.Status)
	assert.Equal(t, responseBody[0]["region"], location.Region)
}

func TestGetOsImages(t *testing.T) {
	responseBody := []map[string]interface{}{
		{
			"name":  "Ubuntu 20.04",
			"ready": true,
			"price": map[string]interface{}{
				"amount": 45.00,
				"period": "MONTHLY",
			},
		},
		{
			"name":  "Ubuntu 21.04",
			"ready": false,
			"price": map[string]interface{}{
				"amount": 50.00,
				"period": "MONTHLY",
			},
		},
	}
	apiResponses := HttpResponses{
		{
			StatusCode: 200,
			Body:       responseBody,
		},
	}
	testServer, apigeeCallCount := setupTestServerWithDefaultApigeeResponse(t, apiResponses)
	defer testServer.Close()

	client := NewBareMetalClient(testServer.URL, "test_user", "test_password", "test_account")

	osImages, err := client.GetOsImages("testLocation")
	deref := *osImages
	assert.Nil(t, err)
	assert.Equal(t, 1, *apigeeCallCount)
	assert.Equal(t, 1, len(deref))

	osImage := deref[0]
	assert.Equal(t, responseBody[0]["name"], osImage.Name)
	assert.Equal(t, responseBody[0]["ready"], osImage.Ready)
	assert.Equal(t, "$45.00/MONTHLY", osImage.Price.String())
}

func TestGetServer_NotFound(t *testing.T) {
	apiResponses := HttpResponses{
		{
			StatusCode: 404,
		},
	}

	testServer, _ := setupTestServerWithDefaultApigeeResponse(t, apiResponses)
	defer testServer.Close()

	client := NewBareMetalClient(testServer.URL, "test_user", "test_password", "test_account")

	server, err := client.GetServer("test-id")
	assert.Nil(t, err)
	assert.Nil(t, server)
}

func TestGetServer_ServerError(t *testing.T) {
	retryWaitTime = 1 * time.Second
	retryMaxWaitTime = 1 * time.Second

	apiResponses := HttpResponses{
		{
			StatusCode: 500,
		},
	}

	testServer, _ := setupTestServerWithDefaultApigeeResponse(t, apiResponses)
	defer testServer.Close()

	client := NewBareMetalClient(testServer.URL, "test_user", "test_password", "test_account")

	server, err := client.GetServer("test-id")
	assert.NotNil(t, err)
	assert.Nil(t, server)
}

func TestDeleteServer_NotFoundResponse(t *testing.T) {
	responses := HttpResponses{
		{
			StatusCode: 404,
		},
	}
	testServer, _ := setupTestServerWithDefaultApigeeResponse(t, responses)
	defer testServer.Close()

	client := NewBareMetalClient(testServer.URL, "test_user", "test_password", "test_account")

	server, err := client.DeleteServer("test-id")
	assert.Nil(t, err)
	assert.Nil(t, server)
}

func TestDeleteServer_ConflictResponse(t *testing.T) {
	responses := HttpResponses{
		{
			StatusCode: 409,
		},
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"status": "releasing",
			},
		},
	}
	testServer, _ := setupTestServerWithDefaultApigeeResponse(t, responses)
	defer testServer.Close()

	client := NewBareMetalClient(testServer.URL, "test_user", "test_password", "test_account")

	server, err := client.DeleteServer("test-id")
	assert.Nil(t, err)
	assert.NotNil(t, server)
	assert.Equal(t, "releasing", server.Status)
}

func TestDeleteServer_ConflictButServerNotFoundAfterPolling(t *testing.T) {
	responses := HttpResponses{
		{
			StatusCode: 409,
		},
		{
			StatusCode: 404,
		},
	}
	testServer, _ := setupTestServerWithDefaultApigeeResponse(t, responses)
	defer testServer.Close()

	client := NewBareMetalClient(testServer.URL, "test_user", "test_password", "test_account")

	server, err := client.DeleteServer("test-id")
	assert.Nil(t, err)
	assert.Nil(t, server)
}

func TestDeleteServer_ConflictButServerInFailedStatus(t *testing.T) {
	responses := HttpResponses{
		{
			StatusCode: 409,
		},
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"status": "failed",
			},
		},
	}
	testServer, _ := setupTestServerWithDefaultApigeeResponse(t, responses)
	defer testServer.Close()

	client := NewBareMetalClient(testServer.URL, "test_user", "test_password", "test_account")

	server, err := client.DeleteServer("test-id")
	assert.NotNil(t, err)
	assert.NotNil(t, server)
	assert.Equal(t, "failed", server.Status)
}
