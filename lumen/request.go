package lumen

import (
	"fmt"
)

// Request parameters for API calls
type Request struct {
	Method      string
	Path        string
	QueryParams map[string]string
	Headers     map[string]string
	Body        map[string]interface{}
	FormData    map[string]string
	Timeout     int

	Result interface{}
}

func (req *Request) String() string {
	return fmt.Sprintf("Request Method: %v Path: %v QueryParams: %v Body: %s",
		req.Method, req.Path, req.QueryParams, req.Body)
}
