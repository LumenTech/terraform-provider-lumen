package lumen

import (
	"context"
	//"encoding/json"
	"net"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceBareMetalInstance() *schema.Resource {
	return &schema.Resource{
		Description:   "Creates Lumen instance resource",
		CreateContext: ResourceBareMetalInstanceCreate,
		ReadContext:   ResourceBareMetalInstanceRead,
		UpdateContext: ResourceBareMetalInstanceUpdate,
		DeleteContext: ResourceBareMetalInstanceDelete,
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
			"instance_cloud_name": {
				Description: "The instance cloud name",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"instance_type": {
				Description: "The instance layout type to provision",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"instance_type_code": {
				Description: "The instance layout code to provision",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"instance_layout_id": {
				Description: "The layout to provision the instance from",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"instance_resource_pool_id": {
				Description: "The ID of the resource pool to provision the instance to",
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
			"bandwidth": {
				Description: "The bandwidth assigned to instance",
				Type:        schema.TypeFloat,
				Optional:    true,
			},
			"network_type": {
				Description: "The instance network type",
				Type:        schema.TypeString,
				Optional:    true,
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
			"version": {
				Type:     schema.TypeString,
				Computed: true,
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

	// Initializing client
	c := m.(*Client)

	// Getting instance name, cloud, and group id
	instancename := d.Get("name").(string)
	instancegroup := d.Get("group_id").(int)
	instancecloud := d.Get("cloud_id").(int)

	// payload
	payload := make(map[string]interface{})
	payload["zoneId"] = instancecloud

	// instance details
	instance := make(map[string]interface{})
	instance["name"] = instancename

	// instance description
	instancedescription := d.Get("description").(string)
	instance["description"] = instancedescription

	// Creating instance site details
	site := make(map[string]interface{})
	site["id"] = instancegroup
	// Adding site to instance payload
	instance["site"] = site

	// Getting instance type and adding to instance payload
	instancetype := d.Get("instance_type").(string)
	instance["type"] = instancetype

	// Getting instance layout code and adding to instance payload
	instancetypecode := make(map[string]interface{})
	instancelayoutcode := d.Get("instance_type_code").(string)
	instancetypecode["code"] = instancelayoutcode
	instance["instanceType"] = instancetypecode

	// Getting instance layout id and adding to instance payload
	instancelayout := make(map[string]interface{})
	instancelayoutid := d.Get("instance_layout_id").(int)
	instancelayout["id"] = instancelayoutid
	instance["layout"] = instancelayout

	// Getting instance plan and adding to instance payload
	instanceplan := make(map[string]interface{})
	instanceplanid := d.Get("plan_id").(int)
	instanceplan["id"] = instanceplanid
	instance["plan"] = instanceplan

	// Adding instance details to payload
	payload["instance"] = instance

	// Creating instance config
	config := make(map[string]interface{})

	// Getting config data
	resourcepoolid := d.Get("instance_resource_pool_id").(int)
	config["resourcePoolId"] = resourcepoolid
	// Custom configs
	edgelocation := d.Get("location").(string)
	edgebandwidth := d.Get("bandwidth").(float64)
	instancenetworktype := d.Get("network_type").(string)

	customoptions := make(map[string]interface{})
	customoptions["edgeLocation"] = edgelocation
	customoptions["edgeBandwidth"] = edgebandwidth
	customoptions["centuryLinkNetworkType"] = instancenetworktype
	// Adding custom config to config
	config["customOptions"] = customoptions
	// Create user and add to payload
	createuser := d.Get("create_user")
	config["createUser"] = createuser

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
	// Fetching Instance details and storing response for user
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
	var resp *Response
	var err error

	// Initializing client
	c := m.(*Client)

	instanceid := d.Id()
	instancename := d.Get("name").(string)
	if instanceid == "" && instancename != "" {
		resp, err = c.FindInstanceByName(instancename)
	} else if instanceid != "" {
		resp, err = c.GetInstance(toInt64(instanceid), &Request{})
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
	result := resp.Result.(*GetInstanceResult)
	instanceDetails := result.Instance
	if instanceDetails == nil {
		return diag.Errorf("ERROR: Instance details not retrieved in response data")
	}

	// Populating schema with morpheus response
	d.Set("description", instanceDetails.Description)
	d.Set("cloud_id", instanceDetails.Cloud["id"])
	d.Set("group_id", instanceDetails.Group["id"])
	d.Set("instance_type_id", instanceDetails.InstanceType["id"])
	d.Set("instance_type_layout", instanceDetails.Layout["id"])
	d.Set("plan_id", instanceDetails.Plan.ID)
	d.Set("resource_pool_id", instanceDetails.Config["resourcePoolId"])
	d.Set("environment", instanceDetails.Environment)
	d.Set("labels", instanceDetails.Labels)
	d.Set("version", instanceDetails.Version)
	d.Set("status", instanceDetails.Status)

	// Setting location and bandwidth
	customOptions := instanceDetails.Config["customOptions"]
	v := reflect.ValueOf(customOptions)
	if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			strct := v.MapIndex(key)
			if key.Interface().(string) == "edgeLocation" {
				d.Set("location", strct.Interface().(string))
			} else if key.Interface().(string) == "edgeBandwidth" {
				d.Set("bandwidth", strct.Interface().(float64))
			}
		}
	}

	// Setting instance ip address
	envVars := instanceDetails.EnvironmentVariables
	for _, envVarsItems := range *envVars {
		envVarIP := envVarsItems["value"].(string)
		varIp := net.ParseIP(envVarIP)
		if varIp.To4() != nil {
			d.Set("instance_ip", envVarIP)
			break
		}
	}

	// Setting tags
	tags := make(map[string]interface{})
	if instanceDetails.Tags != nil {
		tagDetails := instanceDetails.Tags
		tagsList := *tagDetails
		for i := 0; i < len(tagsList); i++ {
			tag := tagsList[i]
			tagName := tag["name"]
			tags[tagName.(string)] = tag["value"]

		}
	}
	d.Set("tags", tags)

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
	var resp *Response
	var err error

	// Initializing client
	c := m.(*Client)

	instanceid := d.Id()
	morphRequest := &Request{
		QueryParams: map[string]string{},
	}

	resp, err = c.DeleteInstance(toInt64(instanceid), morphRequest)

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

/* Function to update instance based on instance id.
This updates instance name, description, labels, tags.
*/
func ResourceBareMetalInstanceUpdate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{}) diag.Diagnostics {

	// To collect errors, warnings
	var diags diag.Diagnostics
	var err error

	// Initializing client
	c := m.(*Client)

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
	morphRequest := &Request{Body: payload}
	resp, err := c.UpdateInstance(toInt64(id), morphRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	updateInstanceResult := resp.Result.(*UpdateInstanceResult)
	instanceDetails := updateInstanceResult.Instance

	// Updated resource successfully
	// Setting instance id
	d.SetId(int64ToString(instanceDetails.ID))
	ResourceBareMetalInstanceRead(ctx, d, m)
	return diags
}
