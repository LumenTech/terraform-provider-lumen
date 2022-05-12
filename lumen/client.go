package lumen

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	ApiUrl    string
	AuthUrl   string
	UserAgent string

	Username        string
	Password        string
	AuthToken       string
	ApiAccessToken  string
	ApiRefreshToken string

	ExpiresIn     string
	AuthTimestamp string
	Scope         string
}

// Setting up new client
func NewClient(apiUrl string, authUrl string) (client *Client) {
	var userAgent = "lumen-terraform-plugin v0.3.5"
	return &Client{
		ApiUrl:    apiUrl,
		AuthUrl:   authUrl,
		UserAgent: userAgent,
	}
}

// Setting up token for Lumen APIGee auth
func (client *Client) SetAuthToken(authToken string) {
	client.AuthToken = authToken
}

// Clearing Lumen APIGee auth
func (client *Client) ClearAuthToken() {
	client.AuthToken = ""
}

// Setting up credentials and api token for lumen API access
func (client *Client) SetCredsAndTokens(
	username string, password string,
	apiAccessToken string, apiRefreshToken string) {
	client.Username = username
	client.Password = password
	client.ApiAccessToken = apiAccessToken
	client.ApiRefreshToken = apiRefreshToken
}

// Executing authentication and generating authentication token
func (client *Client) GenerateAuthToken() error {
	// The transient resty response object
	var authRestyResp *resty.Response
	// potential error to be returned
	var err error

	// constructing request
	var httpMethod string = "POST"
	var url string = client.AuthUrl

	// iniliazing client
	authRestyClient := resty.New()

	// Ignoring ssl cert errors for now
	// TODO: get this into client config
	if strings.HasPrefix(url, "https") {
		authRestyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	// Constructing resty.Request
	authRestyReq := authRestyClient.R()
	authRestyReq.Method = httpMethod

	// Set default Content type
	if authRestyReq.Header["Content-Type"] == nil {
		authRestyReq.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	}
	// Set default Accept header
	if authRestyReq.Header["Accept"] == nil {
		authRestyReq.SetHeader("Accept", "application/json")
	}

	payload := make(map[string]string)

	// Setting username password
	authRestyClient.SetBasicAuth(client.Username, client.Password)
	// Setting 'grant_type' 'client_credentials' in payload
	payload["grant_type"] = "client_credentials"
	authRestyClient.SetQueryParams(payload)

	// Debug - requests
	// log.Printf(fmt.Sprintf("==> Auth API Request: %s %s JSON: %s %s", authRestyReq.Method,
	//	url, authRestyReq.Header, authRestyReq.Body))

	// Make the request
	if httpMethod == "POST" {
		authRestyResp, err = authRestyReq.Post(url)
	} else {
		return errors.New(fmt.Sprintf("Invalid Request.  Unknown HTTP Method: %v", httpMethod))
	}

	// The response object to be returned
	var authResp *AuthResponse
	// Converting resty response into local Response object
	authResp = &AuthResponse{
		Success:    authRestyResp.IsSuccess(),
		StatusCode: authRestyResp.StatusCode(),
		Status:     authRestyResp.Status(),
		Size:       authRestyResp.Size(),
		Body:       authRestyResp.Body(), // byte[]
	}

	// Debug - response
	// log.Printf("==> Auth API Response: [%v] %d %s", authResp.Success, authResp.StatusCode, authResp.Body)
	if err != nil {
		log.Printf("API Error: %v", err)
	}

	// Determining success and setting error accordingly
	if authResp.Success != true {
		err = errors.New(fmt.Sprintf("API returned HTTP %d", authResp.StatusCode))
		return err
	}

	// The http response will be at RestyResponse.RawResponse
	authResp.AuthRestyResponse = authRestyResp

	// Parsing json blob, populates JsonData
	var parsedResult map[string]interface{}
	var jsonParseError error
	jsonParseError = json.Unmarshal([]byte(authResp.Body), &parsedResult)

	authResp.JsonData = parsedResult
	authResp.JsonParseError = jsonParseError

	if jsonParseError != nil {
		log.Printf("Failed to parse JSON result for type %T. Parse Error: %v", authResp.Result, jsonParseError)
	} else {
		log.Printf("Parsed JSON result for type %s, type %s", parsedResult, reflect.TypeOf(parsedResult))
	}

	var val interface{}
	// Client access token
	val = parsedResult["access_token"]
	client.AuthToken = val.(string)

	// Set timestamp and expires_in timing for access_token
	val = parsedResult["issued_at"]
	client.AuthTimestamp = val.(string)
	val = parsedResult["expires_in"]
	client.ExpiresIn = val.(string)
	return err
}

func (client *Client) Execute(req *Request) (*Response, error) {
	// The transient resty response object
	var restyResponse *resty.Response
	// The response object to be returned
	var resp *Response
	// potential error to be returned
	var err error

	// TODO: renew access_token only when it's close to expiring

	// Generating access_token
	err = client.GenerateAuthToken()
	if err != nil {
		log.Printf("Error in generating APIGee auth access_token %s", err)
		return nil, errors.New("Invalid APIGee request")
	}

	// construct the request
	var httpMethod = req.Method
	if httpMethod == "" {
		return nil, errors.New("Invalid Request: Method is required eg. GET,POST,PUT,DELETE")
	}

	// Setting url for API
	var url string = client.ApiUrl + req.Path
	restyClient := resty.New()

	// Ignoring ssl cert errors for now
	// TODO: get this into client config
	if strings.HasPrefix(url, "https") {
		restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	// Constructing resty.Request
	restyReq := restyClient.R()

	// set query params
	restyReq.SetQueryParams(req.QueryParams)

	// Set default headers: application/json
	if req.Headers != nil {
		for k, v := range req.Headers {
			restyReq.SetHeader(k, v)
		}
	}

	// add Authorization Header with our access token
	if restyReq.Header["Authorization"] == nil {
		if client.AuthToken != "" {
			restyReq.SetHeader("Authorization", "Bearer "+client.AuthToken)
		}
	}
	// add Morpheus Token with ApiToken
	if restyReq.Header["morpheusToken"] == nil {
		if client.ApiAccessToken != "" {
			restyReq.SetHeader("morpheusToken", client.ApiAccessToken)
		}
	}

	// set body
	if httpMethod == "POST" || httpMethod == "PUT" {
		// FormData uses application/x-www-form-urlencoded
		if req.FormData != nil {
			log.Printf("REQUEST FORM DATA: ", req.FormData)
			restyReq.SetFormData(req.FormData)
			if restyReq.Header["Content-Type"] == nil {
				restyReq.SetHeader("Content-Type", "application/x-www-form-urlencoded")
			}
		}
		if req.Body != nil {
			log.Printf("REQUEST BODY: ", req.Body)
			restyReq.SetBody(req.Body)
			if restyReq.Header["Content-Type"] == nil {
				restyReq.SetHeader("Content-Type", "application/json")
			}
		}

		// Set default headers: application/json
		if restyReq.Header["Content-Type"] == nil {
			restyReq.SetHeader("Content-Type", "application/json")
		}
	}

	// Set default Accept header
	if restyReq.Header["Accept"] == nil {
		restyReq.SetHeader("Accept", "application/json")
	}

	// Make the request
	if httpMethod == "GET" {
		restyResponse, err = restyReq.Get(url)
	} else if httpMethod == "POST" {
		restyResponse, err = restyReq.Post(url)
	} else if httpMethod == "PUT" {
		restyResponse, err = restyReq.Put(url)
	} else if httpMethod == "DELETE" {
		restyResponse, err = restyReq.Delete(url)
	} else {
		return nil, errors.New(fmt.Sprintf("Invalid Request.  Unknown HTTP Method: %v", httpMethod))
	}

	// Converting resty response into local Response object
	resp = &Response{
		Success:    restyResponse.IsSuccess(),
		StatusCode: restyResponse.StatusCode(),
		Status:     restyResponse.Status(),
		Size:       restyResponse.Size(),
		Body:       restyResponse.Body(), // byte[]
	}

	// Clearing up access token
	client.ClearAuthToken()

	// Determining success and setting error accordingly
	if resp.Success != true {
		err = errors.New(fmt.Sprintf("API returned HTTP %d", resp.StatusCode))
		var standardResult StandardResult
		standardResultParseErr := json.Unmarshal(resp.Body, &standardResult)
		if standardResult.Message != "" {
			err = errors.New(standardResult.Message)
		}
		if standardResultParseErr != nil {
			err = standardResultParseErr
		}
	}

	// The http response will be at RestyResponse.RawResponse
	resp.RestyResponse = restyResponse

	// Parsing json blob, populates JsonData
	var parsedResult interface{}
	jsonError := parseJsonToResult(resp.Body, &parsedResult)
	resp.JsonData = parsedResult
	resp.JsonParseError = jsonError

	// Parsed json into interface{}
	resp.Result = req.Result
	if resp.Result != nil {
		jsonParseResultError := parseJsonToResult(resp.Body, &resp.Result)
		if jsonParseResultError != nil {
			log.Printf("Failed to parse JSON result for type %T. Parse Error: %v", resp.Result, jsonParseResultError)
		} else {
			log.Printf("Parsed JSON result for type %T", resp.Result)
		}
	}

	// Debug - requests
	// log.Printf(fmt.Sprintf("==> Request: %s %s JSON: %s", req.Method, url, req.Body))

	// Debug - response
	// log.Printf("API Response: [%v] %d %s", resp.Success, resp.StatusCode, resp.Body)
	if err != nil {
		log.Printf("API Error: %v", err)
	}
	return resp, err
}

func (client *Client) Get(req *Request) (*Response, error) {
	req.Method = "GET"
	return client.Execute(req)
}

func (client *Client) Post(req *Request) (*Response, error) {
	req.Method = "POST"
	return client.Execute(req)
}

func (client *Client) Put(req *Request) (*Response, error) {
	req.Method = "PUT"
	return client.Execute(req)
}

func (client *Client) Delete(req *Request) (*Response, error) {
	req.Method = "DELETE"
	return client.Execute(req)
}

// helper function: set success accordingly
func IsSuccess(statusCode int) bool {
	var flag bool

	// statusCode should be 2xx
	if statusCode >= 200 && statusCode <= 299 {
		flag = true
	} else {
		flag = false
	}
	return flag
}

// helper function: to parse json
func parseJsonToResult(data []byte, output interface{}) error {
	var err error
	if data != nil {
		err = json.Unmarshal(data, &output)
	}
	return err
}
