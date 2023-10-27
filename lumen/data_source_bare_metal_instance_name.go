package lumen

import (
	"context"
	"net"
	"reflect"
	"terraform-provider-lumen/lumen/client"
	"terraform-provider-lumen/lumen/client/model/morpheus"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceBareMetalInstanceName() *schema.Resource {
	/*
		return schema for reading bare metal instance
		 details based on instance name.
	*/
	return &schema.Resource{
		Description:        "Provides Lumen instance details based on instance name",
		DeprecationMessage: CustomerDeprecationNotice,
		ReadContext:        DataSourceBareMetalInstanceNameRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Instance Name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"id": {
				Description: "Instance ID",
				Type:        schema.TypeInt,
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
				Description: "The layout to provision the instance from",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"plan_id": {
				Description: "The service plan associated with the instance",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"resource_pool_id": {
				Description: "The ID of the resource pool to provision the instance to",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"environment": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_location": {
				Description: "The instance location",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"instance_bandwidth": {
				Description: "The instance bandwidth",
				Type:        schema.TypeFloat,
				Computed:    true,
			},
			"instance_ip": {
				Description: "The instance ip address",
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
			"evars": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"export": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"masked": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

/*
Retrieve instance details based on instance name
*/
func DataSourceBareMetalInstanceNameRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{}) diag.Diagnostics {

	// Initializing client
	c := m.(*Client).Morpheus

	// Warnings or error to be collected in a slice type
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

	var instancedetails *client.GetInstanceResult
	instancedetails = resp.Result.(*client.GetInstanceResult)

	// populating schema based on response.
	PopulateSchemaInstanceNameResponse(instancedetails.Instance, d)

	return diags
}

// flattening response
func PopulateSchemaInstanceNameResponse(
	instanceDetails *client.Instance,
	d *schema.ResourceData) {
	if instanceDetails != nil {
		d.Set("name", instanceDetails.Name)
		d.Set("id", instanceDetails.ID)
		d.Set("cloud_id", instanceDetails.Cloud["id"])
		d.Set("group_id", instanceDetails.Group["id"])
		d.Set("instance_type_id", instanceDetails.InstanceType["id"])
		d.Set("instance_layout_id", instanceDetails.Layout["id"])
		d.Set("plan_id", instanceDetails.Plan.ID)
		d.Set("resource_pool_id", instanceDetails.Config["resourcePoolId"])
		d.Set("environment", instanceDetails.Environment)
		d.Set("status", instanceDetails.Status)
		d.Set("labels", instanceDetails.Labels)

		// Setting instance bandwidth and location
		customOptions := instanceDetails.Config["customOptions"]
		v := reflect.ValueOf(customOptions)
		if v.Kind() == reflect.Map {
			for _, key := range v.MapKeys() {
				strct := v.MapIndex(key)
				if key.Interface().(string) == "edgeLocation" {
					d.Set("instance_location", strct.Interface().(string))
				} else if key.Interface().(string) == "edgeBandwidth" {
					d.Set("instance_bandwidth", strct.Interface().(float64))
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
			instanceTags := instanceDetails.Tags
			tagList := *instanceTags
			for i := 0; i < len(tagList); i++ {
				tag := tagList[i]
				tagName := tag["name"]
				tags[tagName.(string)] = tag["value"]
			}
		}
		d.Set("tags", tags)

		// Setting volumes
		volumes := make(map[string]interface{})
		if instanceDetails.Volumes != nil {
			instanceVolumes := instanceDetails.Volumes
			volumeList := *instanceVolumes
			for i := 0; i < len(volumeList); i++ {
				volume := volumeList[i]
				volumeName := volume["name"]
				volumes[volumeName.(string)] = volume["value"]
			}
		}
		d.Set("volumes", volumes)
	}
}
