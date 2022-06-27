package lumen

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceNetworkInstance() *schema.Resource {
	/* Return schema to create,
	read and delete network instance*/
	return &schema.Resource{
		Description:   "Lumen network resources",
		CreateContext: ResourceNetworkInstanceCreate,
		ReadContext:   ResourceNetworkInstanceRead,
		DeleteContext: ResourceNetworkInstanceDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the instance",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The name of the instance",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"description": {
				Description: "The description of the instance",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"cloud_id": {
				Description: "The cloud id associated with the instance",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},
			"group_id": {
				Description: "The group id associated with the instance",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},
			"instance_type_code": {
				Description: "The instance layout code to provision",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"instance_type_id": {
				Description: "The instance type id to provision",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},
			"instance_layout_id": {
				Description: "The layout to provision the instance from",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},
			"plan_id": {
				Description: "The service plan associated with the instance",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},
			"location": {
				Description: "The edge location in which the instance will be created",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"network_type": {
				Description: "The instance network type",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"bandwidth": {
				Description: "The bandwidth associated to instance",
				Type:        schema.TypeFloat,
				Optional:    true,
				ForceNew:    true,
			},
			"network_id": {
				Description: "The instance network id",
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
			},
			"status": {
				Description: "Instance status",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"instance_cidr": {
				Description: "Network instance CIDR",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"transaction_id": {
				Description: "Network transaction id",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"labels": {
				Type:        schema.TypeList,
				Description: "The list of labels to be added to the instance",
				Optional:    true,
				ForceNew:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tags": {
				Description: "Tags to assign to the instance",
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"date_created": {
				Description: "Timestamp on instance creation",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_updated": {
				Description: "Timestamp on last instance update",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"instance_created_by": {
				Description: "User who created the instance",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"instance_owner": {
				Description: "The instance owner",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

// Function to create network instance
func ResourceNetworkInstanceCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error

	// Initializing client
	c := m.(*Client)

	// Payload
	payload := make(map[string]interface{})
	instanceCloud := d.Get("cloud_id").(int)
	payload["zoneId"] = instanceCloud

	// Instance details
	instance := make(map[string]interface{})
	// Instance name
	instanceName := d.Get("name").(string)
	instance["name"] = instanceName
	// Instance description
	instanceDescription := d.Get("description").(string)
	instance["description"] = instanceDescription

	// Creating instance site details
	site := make(map[string]interface{})
	instanceGroup := d.Get("group_id").(int)
	site["id"] = instanceGroup
	// Adding site to instance payload
	instance["site"] = site

	// Creating instance type payload
	instanceType := make(map[string]interface{})
	// Getting instance type id and code adding to instance payload
	instanceTypeCode := d.Get("instance_type_code").(string)
	instanceType["code"] = instanceTypeCode
	// Adding instance type in instance payload
	instance["instanceType"] = instanceType

	// Getting instance layout id and adding to instance payload
	instanceLayout := make(map[string]interface{})
	instanceLayoutId := d.Get("instance_layout_id").(int)
	instanceLayout["id"] = instanceLayoutId
	instance["layout"] = instanceLayout

	// Getting instance plan and adding to instance payload
	instancePlan := make(map[string]interface{})
	instancePlanId := d.Get("plan_id").(int)
	instancePlan["id"] = instancePlanId
	instance["plan"] = instancePlan
	payload["instance"] = instance

	// Creating instance configs
	config := make(map[string]interface{})
	// Custom options
	edgeLocation := d.Get("location").(string)
	edgeBandwidth := d.Get("bandwidth").(float64)
	instanceNetworkType := d.Get("network_type").(string)

	customOptions := make(map[string]interface{})
	customOptions["edgeLocation"] = edgeLocation
	customOptions["centuryLinkNetworkType"] = instanceNetworkType
	customOptions["edgeBandwidth"] = edgeBandwidth
	customOptions["networkName"] = instanceName
	customOptions["NetworkSize"] = instancePlanId

	// Adding custom config to config
	config["customOptions"] = customOptions
	// Adding config to payload
	payload["config"] = config

	// Getting tags
	if d.Get("tags") != nil {
		tagsinput := d.Get("tags").(map[string]interface{})
		var tags []map[string]interface{}
		for key, value := range tagsinput {
			tag := make(map[string]interface{})
			tag["name"] = key
			tag["value"] = value.(string)
			tags = append(tags, tag)
		}
		payload["tags"] = tags
	}

	// Getting labels
	if d.Get("labels") != nil {
		payload["labels"] = d.Get("labels")
	}

	// Sending request to create network instance
	req := &Request{Body: payload}
	//slcB, _ := json.Marshal(req.Body)
	resp, err := c.CreateInstance(req)
	if err != nil {
		return diag.FromErr(err)
	}
	instanceresult := resp.Result.(*CreateInstanceResult)
	instanceDetails := instanceresult.Instance

	// Instance state confirmation
	instanceStateConf := &resource.StateChangeConf{
		Pending: []string{"provisioning", "starting", "stopping"},
		Target:  []string{"running", "failed", "warning"},
		Refresh: func() (interface{}, string, error) {
			instancedetails, err := c.GetInstance(
				instanceDetails.ID, &Request{})
			if err != nil {
				return "", "", err
			}
			result := instancedetails.Result.(*GetInstanceResult)
			instance := result.Instance
			return result, instance.Status, nil

		},
		Timeout:      25 * time.Minute,
		MinTimeout:   10 * time.Minute,
		Delay:        4 * time.Minute,
		PollInterval: 4 * time.Minute,
	}
	// Checking instance state to catch any errors
	_, err = instanceStateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("ERROR creating instance: %s", err)
	}

	// If instance creation is successfull, setting instance id and name
	d.SetId(int64ToString(instanceDetails.ID))
	d.Set("name", instanceDetails.Name)
	// Fetching created instance details and storing response in output schema
	ResourceNetworkInstanceRead(ctx, d, m)
	return diags
}

// Function to read network instance
func ResourceNetworkInstanceRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{}) diag.Diagnostics {

	// declaring variables
	var diags diag.Diagnostics
	var resp *Response
	var err error

	// Initializing client
	c := m.(*Client)

	instanceId := d.Id()
	instanceName := d.Get("name").(string)
	if instanceId != "" {
		resp, err = c.GetInstance(toInt64(instanceId), &Request{})
	} else if instanceName != "" {
		resp, err = c.FindInstanceByName(instanceName)
	} else {
		return diag.Errorf("INFO: Instance details cannot be retrieved without id or name")
	}

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return diag.FromErr(err)
		} else {
			return diag.FromErr(err)
		}
	}

	// Storing response
	result := resp.Result.(*GetInstanceResult)
	instanceDetails := result.Instance
	if instanceDetails == nil {
		return diag.Errorf("ERROR: Instance details not retrieved in response data")
	}

	// Populating schema with get instance response
	d.Set("description", instanceDetails.Description)
	d.Set("cloud_id", instanceDetails.Cloud["id"])
	d.Set("group_id", instanceDetails.Group["id"])
	d.Set("plan_id", instanceDetails.Plan.ID)
	d.Set("instance_type_id", instanceDetails.InstanceType["id"])
	d.Set("instance_layout_id", instanceDetails.Layout["id"])
	d.Set("status", instanceDetails.Status)

	// Setting instance bandwidth and location
	SetNetworkInstanceCustomConfigs(instanceDetails, d)

	// Setting timestamps for instance creation, last updated
	SetNetworkInstanceTimestamps(instanceDetails, d)

	// Setting user for instance creation
	SetNetworkInstanceUsers(instanceDetails, d)
	return diags
}

// Function to delete network instance
func ResourceNetworkInstanceDelete(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{}) diag.Diagnostics {

	// To collect errors, response and warnings
	var diags diag.Diagnostics
	var resp *Response
	var err error

	// Initializing client
	c := m.(*Client)

	instanceId := d.Id()
	instanceDelRequest := &Request{
		QueryParams: map[string]string{},
	}

	resp, err = c.DeleteInstance(toInt64(instanceId), instanceDelRequest)

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return diag.FromErr(err)
		} else {
			return diag.FromErr(err)
		}
	}
	d.SetId("")
	return diags
}
