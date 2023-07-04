package lumen

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceBareMetalInstanceId() *schema.Resource {
	/*
		return schema for an instance read
		 order based on instance id
	*/
	return &schema.Resource{
		Description: "Provides Lumen instance details",
		ReadContext: DataSourceBareMetalInstanceIdRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the instance",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"name": {
				Description: "The name of the instance",
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
				Description: "The layout to provision the instance",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"plan_id": {
				Description: "The service plan associated with the instance",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"location": {
				Description: "The instance location",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"instance_ip": {
				Description: "The instance IP address",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"labels": {
				Type:        schema.TypeList,
				Description: "The list of labels to add to the instance",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tags": {
				Description: "Tags to assign to the instance",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString},
			},
			"status": {
				Description: "Instance status",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"volumes": {
				Description: "The instance volumes",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_storage": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"short_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resizeable": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"root_volume": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"interfaces": {
				Description: "The instance network interfaces to create",
				Type:        schema.TypeList,
				Computed:    true,
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
			"network_type": {
				Description: "The network type associated with the resource",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"network_id": {
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

/*
Retrieve instance details with Id
*/
func DataSourceBareMetalInstanceIdRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{}) diag.Diagnostics {

	// Initializing client
	c := m.(*Client)

	// Warnings or error to be collected in a slice type
	var diags diag.Diagnostics
	var resp *Response
	var err error

	instanceid := strconv.Itoa(d.Get("id").(int))
	resp, err = c.GetInstance(toInt64(instanceid), &Request{})

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return nil
		} else {
			return diag.FromErr(err)
		}
	}

	// storing get instance response data
	instancedetails := resp.Result.(*GetInstanceResult)

	// populating schema with instance response
	PopulateSchemaInstanceIdResponse(instancedetails.Instance, d)
	// Setting instance id
	d.SetId(instanceid)
	return diags
}

// helper function to populate response schema
func PopulateSchemaInstanceIdResponse(
	instanceDetails *Instance,
	d *schema.ResourceData) {
	if instanceDetails != nil {
		d.Set("name", instanceDetails.Name)
		d.Set("cloud_id", instanceDetails.Cloud["id"])
		d.Set("group_id", instanceDetails.Group["id"])
		d.Set("instance_type_id", instanceDetails.InstanceType["id"])
		d.Set("instance_layout_id", instanceDetails.Layout["id"])
		d.Set("plan_id", instanceDetails.Plan.ID)
		d.Set("status", instanceDetails.Status)
		d.Set("labels", instanceDetails.Labels)

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
	}
}
