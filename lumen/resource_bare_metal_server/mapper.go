package resource_bare_metal_server

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-lumen/lumen/client/model/bare_metal"
)

func populateServerSchema(d *schema.ResourceData, server bare_metal.Server) {
	d.Set("name", server.Name)
	d.Set("location_id", server.LocationID)
	d.Set("configuration_name", server.Configuration.Name)
	d.Set("os_image", server.OSImage)
	d.Set("machine_id", server.MachineID)
	d.Set("machine_name", server.MachineName)
	d.Set("location", server.Location)
	d.Set("configuration_cores", server.Configuration.Cores)
	d.Set("configuration_memory", server.Configuration.Memory)
	d.Set("configuration_storage", server.Configuration.Storage)
	d.Set("configuration_disks", server.Configuration.Disks)
	d.Set("configuration_nics", server.Configuration.NICs)
	d.Set("configuration_processors", server.Configuration.Processors)
	networks := make([]map[string]interface{}, len(server.Networks))
	networkIds := make([]string, len(server.Networks))
	for i, network := range server.Networks {
		networks[i] = map[string]interface{}{
			"id":             network.ID,
			"network_id":     network.NetworkID,
			"network_name":   network.NetworkName,
			"network_type":   network.NetworkType,
			"status":         network.Status,
			"status_message": network.StatusMessage,
			"ip":             network.IP,
			"vlan":           network.VLAN,
		}
		networkIds[i] = network.NetworkID
	}
	d.Set("networks", networks)
	d.Set("network_ids", networkIds)
	d.Set("status", server.Status)
	d.Set("status_message", server.StatusMessage)
	d.Set("boot_disk", server.BootDisk)
	d.Set("service_id", server.ServiceID)
	prices := make([]map[string]interface{}, len(server.Prices))
	for i, price := range server.Prices {
		prices[i] = map[string]interface{}{
			"type":  price.Type,
			"price": price.Price.String(),
		}
	}
	d.Set("prices", prices)
	d.Set("account_id", server.AccountID)
	d.Set("created", server.Created)
	d.Set("updated", server.Updated)
}
