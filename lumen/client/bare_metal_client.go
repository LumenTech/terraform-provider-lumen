package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"
	"terraform-provider-lumen/lumen/client/model/bare_metal"
	"time"
)

type BareMetalClient struct {
	BaseURL            string
	ApigeeAuthEndpoint string
	ApigeeUsername     string
	ApigeePassword     string
	ApigeeToken        string
	ExpireTime         int64
	AccountAlias       string
	defaultClient      *resty.Client
}

func NewBareMetalClient(username, password, account string) *BareMetalClient {
	client := resty.New()
	client.SetHeader("User-Agent", "lumen-terraform-plugin v0.5.3")
	client.SetHeader("x-billing-account-number", account)
	return &BareMetalClient{
		BaseURL:            "https://api-dev1.lumen.com/EdgeServices/v2/Compute/bareMetal",
		ApigeeAuthEndpoint: "https://api-dev1.lumen.com/oauth/token",
		ApigeeUsername:     username,
		ApigeePassword:     password,
		AccountAlias:       account,
		defaultClient:      client,
	}
}

func (bm *BareMetalClient) GetLocations() ([]bare_metal.Location, error) {
	resp, err := bm.execute("GET", fmt.Sprintf("%s/locations", bm.BaseURL))
	if err != nil || !resp.IsSuccess() {
		return nil, errors.New("bare metal api failure")
	}

	var locations []bare_metal.Location
	if jsonErr := json.Unmarshal(resp.Body(), &locations); jsonErr != nil {
		return nil, errors.New("unable to parse location response")
	}

	return locations, nil
}

func (bm *BareMetalClient) execute(method, url string) (*resty.Response, error) {
	// TODO: Should this handle some default retry policy
	if err := bm.refreshApigeeToken(); err != nil {
		return nil, err
	}

	request := bm.defaultClient.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", bm.ApigeeToken)).
		SetHeader("Accept", "application/json")

	return request.Execute(method, url)
}

func (bm *BareMetalClient) refreshApigeeToken() error {
	expireTime := time.UnixMilli(bm.ExpireTime - 60000)
	if len(bm.ApigeeToken) == 0 || time.Now().After(expireTime) {
		resp, err := bm.defaultClient.R().
			SetBasicAuth(bm.ApigeeUsername, bm.ApigeePassword).
			SetHeader("Accept", "application/json").
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
		if issuerErr == nil || expiresErr == nil {
			bm.ExpireTime = issueAt + (expiresInSeconds * 1000)
		}
	}

	return nil
}
