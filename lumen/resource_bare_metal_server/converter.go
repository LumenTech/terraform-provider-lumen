package resource_bare_metal_server

import "terraform-provider-lumen/lumen/client/model/bare_metal"

func convertListOfInterfaceToListOfString(items []interface{}) []string {
	networkIds := make([]string, len(items))
	for idx, item := range items {
		networkIds[idx] = item.(string)
	}
	return networkIds
}

func convertNetworksToListOfNetworkIds(networks []bare_metal.ServerNetwork) []string {
	networkIds := make([]string, len(networks))
	for idx, n := range networks {
		networkIds[idx] = n.NetworkID
	}
	return networkIds
}
