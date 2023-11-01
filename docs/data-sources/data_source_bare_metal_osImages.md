| Page_Title                      | Description                                                   |
|---------------------------------|---------------------------------------------------------------|
| Data_Source_Bare_Metal_OsImages | Provides a list of available OS images at a specific location |

## Introduction
This document provides usage details on a data source that can be used to access a list of available OS images at a specific location.

## Example Usage
`main.tf`
```hcl
data "lumen_bare_metal_os_images" "os_images" {
  location_id = var.location_id
}

output "os_images" {
  value = data.lumen_bare_metal_os_images.os_images
  description = "Lumen bare metal OS images"
}
```

## Schema

### Required
- location_id (String) "The id of a location"

### Computed
- name (String) "The name of this OS image"
- price (String) "The price for using this OS image"

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
