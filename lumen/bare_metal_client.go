package lumen

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"
	"time"
)

type BareMetalClient struct {
	BaseURL            string
	ApigeeAuthEndpoint string
	ApigeeUsername     string
	ApigeePassword     string
	ApigeeToken        string
	ExpireTime         int64
	DefaultClient      *resty.Client
}

func NewBareMetalClient(username string, password string) *BareMetalClient {
	client := resty.New()
	client.SetHeader("User-Agent", "lumen-terraform-plugin v0.5.2")
	return &BareMetalClient{
		BaseURL:            "https://api-dev1.lumen.com/Infrastructure/v1/BMC/",
		ApigeeAuthEndpoint: "https://api-dev1.lumen.com/oauth/token",
		ApigeeUsername:     username,
		ApigeePassword:     password,
		DefaultClient:      client,
	}
}

func (bm *BareMetalClient) refreshApigeeToken() error {
	expireTime := time.Unix(bm.ExpireTime-60000, 0)
	if len(bm.ApigeeToken) == 0 || time.Now().After(expireTime) {
		client := resty.New()
		resp, err := client.R().
			SetBasicAuth(bm.ApigeeUsername, bm.ApigeePassword).
			SetFormData(map[string]string{
				"grant_type": "client_credentials",
			}).Post(bm.ApigeeAuthEndpoint)

		if err != nil || !resp.IsSuccess() {
			return errors.New("apigee authentication failure")
		}

		var data map[string]interface{}
		if jsonErr := json.Unmarshal(resp.Body(), &data); jsonErr != nil {
			return fmt.Errorf("unable to parse apigee response: %s", jsonErr)
		}

		bm.ApigeeToken = data["access_token"].(string)
		issueAt, issuerErr := strconv.ParseInt(data["issued_at"].(string), 10, 64)
		expiresInSeconds, expiresErr := strconv.ParseInt(data["expires_in"].(string), 10, 64)
		if issuerErr != nil || expiresErr != nil {
			bm.ExpireTime = issueAt + (expiresInSeconds * 1000)
		}
		// TODO: I'm pretty sure int64 defaults to 0 if not set so if I can't calculate I should still
		// be able to refresh the token - Test
	}

	return nil
}
