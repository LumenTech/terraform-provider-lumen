| Page_Title                 | Description                           |
|----------------------------|---------------------------------------|
| resource_bare_metal_server | Details on bare metal server creation |

## Introduction
This document provides details on the resource to create Lumen bare metal server(s). In order to create a bare metal server,
you can either provision a new network or use an existing network.  If you are using an existing network you will provide
the id for the network_id variable.  If you are create a new network you will provide the network_size_id and network_name variables.

## Example Usage
`main.tf`
```hcl
resource "lumen_bare_metal_server" "server" {
  #  Example request data
  name = "BASTION01"
  location_id = "DNVTCO56LEC"
  configuration_name = "small_plus"
  os_image_name = "Ubuntu 20.04"
  network_name = "bastion-network"
  network_size_id = "64ef800284b3f203bc27f9e2"
  username = "admin"
  password = "**********"
}

output "server" {
  sensitive = true
  value = lumen_bare_metal_server.server
  description = "Lumen bare metal server details"
}
```

## Schema

### Required
- name (String) "Server hostname"
- location_id (String) "A location id (can be retrieved with data_source_bare_metal_locations)"
- configuration_name (String) "A configuration name (can be retrieved with data_source_bare_metal_configurations)"
- os_image_name (String) "A os image name (can be retrieved with data_source_bare_metal_os_images)"
- username (String) "Username that should be created on the server"


### Optional
#### network_id or (network_name, network_size_id)
- network_id (String) "ID of network if you are using an existing network"
- network_name (String) "The name of the network you wish to create with the server"
- network_size_id (String) "The network size id you wish to create with the server (can be retrieved with data_source_bare_metal_network_sizes)"
#### at least one (password and ssh_public_key)
- password (String) "The password you wish to have associated with the user account"
- ssh_public_key (String) "The ssh public key you wish to have associated with the username"

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