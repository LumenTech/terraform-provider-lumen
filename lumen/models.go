package lumen

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func toInt64(i interface{}) int64 {
	return stringToInt64(i.(string))
}

func stringToInt64(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

func intToString(n int) string {
	return strconv.FormatInt(int64(n), 10)
}

func int64ToString(n int64) string {
	return strconv.FormatInt(n, 10)
}

// Helper function to parse environment
// variables. This will be used to parse
// user provided environment variables.
func parseEnvironmentVariables(variables []interface{}) []map[string]interface{} {
	var evars []map[string]interface{}
	// iterate over the array of evars
	for i := 0; i < len(variables); i++ {
		row := make(map[string]interface{})
		evarconfig := variables[i].(map[string]interface{})
		for k, v := range evarconfig {
			switch k {
			case "name":
				row["name"] = v.(string)
			case "value":
				row["value"] = v.(string)
			case "export":
				row["export"] = v.(bool)
			case "masked":
				row["masked"] = v
			}
		}
		evars = append(evars, row)
	}
	return evars
}

// Helper function to parse network
// interface data. This will be used
// to parse user provided network
// interface variables.
func parseNetworkInterfaces(interfaces []interface{}) []map[string]interface{} {
	var networkInterfaces []map[string]interface{}
	for i := 0; i < len(interfaces); i++ {
		row := make(map[string]interface{})
		item := (interfaces)[i].(map[string]interface{})
		if item["network_id"] != nil {
			row["network"] = map[string]interface{}{
				"id": fmt.Sprintf("network-%d", item["network_id"].(int)),
			}
		}
		if item["ip_address"] != nil {
			row["ipAddress"] = item["ip_address"]
		}
		if item["ip_mode"] != nil {
			row["ipMode"] = item["ip_mode"]
		}
		if item["network_interface_type_id"] != nil {
			row["networkInterfaceTypeId"] = item["network_interface_type_id"]
		}
		networkInterfaces = append(networkInterfaces, row)
	}
	return networkInterfaces
}

// Helper function to parse storage
// volume data. This will be used
// to parse user provided storage
// volume variables.
func parseStorageVolumes(volumes []interface{}) []map[string]interface{} {
	var storageVolumes []map[string]interface{}
	for i := 0; i < len(volumes); i++ {
		row := make(map[string]interface{})
		item := (volumes)[i].(map[string]interface{})
		if item["root"] != nil {
			row["root_volume"] = item["root"]
		}
		if item["name"] != nil {
			row["name"] = item["name"]
		}
		if item["size"] != nil {
			row["size"] = item["size"]
		}
		if item["size_id"] != nil {
			row["sizeId"] = item["size_id"]
		}
		if item["storage_type"] != nil {
			row["storageType"] = item["storage_type"]
		}
		if item["datastore_id"] != nil {
			row["datastoreId"] = item["datastore_id"]
		}
		storageVolumes = append(storageVolumes, row)
	}
	return storageVolumes
}

// Populate custom configs in network resource schema.
func SetNetworkInstanceCustomConfigs(
	instanceDetails *Instance,
	d *schema.ResourceData) {
	customOptions := instanceDetails.Config["customOptions"]
	v := reflect.ValueOf(customOptions)
	if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			strct := v.MapIndex(key)
			if key.Interface().(string) == "transactionId" {
				d.Set("transaction_id", strct.Interface().(string))
			} else if key.Interface().(string) == "cidr" {
				d.Set("network_cidr", strct.Interface().(string))
			} else if key.Interface().(string) == "network" {
				if _, ok := strct.Interface().(interface{}); ok {
					network := strct.Interface().(interface{})
					networkId := reflect.ValueOf(network)
					if networkId.Kind() == reflect.Map {
						for _, nwkey := range networkId.MapKeys() {
							nwstrct := networkId.MapIndex(nwkey)
							if nwkey.Interface().(string) == "id" {
								d.Set("network_resource_id", nwstrct.Interface().(float64))
							}
						}
					}
				}
			}
		}
	}
}

// Populate timestamps in network resource schema.
func SetNetworkInstanceTimestamps(
	instanceDetails *Instance,
	d *schema.ResourceData) {
	d.Set("date_created", instanceDetails.DateCreated)
	d.Set("last_updated", instanceDetails.LastUpdated)
}

// Helper function to populate instance owner and creator
func SetNetworkInstanceUsers(
	instanceDetails *Instance,
	d *schema.ResourceData) {
	valueCreatedBy := reflect.ValueOf(instanceDetails.CreatedBy)
	if valueCreatedBy.Kind() == reflect.Map {
		for _, key := range valueCreatedBy.MapKeys() {
			strct := valueCreatedBy.MapIndex(key)
			if key.Interface().(string) == "username" {
				d.Set("instance_created_by", strct.Interface().(string))
			}
		}
	}
}

// Helper function to populate custom configs in schema
func SetBareMetalInstanceCustomConfigs(
	instanceDetails *Instance,
	d *schema.ResourceData) {
	customOptions := instanceDetails.Config["customOptions"]
	v := reflect.ValueOf(customOptions)
	if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			strct := v.MapIndex(key)
			if key.Interface().(string) == "edgeLocation" {
				d.Set("location", strct.Interface().(string))
			} else if key.Interface().(string) == "networkId" {
				d.Set("network_id", strct.Interface().(string))
			} else if key.Interface().(string) == "centuryLinkNetworkType" {
				d.Set("network_type", strct.Interface().(string))
			}
		}
	}
}

// Helper function to populate connection info in schmea
func SetBareMetalInstanceConnectionInfo(
	instanceDetails *Instance,
	d *schema.ResourceData) {
	connectionInfo := instanceDetails.ConnectionInfo
	data := reflect.ValueOf(connectionInfo[0])
	if data.Kind() == reflect.Map {
		for _, key := range data.MapKeys() {
			strct := data.MapIndex(key)
			if key.Interface().(string) == "ip" {
				d.Set("instance_ip", strct.Interface().(string))
				break
			}
		}
	}
}

// Helper function to populate timestamps in schema
func SetBareMetalInstanceTimestamps(
	instanceDetails *Instance,
	d *schema.ResourceData) {
	d.Set("date_created", instanceDetails.DateCreated)
	d.Set("last_updated", instanceDetails.LastUpdated)
}

// Helper function to populate instance owner and creator
func SetBareMetalInstanceUsers(
	instanceDetails *Instance,
	d *schema.ResourceData) {
	// Setting instance creator
	valueCreatedBy := reflect.ValueOf(instanceDetails.CreatedBy)
	if valueCreatedBy.Kind() == reflect.Map {
		for _, key := range valueCreatedBy.MapKeys() {
			strct := valueCreatedBy.MapIndex(key)
			if key.Interface().(string) == "username" {
				d.Set("instance_created_by", strct.Interface().(string))
			}
		}
	}

	// Setting instance owner
	valueOwner := reflect.ValueOf(instanceDetails.CreatedBy)
	if valueOwner.Kind() == reflect.Map {
		for _, key := range valueOwner.MapKeys() {
			strct := valueOwner.MapIndex(key)
			if key.Interface().(string) == "username" {
				d.Set("instance_owner", strct.Interface().(string))
			}
		}
	}
}

// Helper function to populate volumes in schema
func SetBareMetalInstanceVolumes(
	instanceDetails *Instance,
	d *schema.ResourceData) {
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

// Helper function to populate tags in schema
func SetBareMetalInstanceTags(
	instanceDetails *Instance,
	d *schema.ResourceData) {
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
}
