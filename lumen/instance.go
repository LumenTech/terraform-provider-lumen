package lumen

import (
	"fmt"
)

var (
	// InstancesPath is the API endpoint for instances
	InstancesPath = "/api/instances"
)

// Instance structures for use in request and response payloads
type Instance struct {
	ID           int64                  `json:"id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	InstanceType map[string]interface{} `json:"instanceType"`
	Layout       map[string]interface{} `json:"layout"`
	Group        map[string]interface{} `json:"group"`
	Cloud        map[string]interface{} `json:"cloud"`
	Environment  string                 `json:"instanceContext"`
	Plan         InstancePlan           `json:"plan"`
	Config       map[string]interface{} `json:"config"`
	Labels       []string               `json:"labels"`
	Version      string                 `json:"instanceVersion"`
	Status       string                 `json:"status"`

	Volumes              *[]map[string]interface{} `json:"volumes"`
	Interfaces           *[]map[string]interface{} `json:"interfaces"`
	Controllers          *[]map[string]interface{} `json:"controllers"`
	Tags                 *[]map[string]interface{} `json:"tags"`
	Metadata             *[]map[string]interface{} `json:"metadata"`
	EnvironmentVariables *[]map[string]interface{} `json:"evars"`
}

type StandardResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
}

type InstancePlan struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type MetaResult struct {
	Total  int64       `json:"total"`
	Size   int64       `json:"size"`
	Max    interface{} `json:"max"`
	Offset int64       `json:"offset"`
}

type ListInstancesResult struct {
	Instances *[]Instance `json:"instances"`
	Meta      *MetaResult `json:"meta"`
}

type GetInstanceResult struct {
	Instance *Instance `json:"instance"`
}

type CreateInstanceResult struct {
	Success  bool              `json:"success"`
	Message  string            `json:"msg"`
	Errors   map[string]string `json:"errors"`
	Instance *Instance         `json:"instance"`
}

type UpdateInstanceResult struct {
	CreateInstanceResult
}

type DeleteInstanceResult struct {
	StandardResult
}

// API endpoints
func (client *Client) ListInstances(req *Request) (*Response, error) {
	// List instances
	return client.Execute(&Request{
		Method:      "GET",
		Path:        InstancesPath,
		QueryParams: req.QueryParams,
		Result:      &ListInstancesResult{},
	})
}

func (client *Client) GetInstance(id int64, req *Request) (*Response, error) {
	// Get instance details based on instance id
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", InstancesPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetInstanceResult{},
	})
}

func (client *Client) CreateInstance(req *Request) (*Response, error) {
	// Create instance
	return client.Execute(&Request{
		Method:      "POST",
		Path:        InstancesPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateInstanceResult{},
	})
}

func (client *Client) UpdateInstance(id int64, req *Request) (*Response, error) {
	// Update instance
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", InstancesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateInstanceResult{},
	})
}

func (client *Client) DeleteInstance(id int64, req *Request) (*Response, error) {
	// Delete instance
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", InstancesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteInstanceResult{},
	})
}

func (client *Client) FindInstanceByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListInstances(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListInstancesResult)
	instanceCount := len(*listResult.Instances)
	if instanceCount != 1 {
		return resp, fmt.Errorf("found %d Instances for %v", instanceCount, name)
	}
	firstRecord := (*listResult.Instances)[0]
	instanceId := firstRecord.ID
	return client.GetInstance(instanceId, &Request{})
}
