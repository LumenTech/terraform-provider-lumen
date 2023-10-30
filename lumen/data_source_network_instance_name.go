package lumen

import (
	"context"
	"log"
	"reflect"
	"terraform-provider-lumen/lumen/client"
	"terraform-provider-lumen/lumen/client/model/morpheus"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceNetworkInstanceName() *schema.Resource {
	/* return schema for network instance
	read based on instance name.
	*/
	return &schema.Resource{
		Description:        "Provides Lumen network instance details based on instance name",
		DeprecationMessage: CustomerDeprecationNotice,
		ReadContext:        DataSourceNetworkInstanceNameRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The network instance name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"id": {
				Description: "The network instance id",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"description": {
				Description: "Instance description",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"cloud_id": {
				Description: "The cloud id associated with the instance",
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
			"instance_type_code": {
				Description: "The network instance type",
				Type:        schema.TypeString,
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

// Retrieve network instance details based on instance name
func DataSourceNetworkInstanceNameRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{}) diag.Diagnostics {

	// Initializing clients
	c := m.(*Clients).Morpheus

	// To collect warnings and errors in a slice type
	var diags diag.Diagnostics
	var resp *morpheus.Response
	var err error

	instanceName := d.Get("name").(string)
	resp, err = c.FindInstanceByName(instanceName)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return nil
		} else {
			return diag.FromErr(err)
		}
	}
	instanceDetails := resp.Result.(*client.GetInstanceResult)

	// populating schema based on response
	PopulateSchemaNetworkInstanceNameResponse(instanceDetails.Instance, d)
	return diags
}

// flattening response
func PopulateSchemaNetworkInstanceNameResponse(
	instanceDetails *client.Instance,
	d *schema.ResourceData) {
	if instanceDetails != nil {
		log.Printf("instance name: ", instanceDetails.Name)
		d.Set("name", instanceDetails.Name)
		log.Printf("instance id: ", instanceDetails.ID)
		d.Set("id", instanceDetails.ID)
		log.Printf("instance description: ", instanceDetails.Description)
		d.Set("description", instanceDetails.Description)
		d.Set("cloud_id", instanceDetails.Cloud["id"])
		d.Set("group_id", instanceDetails.Group["id"])
		d.Set("plan_id", instanceDetails.Plan.ID)
		d.Set("instance_type_id", instanceDetails.InstanceType["id"])
		d.Set("instance_type_code", instanceDetails.InstanceType["code"])
		d.Set("instance_layout_id", instanceDetails.Layout["id"])
		d.Set("status", instanceDetails.Status)

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
					if _, ok := strct.Interface().(float64); ok {
						d.Set("network_bandwidth", strct.Interface().(float64))
					}
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
