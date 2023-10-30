package lumen

import (
	"context"
	"terraform-provider-lumen/lumen/client"
	"terraform-provider-lumen/lumen/client/model/morpheus"

	//"encoding/json"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceBareMetalInstance() *schema.Resource {
	return &schema.Resource{
		Description:        "Creates Lumen instance resource",
		DeprecationMessage: CustomerDeprecationNotice,
		CreateContext:      ResourceBareMetalInstanceCreate,
		ReadContext:        ResourceBareMetalInstanceRead,
		UpdateContext:      ResourceBareMetalInstanceUpdate,
		DeleteContext:      ResourceBareMetalInstanceDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Read:   schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(90 * time.Minute),
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
			},
			"description": {
				Description: "The description of the instance",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"group_id": {
				Description: "The group id associated with the instance",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"cloud_id": {
				Description: "The cloud id associated with the instance",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"plan_id": {
				Description: "The service plan associated with the instance",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"instance_type_id": {
				Description: "The instance type id to provision",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"instance_type_code": {
				Description: "The instance type layout code to provision",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"instance_layout_id": {
				Description: "The layout to provision the instance from",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"create_user": {
				Description: "Whether to create an user account on the instance",
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Default:     true,
			},
			"location": {
				Description: "The edge location in which the instance will be created",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"is_vpc_selectable": {
				Description: "Vpc Selectable option",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"bare_metal_network_type": {
				Description: "The bare metal network type",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"network_resource_id": {
				Description: "The network instance id",
				Type:        schema.TypeString,
				Required:    true,
			},
			"instance_ip": {
				Description: "The instance ip address",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"labels": {
				Type:        schema.TypeList,
				Description: "The list of labels to be added to the instance",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tags": {
				Description: "Tags to assign to the instance",
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"addtags": {
				Description: "Tags to be added when updated",
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"removetags": {
				Description: "Tags to be removed when updated",
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Description: "Instance status",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"evar": {
				Type:        schema.TypeList,
				Description: "The environment variables to create",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the environment variable",
							Optional:    true,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "The value of the environment variable",
							Optional:    true,
						},
						"export": {
							Type:        schema.TypeBool,
							Description: "Whether the environment variable is exported as an instance tag",
							Optional:    true,
						},
						"masked": {
							Type:        schema.TypeBool,
							Description: "Whether the environment variable is masked for security purposes",
							Optional:    true,
						},
					},
				},
			},
			"volumes": {
				Description: "The instance volumes to create",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "The name/type of the volume being created",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"root_volume": {
							Description: "Whether the volume is the root volume of the instance",
							Type:        schema.TypeBool,
							Optional:    true,
						},
						"size": {
							Description: "The size of the volume being created",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"size_id": {
							Description: "The ID of an existing volume to assign to the instance",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"storage_type": {
							Description: "The ID of the volume type",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"datastore_id": {
							Description: "The ID of the datastore",
							Type:        schema.TypeInt,
							Optional:    true,
						},
					},
				},
			},
			"interfaces": {
				Description: "The instance network interfaces to create",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_id": {
							Description: "The network to assign the network interface to",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"ip_address": {
							Description: "",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"ip_mode": {
							Description: "",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"network_interface_type_id": {
							Description: "The network interface type",
							Type:        schema.TypeInt,
							Optional:    true,
						},
					},
				},
			},
			"date_created": {
				Description: "Instance creation date",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_updated": {
				Description: "Instance last updated",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"instance_created_by": {
				Description: "Instance created by user",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

/*
Function to create instance
*/
func ResourceBareMetalInstanceCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error

	// Initializing clients
	c := m.(*Clients).Morpheus

	// payload
	payload := make(map[string]interface{})

	// Getting instance name, cloud, and group id
	instanceName := d.Get("name").(string)
	instanceCloudId := d.Get("cloud_id").(int)
	instanceGroupId := d.Get("group_id").(int)

	payload["zoneId"] = instanceCloudId

	// instance details
	instance := make(map[string]interface{})
	instance["name"] = instanceName

	// instance description
	instanceDescription := d.Get("description").(string)
	instance["description"] = instanceDescription

	// Adding instance site details to payload
	site := make(map[string]interface{})
	site["id"] = instanceGroupId
	instance["site"] = site

	// Adding instance type id and code.
	instanceType := make(map[string]interface{})
	instanceType["id"] = d.Get("instance_type_id").(string)
	instanceType["code"] = d.Get("instance_type_code").(string)
	instance["instanceType"] = instanceType

	// Getting instance layout id and adding to instance payload
	instanceLayout := make(map[string]interface{})
	instanceLayoutId := d.Get("instance_layout_id")
	instanceLayout["id"] = instanceLayoutId
	instance["layout"] = instanceLayout

	// Adding instance hostname
	instance["hostName"] = instanceName

	// Getting instance plan and adding to instance payload
	instancePlan := make(map[string]interface{})
	instancePlanId := d.Get("plan_id").(int)
	instancePlan["id"] = instancePlanId
	instance["plan"] = instancePlan

	// Adding instance details to payload
	payload["instance"] = instance

	// Creating instance config
	config := make(map[string]interface{})

	// Custom configs
	customOptions := make(map[string]interface{})
	customOptions["edgeLocation"] = d.Get("location").(string)
	customOptions["networkId"] = d.Get("network_resource_id").(string)
	customOptions["centuryLinkNetworkType"] = d.Get("bare_metal_network_type").(string)
	config["is_vpc_selectable"] = d.Get("is_vpc_selectable")
	config["create_user"] = d.Get("create_user")
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

	// Getting environment variables
	if d.Get("evar") != nil {
		payload["evars"] = parseEnvironmentVariables(d.Get("evar").([]interface{}))
	}

	// Getting Network interfaces
	if d.Get("interfaces") != nil {
		payload["networkInterfaces"] = parseNetworkInterfaces(d.Get("interfaces").([]interface{}))
	}

	// Volumes
	if d.Get("volumes") != nil {
		payload["volumes"] = parseStorageVolumes(d.Get("volumes").([]interface{}))
	}

	req := &morpheus.Request{Body: payload}
	//slcB, _ := json.Marshal(req.Body)
	resp, err := c.CreateInstance(req)
	if err != nil {
		return diag.FromErr(err)
	}
	instanceresult := resp.Result.(*client.CreateInstanceResult)
	instanceDetails := instanceresult.Instance

	// Instance state confirmation
	instanceStateConf := &resource.StateChangeConf{
		Pending: []string{"provisioning", "starting", "stopping"},
		Target:  []string{"running", "failed", "warning"},
		Refresh: func() (interface{}, string, error) {
			instancedetails, err := c.GetInstance(
				instanceDetails.ID, &morpheus.Request{})
			if err != nil {
				return "", "", err
			}
			result := instancedetails.Result.(*client.GetInstanceResult)
			instance := result.Instance
			return result, instance.Status, nil

		},
		Timeout:      3 * time.Hour,
		MinTimeout:   10 * time.Minute,
		Delay:        4 * time.Minute,
		PollInterval: 3 * time.Minute,
	}
	// Checking instance state to catch any errors
	_, err = instanceStateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("ERROR creating instance: %s", err)
	}

	// Successfully created instance, setting instance id and name
	d.SetId(int64ToString(instanceDetails.ID))
	d.Set("name", instanceDetails.Name)
	// Fetching created instance details and storing response in output schema
	ResourceBareMetalInstanceRead(ctx, d, m)
	return diags
}

/*
Function to read created / modified instance
based on instance id or instance Name.
Returns error if details not retrieved.
*/
func ResourceBareMetalInstanceRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{}) diag.Diagnostics {

	// declaring variables
	var diags diag.Diagnostics
	var resp *morpheus.Response
	var err error

	// Initializing clients
	c := m.(*Clients).Morpheus

	instanceid := d.Id()
	instancename := d.Get("name").(string)
	if instanceid != "" {
		resp, err = c.GetInstance(toInt64(instanceid), &morpheus.Request{})
	} else if instancename != "" {
		resp, err = c.FindInstanceByName(instancename)
	} else {
		return diag.Errorf(
			"INFO: Instance details cannot be retrieved without id or name")
	}

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return diag.FromErr(err)
		} else {
			return diag.FromErr(err)
		}
	}

	// Storing resource data
	result := resp.Result.(*client.GetInstanceResult)
	instanceDetails := result.Instance
	if instanceDetails == nil {
		return diag.Errorf("ERROR: Instance details not retrieved in response data")
	}

	// Populating schema with get instance response
	d.Set("description", instanceDetails.Description)
	d.Set("cloud_id", instanceDetails.Cloud["id"])
	d.Set("group_id", instanceDetails.Group["id"])
	d.Set("instance_type_layout", instanceDetails.Layout["id"])
	d.Set("plan_id", instanceDetails.Plan.ID)
	d.Set("labels", instanceDetails.Labels)
	d.Set("status", instanceDetails.Status)

	// Setting instance custom configs
	SetBareMetalInstanceCustomConfigs(instanceDetails, d)

	// Setting instance connection info
	SetBareMetalInstanceConnectionInfo(instanceDetails, d)

	// Setting timestamps for instance creation, last updated
	SetBareMetalInstanceTimestamps(instanceDetails, d)

	// Setting user for instance user and owner
	SetBareMetalInstanceUsers(instanceDetails, d)

	// Setting tags
	SetBareMetalInstanceTags(instanceDetails, d)

	// Setting volumes
	SetBareMetalInstanceVolumes(instanceDetails, d)
	return diags
}

/*
Function to delete instance based on instance id.
Throws error in case of functional issues.
*/
func ResourceBareMetalInstanceDelete(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{}) diag.Diagnostics {

	// To collect errors, response and warnings
	var diags diag.Diagnostics
	var resp *morpheus.Response
	var err error

	// Initializing clients
	c := m.(*Clients).Morpheus

	instanceid := d.Id()
	instanceDelRequest := &morpheus.Request{
		QueryParams: map[string]string{},
	}

	resp, err = c.DeleteInstance(toInt64(instanceid), instanceDelRequest)

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

/*
Function to update instance based on instance id.
This updates instance name, description, labels, tags.
*/
func ResourceBareMetalInstanceUpdate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{}) diag.Diagnostics {

	// To collect errors, warnings
	var diags diag.Diagnostics
	var err error

	// Initializing clients
	c := m.(*Clients).Morpheus

	id := d.Id()

	// Instance payload
	instancePayload := make(map[string]interface{})

	// Check name
	if d.HasChange("name") {
		name := d.Get("name").(string)
		instancePayload["name"] = name
	}
	// Check description
	if d.HasChange("description") {
		description := d.Get("description").(string)
		instancePayload["description"] = description
	}
	// Check labels
	if d.HasChange("labels") {
		labels := d.Get("labels")
		instancePayload["labels"] = labels
	}

	// Update tags
	// Only for replacing existing tags and add the new ones
	var tags []map[string]interface{}
	if d.HasChange("tags") {
		tagsInput := d.Get("tags").(map[string]interface{})
		for key, value := range tagsInput {
			tag := make(map[string]interface{})
			tag["name"] = key
			tag["value"] = value.(string)
			tags = append(tags, tag)
		}
		instancePayload["tags"] = tags
	}

	// Add tags
	var atags []map[string]interface{}
	if d.Get("addtags") != nil {
		addTags := d.Get("addtags").(map[string]interface{})
		for key, value := range addTags {
			atag := make(map[string]interface{})
			atag["name"] = key
			atag["value"] = value.(string)
			atags = append(atags, atag)
		}
		// Setting instance payload for adding tags
		instancePayload["addTags"] = atags
	}

	// Remove tags
	var rtags []map[string]interface{}
	if d.Get("removetags") != nil {
		removeTags := d.Get("removetags").(map[string]interface{})
		for key, value := range removeTags {
			rtag := make(map[string]interface{})
			rtag["name"] = key
			rtag["value"] = value.(string)
			rtags = append(rtags, rtag)
		}
		// Setting instance payload for removing tags
		instancePayload["removeTags"] = rtags
	}

	// Payload
	payload := make(map[string]interface{})
	payload["instance"] = instancePayload

	// API request
	morphRequest := &morpheus.Request{Body: payload}
	resp, err := c.UpdateInstance(toInt64(id), morphRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	updateInstanceResult := resp.Result.(*client.UpdateInstanceResult)
	instanceDetails := updateInstanceResult.Instance

	// Updated resource successfully
	// Setting instance id
	d.SetId(int64ToString(instanceDetails.ID))
	ResourceBareMetalInstanceRead(ctx, d, m)
	return diags
}
