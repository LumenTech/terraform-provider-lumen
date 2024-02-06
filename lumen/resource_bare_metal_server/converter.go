package resource_bare_metal_server

import "terraform-provider-lumen/lumen/client/model/bare_metal"

func convertDataToAttachedNetworks(items []interface{}) []bare_metal.AttachNetwork {
	networks := make([]bare_metal.AttachNetwork, len(items))
	for i, item := range items {
		net := item.(map[string]interface{})
		networkID := net["network_id"].(string)
		assignIPV6 := net["assign_ipv6_address"].(bool)
		networks[i] = bare_metal.AttachNetwork{
			NetworkID:  networkID,
			AssignIPV6: assignIPV6,
		}
	}
	return networks
}

func convertNetworksToListOfAttachNetworks(networks []bare_metal.ServerNetwork) []bare_metal.AttachNetwork {
	attachNetworks := make([]bare_metal.AttachNetwork, len(networks))
	for idx, n := range networks {
		attachNetworks[idx] = bare_metal.AttachNetwork{
			NetworkID:  n.NetworkID,
			AssignIPV6: n.AssignIPV6Address,
		}
	}
	return attachNetworks
}
