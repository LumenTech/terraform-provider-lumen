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

var retryWaitTime = 1 * time.Second
var retryMaxWaitTime = 30 * time.Second

type BareMetalClient struct {
	URL                string
	ApigeeAuthEndpoint string
	ApigeeUsername     string
	ApigeePassword     string
	ApigeeToken        string
	ExpireTime         int64
	AccountNumber      string
	defaultClient      *resty.Client
}

func NewBareMetalClient(apigeeBaseURL, username, password, accountNumber string) *BareMetalClient {
	client := resty.New()
	client.SetHeader("User-Agent", "lumen-terraform-plugin v0.5.3")
	client.SetHeader("x-billing-account-number", accountNumber)
	client.SetRetryCount(5)
	client.SetRetryWaitTime(retryWaitTime)
	client.SetRetryMaxWaitTime(retryMaxWaitTime)
	client.AddRetryCondition(func(response *resty.Response, err error) bool {
		return err != nil || response.StatusCode() == 429 || response.StatusCode() >= 500
	})
	return &BareMetalClient{
		URL:                fmt.Sprintf("%s/EdgeServices/v2/Compute/bareMetal", apigeeBaseURL),
		ApigeeAuthEndpoint: fmt.Sprintf("%s/oauth/token", apigeeBaseURL),
		ApigeeUsername:     username,
		ApigeePassword:     password,
		AccountNumber:      accountNumber,
		defaultClient:      client,
	}
}

func (bm *BareMetalClient) GetLocations() (*[]bare_metal.Location, error) {
	url := fmt.Sprintf("%s/locations", bm.URL)
	resp, err := bm.execute("GET", url, []bare_metal.Location{})
	if err != nil {
		return nil, err
	}

	return resp.Result().(*[]bare_metal.Location), nil
}

func (bm *BareMetalClient) GetConfigurations(locationId string) (*[]bare_metal.Configuration, error) {
	url := fmt.Sprintf("%s/locations/%s/configurations", bm.URL, locationId)
	resp, err := bm.execute("GET", url, []bare_metal.Configuration{})
	if err != nil {
		return nil, err
	}

	return resp.Result().(*[]bare_metal.Configuration), nil
}

func (bm *BareMetalClient) GetNetworkSizes(locationId string) (*[]bare_metal.NetworkSize, error) {
	url := fmt.Sprintf("%s/locations/%s/networkSizes", bm.URL, locationId)
	resp, err := bm.execute("GET", url, []bare_metal.NetworkSize{})
	if err != nil {
		return nil, err
	}

	return resp.Result().(*[]bare_metal.NetworkSize), nil
}

func (bm *BareMetalClient) execute(method, url string, result interface{}) (*resty.Response, error) {
	if err := bm.refreshApigeeToken(); err != nil {
		return nil, err
	}

	request := bm.defaultClient.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", bm.ApigeeToken)).
		SetHeader("Accept", "application/json")

	if result != nil {
		request.SetResult(result)
	}

	resp, err := request.Execute(method, url)
	if err != nil || !resp.IsSuccess() {
		var reason string
		if err != nil {
			reason = err.Error()
		} else {
			reason = resp.Status()
		}

		return nil, fmt.Errorf("%s (%s) failures reason (%s)", method, url, reason)
	}

	return resp, err
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
