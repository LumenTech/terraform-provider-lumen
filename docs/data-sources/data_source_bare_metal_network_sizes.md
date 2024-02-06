| Page_Title                          | Description                                                   |
|-------------------------------------|---------------------------------------------------------------|
| Data_Source_Bare_Metal_NetworkSizes | Provides a list of Lumen network sizes at a specific location |

## Introduction
This document provides usage details on a data source that can be used to access a list of Lumen network sizes at a specific location.

## Example Usage
`main.tf`
```hcl
data "lumen_bare_metal_network_sizes" "network_sizes" {
  location_id = var.location_id
}

output "network_sizes" {
  value = data.lumen_bare_metal_network_sizes.network_sizes
  description = "Lumen bare metal network sizes"
}
```

## Schema

### Required
- location_id (String) "The id of a location"

### Computed
- id (String) "The id of a network size"
- name (String) "The name of this network size"
- cidr (String) "The CIDR for this network size"
- network_type (String) "The type of network being used"
- available_ips (Integer) "The number of available IPs for this network size"
- price (String) "The price for this network size"

## Terraform Input Variables
### Required
- username "Lumen username"
- password "Lumen password"
- account_number "Customer Account Number"
- location_id "The id of a location"

### Example usage
`terraform.tfvars`
```hcl
username = $lumen_username
password = $lumen_password
account_number = $account_number
location_id = $location_id
```
