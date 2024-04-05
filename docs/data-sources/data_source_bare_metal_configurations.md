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
- machine_count (Integer) "The number of machines in this configuration"
- price (String) "The price for this configuration"
- tier (String) "The tier associated with the configuration which is used for mapping to OS Image pricing"

## Terraform Input Variables
### Required
- consumer_key "Lumen consumer_key"
- consumer_secret "Lumen consumer_secret"
- account_number "Customer Account Number"
- location_id "The id of a location"
