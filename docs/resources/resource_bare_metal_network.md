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

output "network" {
  value = lumen_bare_metal_network.network
  description = "Lumen bare metal network details"
}
```

## Schema

### Required
- name (String) "Network name"
- location_id (String) "A location id (can be retrieved with data_source_bare_metal_locations)"
- network_size_id (String) "A network size id (can be retrieved with data_source_bare_metal_networkSizes)"

### Computed
- id (String)
- account_id (String)
- service_id (String)
- location (String)
- ip_block (String)
- gateway (String)
- available_ips (Integer)
- total_ips (Integer)
- type (String)
- status (String)
- prices (List of Price)
- created (String)
- updated (String)