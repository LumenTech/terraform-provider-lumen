| Page_Title                 | Description                           |
|----------------------------|---------------------------------------|
| resource_bare_metal_server | Details on bare metal server creation |

## Introduction
This document provides details on the resource to create Lumen bare metal server(s). In order to create a bare metal
server, you can either provision a new network or use an existing network(s). If you are creating a new network you will
provide the network_size_id and network_name variables. If you are using an existing network you will provide a list of
network IDs in the attach_networks object (the first network will be configured on the server all others will need manual
server configuration changes performed by the user). After the server has been created we support the ability to change
the name, which only affects the data stored in our system, and add/remove networks through altering the attach_networks
field. We currently only manage the configuration of our network infrastructure. Any changes to the networks will require
server configuration changes performed by the user. This is due to us not running an agent or having any access to the
host system. If you provision a server where we created the network if you decide to add networks you will need to
remove the old fields and populate the attach_networks object with your current network id.

## Example Usage
`main.tf`
```hcl
resource "lumen_bare_metal_server" "server" {
  #  Example request data with a new created network
  name = "BASTION01"
  location_id = "DNVTCO56LEC"
  configuration_name = "small_plus"
  os_image_name = "Ubuntu 20.04"
  network_name = "bastion-network"
  network_size_id = "64ef800284b3f203bc27f9e2"
  username = "admin"
  password = "**********"
}

resource "lumen_bare_metal_server" "server2" {
  #  Example request data using multiple existing networks
  name = "BASTION02"
  location_id = "DNVTCO56LEC"
  configuration_name = "small_plus"
  os_image_name = "Ubuntu 20.04"
  attach_networks {
    network_id = "65c283e988f85707cc53b308"
    assign_ipv6_address = true
  }
  attach_networks {
    network_id = "6553dc0d861724132e81ce42"
  }
  attach_networks {
    network_id = "6553dc22861724132e81ce43"
  }
  username = "admin"
  password = "**********"
}

resource "lumen_bare_metal_server" "server3" {
  #  Example server provisioning request with hyperthreading true/false
  name = "BASTION02"
  location_id = "DNVTCO56LEC"
  configuration_name = "small_plus"
  os_image_name = "Ubuntu 20.04"
  network_size_id = "64ef800284b3f203bc27f9e2"
  enable_hyperthreading = true # or false
  username = "admin"
  password = "**********"
}

resource "lumen_bare_metal_server" "server4" {
  #  Example request data with a new private network using a new VRF
  name = "BASTION01"
  location_id = "DNVTCO56LEC"
  configuration_name = "small_plus"
  os_image_name = "Ubuntu 20.04"
  network_name = "bastion-network"
  network_type = "PRIVATE"
  vrf_description = "testVPE"
  username = "admin"
  password = "**********"
}

output "server" {
  sensitive = true
  value = lumen_bare_metal_server.server
  description = "Lumen bare metal server details"
}

output "server2" {
  sensitive = true
  value = lumen_bare_metal_server.server2
  description = "Server with multiple attached networks"
}

output "server3" {
  sensitive = true
  value = lumen_bare_metal_server.server3
  description = "Server with hyperthreading enabled/disabled"
}

output "server4" {
  sensitive = true
  value = lumen_bare_metal_server.server3
  description = "Server with attached private network"
}
```

## Schema

### Required
- name (String) "Server hostname (updatable)"
- location_id (String) "A location id (can be retrieved with data_source_bare_metal_locations)"
- configuration_name (String) "A configuration name (can be retrieved with data_source_bare_metal_configurations)"
- os_image_name (String) "A os image name (can be retrieved with data_source_bare_metal_os_images)"
- username (String) "Username that should be created on the server"

### Conditionally Required
#### attach_networks or (network_name, network_size_id) or (network_name, vrf) or (network_name, vrf_description)
- network_name (String) "The name of the network you wish to create with the server"
- network_size_id (String) "The network size id you wish to create with the server (can be retrieved with data_source_bare_metal_network_sizes)"
- vrf (String) "For private networks, this is an existing VRF to be used in creating the new network."
- vrf_description (String) "For private networks, create a new VRF with this description and use it in creating the new network."
- attach_networks (List of Object) "List of existing networks to attach to the server being provisioned. (updatable)"
#### at least one (password and ssh_public_key)
- password (String) "The password you wish to have associated with the user account"
- ssh_public_key (String) "The ssh public key you wish to have associated with the username"

### Optional
- network_type (String) "The type of network being used. Three possible values: INTERNET, DUAL_STACK_INTERNET, and PRIVATE"
- assign_ipv6_address (Boolean) "A boolean (true/false) value indicating whether to assign an IPv6 address
  for this server if using a dual stack network. Defaults to false if not set."
- enable_hyperthreading (Boolean) "A boolean (true/false) value indicating whether to enable or disable hyperthreading on the server. Two possible values: true and false, defaults to true if not set."

### Computed
- id (String)
- machine_id (String)
- machine_name (String)
- location (String)
- configuration_cores (Integer)
- configuration_memory (String)
- configuration_storage (String)
- configuration_disks (Integer)
- configuration_nics (Integer)
- configuration_processors (Integer)
- networks (List of Network)
- status (String)
- status_message (String)
- boot_disk (String)
- service_id (String)
- prices (List of Price)
- account_id (String)
- created (String)
- updated (String)
- hyperthreading (Boolean)
