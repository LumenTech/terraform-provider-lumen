package lumen

import (
	"context"
	"net"
	"reflect"
	"strconv"
	"time"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceBareMetalAllInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceBareMetalAllInstancesRead,
		Schema: map[string]*schema.Schema{
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "The IDs of Instances created",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"name": {
							Description: "The names of instances created",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"description": {
							Description: "The instance description",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"cloud_id": {
							Description: "The ID of the cloud associated with the instance",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"group_id": {
							Description: "The ID of the group associated with the instance",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"instance_type_id": {
							Description: "The type of instance to provision",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"instance_layout_id": {
							Description: "The layout to provision the instance from",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"plan_id": {
							Description: "The service plan associated with the instance",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"environment": {
							Description: "The environment to assign the instance to",
							Type:        schema.TypeString,
							Computed:    true,
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
						"instance_location": {
							Description: "The instance location",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"instance_ip": {
							Description: "The instance ip address",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func DataSourceBareMetalAllInstancesRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{}) diag.Diagnostics {

	// Initializing morpheus client
	c := m.(*morpheus.Client)

	var diags diag.Diagnostics

	// List instance call
	resp, err := c.ListInstances(&morpheus.Request{})
	if err != nil {
		return diag.FromErr(err)
	}

	// List of instances
	instances := resp.Result.(*morpheus.ListInstancesResult)

	// Flattening response to fit schema
	instanceItems := FlattenInstances(instances.Instances)

	if err := d.Set("instances", instanceItems); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to set instances",
			Detail:   "Unable to set instances",
		})
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}

func FlattenInstances(instanceList *[]morpheus.Instance) []interface{} {
	if instanceList != nil {
		instances := make([]interface{}, len(*instanceList), len(*instanceList))

		for i, instanceItem := range *instanceList {
			instance := make(map[string]interface{})
			// Populating instance details
			instance["id"] = instanceItem.ID
			instance["name"] = instanceItem.Name
			instance["description"] = instanceItem.Description
			instance["cloud_id"] = instanceItem.Cloud["id"]
			instance["group_id"] = instanceItem.Group["id"]
			instance["plan_id"] = instanceItem.Plan.ID
			instance["instance_type_id"] = instanceItem.InstanceType["id"]
			instance["instance_layout_id"] = instanceItem.Layout["id"]
			instance["environment"] = instanceItem.Environment
			instance["version"] = instanceItem.Version
			instance["status"] = instanceItem.Status

			// Setting location, bandwidth
			customOptions := instanceItem.Config["customOptions"]
			v := reflect.ValueOf(customOptions)
			if v.Kind() == reflect.Map {
				for _, key := range v.MapKeys() {
					strct := v.MapIndex(key)
					if key.Interface().(string) == "edgeLocation" {
						instance["instance_location"] = strct.Interface().(string)
					}
				}
			}

			// Setting instance ip address
			envVars := instanceItem.EnvironmentVariables
			for _, envVarsItems := range *envVars {
				envVarIP := envVarsItems["value"].(string)
				varIp := net.ParseIP(envVarIP)
				if varIp.To4() != nil {
					instance["instance_ip"] = envVarIP
					break
				}
			}
			// Adding instance details
			instances[i] = instance
		}

		return instances
	}

	return make([]interface{}, 0)
}
