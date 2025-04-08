package client

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"strconv"
	"strings"
	"terraform-provider-lumen/lumen/helper"
	"time"

	"terraform-provider-lumen/lumen/client/model/bare_metal"

	"github.com/go-resty/resty/v2"
)

var retryCount = 5
var retryWaitTime = 1 * time.Second
var retryMaxWaitTime = 30 * time.Second

type BareMetalClient struct {
	URL                  string
	ApigeeAuthEndpoint   string
	ApigeeAuthV2Endpoint string
	V2AuthSuccess        bool
	ApigeeConsumerKey    string
	ApigeeConsumerSecret string
	ApigeeToken          string
	ExpireTime           int64
	AccountNumber        string
	defaultClient        *resty.Client
}

func NewBareMetalClient(apigeeBaseURL, consumerKey, consumerSecret, accountNumber string) *BareMetalClient {
	client := resty.New()
	client.SetHeader("User-Agent", "lumen-terraform-plugin v2.5.0")
	client.SetHeader("x-billing-account-number", accountNumber)
	client.SetRetryCount(retryCount)
	client.SetRetryWaitTime(retryWaitTime)
	client.SetRetryMaxWaitTime(retryMaxWaitTime)
	client.AddRetryCondition(func(response *resty.Response, err error) bool {
		return err != nil || response.StatusCode() == 429 || response.StatusCode() >= 500
	})
	return &BareMetalClient{
		URL:                  fmt.Sprintf("%s/EdgeServices/v2/Compute/bareMetal", apigeeBaseURL),
		ApigeeAuthEndpoint:   fmt.Sprintf("%s/oauth/token", apigeeBaseURL),
		ApigeeAuthV2Endpoint: fmt.Sprintf("%s/oauth/v2/token", apigeeBaseURL),
		V2AuthSuccess:        true,
		ApigeeConsumerKey:    consumerKey,
		ApigeeConsumerSecret: consumerSecret,
		AccountNumber:        accountNumber,
		defaultClient:        client,
	}
}

func (bm *BareMetalClient) GetLocations() (*[]bare_metal.Location, diag.Diagnostics) {
	url := fmt.Sprintf("%s/locations", bm.URL)
	resp, diagnostics := bm.execute("GET", url, nil, []bare_metal.Location{})
	if diagnostics.HasError() {
		return nil, diagnostics
	}

	return resp.Result().(*[]bare_metal.Location), diagnostics
}

func (bm *BareMetalClient) GetConfigurations(locationId string) (*[]bare_metal.Configuration, diag.Diagnostics) {
	url := fmt.Sprintf("%s/locations/%s/configurations", bm.URL, locationId)
	resp, diagnostics := bm.execute("GET", url, nil, []bare_metal.Configuration{})
	if diagnostics.HasError() {
		return nil, diagnostics
	}

	return resp.Result().(*[]bare_metal.Configuration), diagnostics
}

func (bm *BareMetalClient) GetNetworkSizes(locationId string) (*[]bare_metal.NetworkSize, diag.Diagnostics) {
	url := fmt.Sprintf("%s/locations/%s/networkSizes", bm.URL, locationId)
	resp, diagnostics := bm.execute("GET", url, nil, []bare_metal.NetworkSize{})
	if diagnostics.HasError() {
		return nil, diagnostics
	}

	return resp.Result().(*[]bare_metal.NetworkSize), diagnostics
}

func (bm *BareMetalClient) GetOsImages(locationId string) (*[]bare_metal.OsImage, diag.Diagnostics) {
	url := fmt.Sprintf("%s/locations/%s/osImages", bm.URL, locationId)
	resp, diagnostics := bm.execute("GET", url, nil, []bare_metal.OsImage{})
	if diagnostics.HasError() {
		return nil, diagnostics
	}

	osImages := *resp.Result().(*[]bare_metal.OsImage)
	retVal := make([]bare_metal.OsImage, 0)
	for _, osImage := range osImages {
		if osImage.Ready {
			retVal = append(retVal, osImage)
		}
	}
	return &retVal, diagnostics
}

func (bm *BareMetalClient) GetServerByName(name string) (*bare_metal.Server, diag.Diagnostics) {
	url := fmt.Sprintf("%s/servers", bm.URL)
	resp, diagnostics := bm.execute("GET", url, nil, []bare_metal.Server{})
	if !diagnostics.HasError() {
		servers := resp.Result().(*[]bare_metal.Server)
		for _, server := range *servers {
			if server.Name == name {
				return &server, diagnostics
			}
		}
	}

	return nil, diagnostics
}

func (bm *BareMetalClient) GetServer(id string) (*bare_metal.Server, diag.Diagnostics) {
	url := fmt.Sprintf("%s/servers/%s", bm.URL, id)
	resp, diagnostics := bm.execute("GET", url, nil, bare_metal.Server{})
	if diagnostics.HasError() {
		if resp != nil && resp.StatusCode() == 404 {
			return nil, helper.ExtractDiagnosticWarnings(diagnostics)
		}

		return nil, diagnostics
	}

	return resp.Result().(*bare_metal.Server), diagnostics
}

func (bm *BareMetalClient) ProvisionServer(provisionRequest bare_metal.ServerProvisionRequest) (*bare_metal.Server, diag.Diagnostics) {
	url := fmt.Sprintf("%s/servers", bm.URL)
	resp, diagnostics := bm.execute("POST", url, provisionRequest, bare_metal.Server{})
	if diagnostics.HasError() {
		return nil, diagnostics
	}

	return resp.Result().(*bare_metal.Server), diagnostics
}

func (bm *BareMetalClient) UpdateServer(serverId string, request bare_metal.ServerUpdateRequest) (*bare_metal.Server, diag.Diagnostics) {
	url := fmt.Sprintf("%s/servers/%s", bm.URL, serverId)
	resp, diagnostics := bm.execute("PUT", url, request, bare_metal.Server{})
	if diagnostics.HasError() {
		return nil, diagnostics
	}
	return resp.Result().(*bare_metal.Server), diagnostics
}

func (bm *BareMetalClient) AttachNetwork(serverId string, request bare_metal.AddNetworkRequest) (*bare_metal.Server, diag.Diagnostics) {
	url := fmt.Sprintf("%s/servers/%s/networks", bm.URL, serverId)
	resp, diagnostics := bm.execute("POST", url, request, bare_metal.Server{})
	if diagnostics.HasError() {
		return nil, diagnostics
	}
	return resp.Result().(*bare_metal.Server), diagnostics
}

func (bm *BareMetalClient) RemoveNetwork(serverId, networkId string) (*bare_metal.Server, diag.Diagnostics) {
	url := fmt.Sprintf("%s/servers/%s/networks/%s", bm.URL, serverId, networkId)
	resp, diagnostics := bm.execute("DELETE", url, nil, bare_metal.Server{})
	if diagnostics.HasError() {
		if resp != nil && resp.StatusCode() == 404 {
			return nil, helper.ExtractDiagnosticWarnings(diagnostics)
		}
		return nil, diagnostics
	}
	return resp.Result().(*bare_metal.Server), diagnostics
}

var deletingStatus = []string{"releasing", "billing_deactivated", "networking_removed", "released"}

func (bm *BareMetalClient) DeleteServer(serverId string) (*bare_metal.Server, diag.Diagnostics) {
	url := fmt.Sprintf("%s/servers/%s", bm.URL, serverId)
	resp, diagnostics := bm.execute("DELETE", url, nil, bare_metal.Server{})
	if diagnostics.HasError() && resp.StatusCode() != 404 && resp.StatusCode() != 409 {
		return nil, diagnostics
	}

	deleteWarnings := helper.ExtractDiagnosticWarnings(diagnostics)
	if resp.StatusCode() == 404 {
		return nil, deleteWarnings
	}

	if resp.StatusCode() == 409 {
		// if server returns 409 then it is in a transitioning status and could be in the process of deleting
		server, getServerDiagnostics := bm.GetServer(serverId)
		if getServerDiagnostics.HasError() {
			return nil, append(deleteWarnings, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("failed to retrieve server (%s) unable to figure out state", serverId),
			})
		} else if server == nil {
			return nil, deleteWarnings
		}

		foundDeletingStatus := false
		for _, status := range deletingStatus {
			if strings.ToLower(server.Status) == status {
				foundDeletingStatus = true
			}
		}

		if foundDeletingStatus {
			return server, deleteWarnings
		}

		return server, append(deleteWarnings, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to delete server (%s) due to pending change", serverId),
		})
	}

	return resp.Result().(*bare_metal.Server), diagnostics
}

func (bm *BareMetalClient) GetNetwork(networkId string) (*bare_metal.Network, diag.Diagnostics) {
	url := fmt.Sprintf("%s/networks/%s", bm.URL, networkId)
	resp, diagnostics := bm.execute("GET", url, nil, bare_metal.Network{})
	if diagnostics.HasError() {
		if resp != nil && resp.StatusCode() == 404 {
			return nil, helper.ExtractDiagnosticWarnings(diagnostics)
		}
		return nil, diagnostics
	}

	return resp.Result().(*bare_metal.Network), diagnostics
}

func (bm *BareMetalClient) ProvisionNetwork(provisionRequest bare_metal.NetworkProvisionRequest) (*bare_metal.Network, diag.Diagnostics) {
	url := fmt.Sprintf("%s/networks", bm.URL)
	resp, diagnostics := bm.execute("POST", url, provisionRequest, bare_metal.Network{})
	if diagnostics.HasError() {
		return nil, diagnostics
	}
	return resp.Result().(*bare_metal.Network), diagnostics
}

func (bm *BareMetalClient) UpdateNetwork(networkId string, request bare_metal.NetworkUpdateRequest) (*bare_metal.Network, diag.Diagnostics) {
	url := fmt.Sprintf("%s/networks/%s", bm.URL, networkId)
	resp, diagnostics := bm.execute("PUT", url, request, bare_metal.Network{})
	if diagnostics.HasError() {
		return nil, diagnostics
	}
	return resp.Result().(*bare_metal.Network), diagnostics
}

func (bm *BareMetalClient) DeleteNetwork(networkId string) (*bare_metal.Network, diag.Diagnostics) {
	url := fmt.Sprintf("%s/networks/%s", bm.URL, networkId)
	resp, diagnostics := bm.execute("DELETE", url, nil, bare_metal.Network{})
	if diagnostics.HasError() {
		return nil, diagnostics
	}
	return resp.Result().(*bare_metal.Network), diagnostics
}

func (bm *BareMetalClient) execute(method, url string, body interface{}, result interface{}) (*resty.Response, diag.Diagnostics) {
	if authErr := bm.refreshApigeeToken(); authErr != nil {
		return nil, authErr
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
	diagnostics := diag.Diagnostics{}
	if err != nil || !resp.IsSuccess() {
		var reason string
		if err != nil {
			reason = err.Error()
		} else {
			reason = fmt.Sprintf("%s - %s", resp.Status(), resp.String())
		}

		diagnostics = append(diagnostics, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("%s (%s) failures reason (%s)", method, url, reason),
		})
	}

	if resp != nil {
		deprecation := resp.Header().Get("Deprecation")
		if len(deprecation) > 0 && strings.ToLower(deprecation) == "true" {
			sunset := resp.Header().Get("Sunset")
			summary := fmt.Sprintf("Deprecation: API Endpoint (%s) is being deprecated the sunset date is (%s).", url, sunset)
			link := resp.Header().Get("Link")
			if len(link) > 0 {
				summary += fmt.Sprintf(" For more information see link (%s).", link)
			}

			diagnostics = append(diagnostics, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  summary,
			})
		}
	}

	return resp, diagnostics
}

func (bm *BareMetalClient) refreshApigeeToken() diag.Diagnostics {
	expireTime := time.UnixMilli(bm.ExpireTime - 60000)
	if len(bm.ApigeeToken) == 0 || time.Now().After(expireTime) {
		authEndpoint := bm.ApigeeAuthV2Endpoint
		if !bm.V2AuthSuccess {
			authEndpoint = bm.ApigeeAuthEndpoint
		}

		resp, err := bm.defaultClient.R().
			SetBasicAuth(bm.ApigeeConsumerKey, bm.ApigeeConsumerSecret).
			SetHeader("Accept", "application/json").
			SetFormData(map[string]string{
				"grant_type": "client_credentials",
			}).Post(authEndpoint)

		if err != nil || !resp.IsSuccess() {
			bm.V2AuthSuccess = false
			resp, err = bm.defaultClient.R().
				SetBasicAuth(bm.ApigeeConsumerKey, bm.ApigeeConsumerSecret).
				SetHeader("Accept", "application/json").
				SetFormData(map[string]string{
					"grant_type": "client_credentials",
				}).Post(bm.ApigeeAuthEndpoint)
			if err != nil || !resp.IsSuccess() {
				return diag.Errorf("authentication failure")
			}
		}

		var data map[string]interface{}
		if jsonErr := json.Unmarshal(resp.Body(), &data); jsonErr != nil {
			return diag.Errorf("unable to parse response from authentication endpoint: %s", jsonErr.Error())
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
