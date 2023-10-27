package lumen

import (
	"context"
	"reflect"
	"strconv"
	"terraform-provider-lumen/lumen/client"
	"terraform-provider-lumen/lumen/client/model/morpheus"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceNetworkInstanceId() *schema.Resource {
	/* Return schema for reading network
	instance details based on instance id.
	*/
	return &schema.Resource{
		Description: "Provides Lumen network instance details based on instance id",
		ReadContext: DataSourceNetworkInstanceIdRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "IDs of network instances created",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"name": {
				Description: "Names of network instances created",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "Instance description",
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
				Description: "The instance layout id",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"plan_id": {
				Description: "The service plan associated with the instance",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"status": {
				Description: "Instance status",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"instance_location": {
				Description: "The network instance location",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"instance_type_code": {
				Description: "The network instance type",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"network_bandwidth": {
				Description: "The network instance bandwidth",
				Type:        schema.TypeFloat,
				Computed:    true,
			},
			"network_cidr": {
				Description: "CIDR associated with network instance",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"network_type": {
				Description: "The network type associated with the resource",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"network_resource_id": {
				Description: "Network resource id",
				Type:        schema.TypeFloat,
				Computed:    true,
			},
			"transaction_id": {
				Description: "The network id associated with the instance",
				Type:        schema.TypeString,
				Computed:    true,
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
		},
	}
}

// Function to get network instance details
func DataSourceNetworkInstanceIdRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{}) diag.Diagnostics {

	// Initializing client
	c := m.(*Client).Morpheus

	// Capture warings and errors
	var diags diag.Diagnostics
	var resp *morpheus.Response
	var err error

	instanceId := strconv.Itoa(d.Get("id").(int))
	resp, err = c.GetInstance(toInt64(instanceId), &morpheus.Request{})

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return nil
		} else {
			return diag.FromErr(err)
		}
	}

	// storing get instance response data
	instanceDetails := resp.Result.(*client.GetInstanceResult)

	// populating schema with instance response
	PopulateSchemaNetworkInstanceIdResponse(instanceDetails.Instance, d)
	// Setting instance id
	d.SetId(instanceId)
	return diags
}

// helper function to populate response schema
func PopulateSchemaNetworkInstanceIdResponse(
	instanceDetails *client.Instance,
	d *schema.ResourceData) {
	if instanceDetails != nil {
		d.Set("name", instanceDetails.Name)
		d.Set("description", instanceDetails.Description)
		d.Set("cloud_id", instanceDetails.Cloud["id"])
		d.Set("group_id", instanceDetails.Group["id"])
		d.Set("plan_id", instanceDetails.Plan.ID)
		d.Set("instance_type_id", instanceDetails.InstanceType["id"])
		d.Set("instance_layout_id", instanceDetails.Layout["id"])
		d.Set("status", instanceDetails.Status)
		d.Set("instance_type_code", instanceDetails.InstanceType["code"])

		customOptions := instanceDetails.Config["customOptions"]
		v := reflect.ValueOf(customOptions)
		if v.Kind() == reflect.Map {
			for _, key := range v.MapKeys() {
				strct := v.MapIndex(key)
				if key.Interface().(string) == "edgeLocation" {
					d.Set("instance_location", strct.Interface().(string))
				} else if key.Interface().(string) == "centuryLinkNetworkType" {
					d.Set("network_type", strct.Interface().(string))
				} else if key.Interface().(string) == "edgeBandwidth" {
					d.Set("network_bandwidth", strct.Interface().(float64))
				}
			}
		}

		// Setting transaction id, cidr, network_resource_id
		SetNetworkInstanceCustomConfigs(instanceDetails, d)

		// Setting timestamps for instance creation, last updated
		SetNetworkInstanceTimestamps(instanceDetails, d)

		// Setting user for instance creation
		SetNetworkInstanceUsers(instanceDetails, d)

	}
}
