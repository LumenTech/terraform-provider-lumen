package lumen

import (
	"context"
	"log"
	"reflect"
	"strconv"
	"terraform-provider-lumen/lumen/client"
	"terraform-provider-lumen/lumen/client/model/morpheus"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceNetworkAllInstances() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: CustomerDeprecationNotice,
		ReadContext:        DataSourceNetworkAllInstancesRead,
		Schema: map[string]*schema.Schema{
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "IDs of network instances created",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"name": {
							Description: "Names of network instances created",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"description": {
							Description: "Instance descriptions",
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
						"instance_type_id": {
							Description: "The type of instance to provision",
							Type:        schema.TypeInt,
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
				},
			},
		},
	}
}

// Func to get all network instances
func DataSourceNetworkAllInstancesRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	// Initializing clients
	c := m.(*Clients).Morpheus

	var diags diag.Diagnostics
	// List network instance call
	resp, err := c.ListInstances(&morpheus.Request{})
	if err != nil {
		return diag.FromErr(err)
	}

	// List network instances
	instances := resp.Result.(*client.ListInstancesResult)

	// Debug
	log.Printf("Instance results: %s", instances)
	// Flattening response to fit schema
	instanceItems := FlattenNetworkInstances(instances.Instances)

	if err := d.Set("instances", instanceItems); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to set instances",
			Detail:   "Unables to set instances",
		})
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}

// Helper function to flatten network instances
func FlattenNetworkInstances(instanceList *[]client.Instance) []interface{} {
	if instanceList != nil {
		instances := make([]interface{}, len(*instanceList)) //, len(*instanceList))
		for i, instanceItem := range *instanceList {
			if instanceItem.InstanceType["code"] == "lumen-network-ip-block" {
				instance := make(map[string]interface{})
				// Populating instance details
				instance["id"] = instanceItem.ID
				instance["name"] = instanceItem.Name
				instance["description"] = instanceItem.Description
				instance["cloud_id"] = instanceItem.Cloud["id"]
				instance["group_id"] = instanceItem.Group["id"]
				instance["plan_id"] = instanceItem.Plan.ID
				instance["instance_type_id"] = instanceItem.InstanceType["id"]
				instance["instance_type_code"] = instanceItem.InstanceType["code"]
				instance["instance_layout_id"] = instanceItem.Layout["id"]
				instance["status"] = instanceItem.Status

				// Setting location, bandwidth
				customOptions := instanceItem.Config["customOptions"]
				v := reflect.ValueOf(customOptions)
				if v.Kind() == reflect.Map {
					for _, key := range v.MapKeys() {
						strct := v.MapIndex(key)
						if key.Interface().(string) == "edgeLocation" {
							instance["instance_location"] = strct.Interface().(string)
						} else if key.Interface().(string) == "transactionId" {
							instance["transaction_id"] = strct.Interface().(string)
						} else if key.Interface().(string) == "centuryLinkNetworkType" {
							instance["network_type"] = strct.Interface().(string)
						} else if key.Interface().(string) == "cidr" {
							instance["network_cidr"] = strct.Interface().(string)
						} else if key.Interface().(string) == "edgeBandwidth" {
							if _, ok := strct.Interface().(float64); ok {
								instance["network_bandwidth"] = strct.Interface().(float64)
							}
						}
					}
				}

				// Setting timestamps for instance creation, last updated
				instance["date_created"] = instanceItem.DateCreated
				instance["last_updated"] = instanceItem.LastUpdated

				// Setting user for instance creation
				valueCreatedBy := reflect.ValueOf(instanceItem.CreatedBy)
				if valueCreatedBy.Kind() == reflect.Map {
					for _, key := range valueCreatedBy.MapKeys() {
						strct := valueCreatedBy.MapIndex(key)
						if key.Interface().(string) == "username" {
							instance["instance_created_by"] = strct.Interface().(string)
						}
					}
				}
				// Setting instance owner
				valueOwner := reflect.ValueOf(instanceItem.CreatedBy)
				if valueOwner.Kind() == reflect.Map {
					for _, key := range valueOwner.MapKeys() {
						strct := valueOwner.MapIndex(key)
						if key.Interface().(string) == "username" {
							instance["instance_owner"] = strct.Interface().(string)
						}
					}
				}
				// Adding instance details in schema
				instances[i] = instance
			}
		}
		return instances
	}
	return make([]interface{}, 0)
}
