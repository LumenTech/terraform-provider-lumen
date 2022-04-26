package lumen

import (
	"fmt"
	"strconv"
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

// helper function to parse environment variables
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

// helper function to parse network interface data
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

// helper function to parse storage volume data
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
