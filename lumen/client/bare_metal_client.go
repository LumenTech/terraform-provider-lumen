package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"
	"strings"
	"terraform-provider-lumen/lumen/client/model/bare_metal"
	"time"
)

var retryCount = 5
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
	client.SetRetryCount(retryCount)
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
	resp, err := bm.execute("GET", url, nil, []bare_metal.Location{})
	if err != nil {
		return nil, err
	}

	return resp.Result().(*[]bare_metal.Location), nil
}

func (bm *BareMetalClient) GetConfigurations(locationId string) (*[]bare_metal.Configuration, error) {
	url := fmt.Sprintf("%s/locations/%s/configurations", bm.URL, locationId)
	resp, err := bm.execute("GET", url, nil, []bare_metal.Configuration{})
	if err != nil {
		return nil, err
	}

	return resp.Result().(*[]bare_metal.Configuration), nil
}

func (bm *BareMetalClient) GetNetworkSizes(locationId string) (*[]bare_metal.NetworkSize, error) {
	url := fmt.Sprintf("%s/locations/%s/networkSizes", bm.URL, locationId)
	resp, err := bm.execute("GET", url, nil, []bare_metal.NetworkSize{})
	if err != nil {
		return nil, err
	}

	return resp.Result().(*[]bare_metal.NetworkSize), nil
}

func (bm *BareMetalClient) GetOsImages(locationId string) (*[]bare_metal.OsImage, error) {
	url := fmt.Sprintf("%s/locations/%s/osImages", bm.URL, locationId)
	resp, err := bm.execute("GET", url, nil, []bare_metal.OsImage{})
	if err != nil {
		return nil, err
	}
	osImages := *resp.Result().(*[]bare_metal.OsImage)
	retVal := make([]bare_metal.OsImage, 0)
	for _, osImage := range osImages {
		if osImage.Ready {
			retVal = append(retVal, osImage)
		}
	}
	return &retVal, nil
}

func (bm *BareMetalClient) GetServerByName(name string) (*bare_metal.Server, error) {
	url := fmt.Sprintf("%s/servers", bm.URL)
	resp, err := bm.execute("GET", url, nil, []bare_metal.Server{})
	if err != nil {
		return nil, err
	}

	servers := resp.Result().(*[]bare_metal.Server)
	for _, server := range *servers {
		if server.Name == name {
			return &server, nil
		}
	}
	return nil, nil
}

func (bm *BareMetalClient) GetServer(id string) (*bare_metal.Server, error) {
	url := fmt.Sprintf("%s/servers/%s", bm.URL, id)
	resp, err := bm.execute("GET", url, nil, bare_metal.Server{})
	if err != nil {
		if resp != nil && resp.StatusCode() == 404 {
			return nil, nil
		}
		return nil, err
	}

	return resp.Result().(*bare_metal.Server), nil
}

func (bm *BareMetalClient) ProvisionServer(provisionRequest bare_metal.ServerProvisionRequest) (*bare_metal.Server, error) {
	url := fmt.Sprintf("%s/servers", bm.URL)
	resp, err := bm.execute("POST", url, provisionRequest, bare_metal.Server{})
	if err != nil {
		return nil, err
	}
	return resp.Result().(*bare_metal.Server), nil
}

func (bm *BareMetalClient) UpdateServer(serverId string, request bare_metal.ServerUpdateRequest) (*bare_metal.Server, error) {
	url := fmt.Sprintf("%s/servers/%s", bm.URL, serverId)
	resp, err := bm.execute("PUT", url, request, bare_metal.Server{})
	if err != nil {
		return nil, err
	}
	return resp.Result().(*bare_metal.Server), nil
}

var deletingStatus = []string{"releasing", "networking_removed", "release"}

func (bm *BareMetalClient) DeleteServer(serverId string) (*bare_metal.Server, error) {
	url := fmt.Sprintf("%s/servers/%s", bm.URL, serverId)
	resp, err := bm.execute("DELETE", url, nil, bare_metal.Server{})
	if err != nil && resp.StatusCode() != 404 && resp.StatusCode() != 409 {
		return nil, err
	}

	if resp.StatusCode() == 404 {
		return nil, nil
	}

	if resp.StatusCode() == 409 {
		// if server returns 409 then it is in a transitioning status and could be in the process of deleting
		server, getServerError := bm.GetServer(serverId)
		if getServerError != nil {
			return nil, fmt.Errorf("failed to retrieve server (%s) unable to figure out state", serverId)
		} else if server == nil {
			return nil, nil
		}

		foundDeletingStatus := false
		for _, status := range deletingStatus {
			if strings.ToLower(server.Status) == status {
				foundDeletingStatus = true
			}
		}

		if foundDeletingStatus {
			return server, nil
		}

		return server, fmt.Errorf("failed to delete server (%s) due to pending change", serverId)
	}

	return resp.Result().(*bare_metal.Server), nil
}

func (bm *BareMetalClient) GetNetwork(networkId string) (*bare_metal.Network, error) {
	url := fmt.Sprintf("%s/networks/%s", bm.URL, networkId)
	resp, err := bm.execute("GET", url, nil, bare_metal.Network{})
	if err != nil && resp.StatusCode() != 404 {
		return nil, err
	} else if resp.StatusCode() == 404 {
		return nil, nil
	}

	return resp.Result().(*bare_metal.Network), nil
}

func (bm *BareMetalClient) ProvisionNetwork(provisionRequest bare_metal.NetworkProvisionRequest) (*bare_metal.Network, error) {
	url := fmt.Sprintf("%s/networks", bm.URL)
	resp, err := bm.execute("POST", url, provisionRequest, bare_metal.Network{})
	if err != nil {
		return nil, err
	}
	return resp.Result().(*bare_metal.Network), nil
}

func (bm *BareMetalClient) UpdateNetwork(networkId string, request bare_metal.NetworkUpdateRequest) (*bare_metal.Network, error) {
	url := fmt.Sprintf("%s/networks/%s", bm.URL, networkId)
	resp, err := bm.execute("PUT", url, request, bare_metal.Network{})
	if err != nil {
		return nil, err
	}
	return resp.Result().(*bare_metal.Network), nil
}

func (bm *BareMetalClient) DeleteNetwork(networkId string) (*bare_metal.Network, error) {
	url := fmt.Sprintf("%s/networks/%s", bm.URL, networkId)
	resp, err := bm.execute("DELETE", url, nil, bare_metal.Network{})
	if err != nil {
		return nil, err
	}
	return resp.Result().(*bare_metal.Network), nil
}

func (bm *BareMetalClient) execute(method, url string, body interface{}, result interface{}) (*resty.Response, error) {
	if err := bm.refreshApigeeToken(); err != nil {
		return nil, err
	}

	request := bm.defaultClient.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", bm.ApigeeToken)).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json")

	if body != nil {
		request.SetBody(body)
	}

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
		err = fmt.Errorf("%s (%s) failures reason (%s)", method, url, reason)
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
