package lumen

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

// Response parameters from Auth API call
type AuthResponse struct {
	AuthRestyResponse *resty.Response
	Success           bool
	StatusCode        int
	Status            string
	Body              []byte
	Size              int64

	// This holds the parsed JSON for convenience
	JsonData interface{}
	// This holds any error encountering JsonData
	JsonParseError error
	Result         interface{}
}

// Response parameters from API call
type Response struct {
	RestyResponse *resty.Response
	Success       bool
	StatusCode    int
	Status        string
	Body          []byte
	Size          int64

	// This holds the parsed JSON for convenience
	JsonData interface{}
	// This holds any error encountering JsonData
	JsonParseError error
	Result         interface{}
}

func (resp *Response) String() string {
	return fmt.Sprintf("Response HTTP: %v Success: %v  Size: %dB Body: %s",
		resp.Status, resp.Success, resp.Size, resp.Body)
}

// Setting API response interface to prevent typecasting.
type APIResponse interface {
	// whether the request was successful (200 OK) or not.
	Success() bool
	// HTTP status code eg. 2xx
	StatusCode() int
	// HTTP status message .eg "OK"
	Status() string
	// response body byte array
	Body() []byte
	// Response Size
	Size() int64

	// This holds the parsed JSON for convenience
	JsonData() interface{}
	// This holds any error encountering JsonData
	JsonParseError() error

	// the parsed result, in the specified type.
	Result() interface{}
}
