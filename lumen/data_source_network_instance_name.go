package lumen

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceNetworkInstanceName() *schema.Resource {
	/* return schema for network instance
	read based on instance name.
	*/
	return &schema.Resource{
		Description: "Provides Lumen network instance details based on instance name",
		ReadContext: DataSourceNetworkInstanceNameRead,
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
			"location": {
				Description: "The network instance location",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"instance_type": {
				Description: "The network instance type",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"instance_bandwidth": {
				Description: "The network instance bandwidth",
				Type:        schema.TypeFloat,
				Computed:    true,
			},
			"instance_cidr": {
				Description: "CIDR associated with network instance",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"network_type": {
				Description: "The network type associated with the resource",
				Type:        schema.TypeString,
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
			"instance_owner": {
				Description: "The instance owner",
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

	// Initializing client
	c := m.(*Client)

	// To collect warnings and errors in a slice type
	var diags diag.Diagnostics
	var resp *Response
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
	instanceDetails := resp.Result.(*GetInstanceResult)

	// populating schema based on response
	PopulateSchemaNetworkInstanceNameResponse(instanceDetails.Instance, d)
	return diags
}

// flattening response
func PopulateSchemaNetworkInstanceNameResponse(
	instanceDetails *Instance,
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
		d.Set("instance_layout_id", instanceDetails.Layout["id"])
		d.Set("status", instanceDetails.Status)
		d.Set("instance_type", instanceDetails.InstanceType["name"])

		// Setting instance bandwidth and location
		SetNetworkInstanceCustomConfigs(instanceDetails, d)

		// Setting timestamps for instance creation, last updated
		SetNetworkInstanceTimestamps(instanceDetails, d)

		// Setting user for instance creation
		SetNetworkInstanceUsers(instanceDetails, d)
	}
}
