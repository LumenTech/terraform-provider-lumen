| Page_Title                  | Description                            |
|-----------------------------|----------------------------------------|
| resource_bare_metal_network | Details on bare metal network creation |

## Introduction
This document provides details on the resource to create Lumen network(s) associated with the bare metal product.

## Example Usage
`main.tf`
```hcl
resource "lumen_bare_metal_network" "network" {
  #  Example request data
  name = "testNetwork5"
  location_id = "DNVTCO56LEC"
  network_size_id = "6529723924b8bf31ebd998e2"
}

resource "lumen_bare_metal_network" "network2" {
  #  Create a private network with an existing VRF
  name = "testNetwork6"
  location_id = "DNVTCO56LEC"
  network_type = "PRIVATE"
  vrf = "88/VP12/003434/ASRT"
}

resource "lumen_bare_metal_network" "network3" {
  #  Create a private network with a new VRF
  name = "testNetwork7"
  location_id = "DNVTCO56LEC"
  network_type = "PRIVATE"
  vrf_description = "testPrivateNetwork"
}

output "network" {
  value = lumen_bare_metal_network.network
  description = "Lumen bare metal network details"
}

output "network2" {
  value = lumen_bare_metal_network.network2
  description = "Private network with existing VRF details"
}

output "network3" {
  value = lumen_bare_metal_network.network3
  description = "Private network with new VRF details"
}
```

## Schema

### Required
- name (String) "Network name (updatable)"
- location_id (String) "A location id (can be retrieved with data_source_bare_metal_locations)"
- network_size_id (String) "A network size id (can be retrieved with data_source_bare_metal_networkSizes)"

### Optional
- network_type (String) "The type of network being used. Three possible values: INTERNET, DUAL_STACK_INTERNET and PRIVATE"
- vrf (String) "For private networks, this is an existing VRF to be used in creating the new network."
- vrf_description (String) "For private networks, create a new VRF with this description and use it in creating the new network."

### Computed
- id (String)
- account_id (String)
- service_id (String)
- location (String)
- ip_block (String)
- ipv6_block (String)
- gateway (String)
- available_ips (Integer)
- total_ips (Integer)
- type (String)
- status (String)
- vrf_value (String)
- vrf_description_value (String)
- prices (List of Price)
- created (String)
- updated (String)