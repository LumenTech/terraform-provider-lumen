| Page_Title                       | Description                                       |
|----------------------------------|---------------------------------------------------|
| Data_Source_Bare_Metal_Locations | Provides a list of locations for Lumen bare metal |

## Introduction
This document provides usage details on a data source that can be used to access a list of locations for Lumen bare metal.

## Example Usage
`main.tf`
```hcl
data "lumen_bare_metal_locations" "locations" {}

output "locations" {
  value = data.lumen_bare_metal_locations.locations
  description = "Lumen bare metal locations"
}
```

## Schema

### Computed
- id (String) "The location id"
- name (String) "The name of the location"
- status (String) "The status of the location"
- region (String) "The region the location is in"

## Terraform Input Variables
### Required
- username "Lumen username"
- password "Lumen password"
- account_number "Customer Account Number"

### Example usage
`terraform.tfvars`
```hcl
username = $lumen_username
password = $lumen_password
account_number = $account_number
```
