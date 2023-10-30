| Page_Title                            | Description                                                               |
|---------------------------------------|---------------------------------------------------------------------------|
| Data_Source_Bare_Metal_Configurations | Provides a list of Lumen bare metal configurations at a specific location |

## Introduction
This document provides usage details on a data source that can be used to access a list of Lumen bare metal configurations at a specific location.

## Example Usage
`main.tf`
```hcl
data "lumen_bare_metal_configurations" "configurations" {
  location_id = var.location_id
}

output "configurations" {
  value = data.lumen_bare_metal_configurations.configurations
  description = "Lumen bare metal configurations"
}
```

## Schema

### Required
- location_id (String) "The id of a location"

### Computed
- name (String) "The type of configuration (ie small, medium, large)"
- cores (Integer) "The number of cores in this configuration"
- memory (String) "The memory available for this configuration"
- storage (String) "The storage available for this configuration"
- disks (Integer) "The number of disks in this configuration"
- nics (Integer) "The number of NICs in this configuration"
- processors (Integer) "The number of processors in this configuration"
- machineCount (Integer) "The number of machines in this configuration"
- price (String) "The price for this configuration"

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
