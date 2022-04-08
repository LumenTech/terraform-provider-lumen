| Page_Title      | Description                                 |
|-----------------|---------------------------------------------|
| Lumen Terraform Provider  | Details on Lumen Terraform provider  |

## Introduction
This document outlines details on lumen terraform provider. It walks through details on provider schema, and lumen provider data_source_order and resource_order that is offered by lumen terraform provider. Also details related to how to use lumen provider is mentioned.

## Schema

### Required
- url (String) "Lumen API endpoint URL where requests will be directed"
- access_token (String) "Access Token of Lumen API user, instead of authenticating with username and password"

### Optional
- username (String) "Lumen API username for authentication"
- password (String) "Password of Lumen API user for authentication",

## Data Sources
```golang
DataSourcesMap: {
    "lumen_bare_metal_instances":     DataSourceBareMetalAllInstances(),
    "lumen_bare_metal_instance_id":   DataSourceBareMetalInstanceId(),
    "lumen_bare_metal_instance_name": DataSourceBareMetalInstanceName(),
},
```
Details on data-sources are provided in [docs](../data-sources).

# Resources
```golang
ResourcesMap: {
    "lumen_bare_metal_instance": ResourceBareMetalInstance(),
},
```
Details on resources are provided in [docs](../resources).

## Example Usage

### Provider configuration
`provider.tf`
```hcl
# Provider configuration
terraform {
  required_providers {
    lumen = {
      source = "lumen.com/edge/lumen"
      version = "0.3.0"
    }
  }
}

# Provider access creds
provider "lumen" {
  url = var.lumen_api_url
  username = null
  password = null
  access_token = var.lumen_access_token
}
```
`variables.tf`
```hcl
- "lumen_api_url" : "https://api.lumen.com/EdgeServices/v1/Compute/api/"
- "lumen_access_token" : "0000-0000-0000-0000"
- "lumen_username" : "Lumen username (user should have API access, optional)"
- "lumen_password" : "Lumen password" (optional)
```

`terraform.tfvars`
```hcl
lumen_api_url = "https://api.lumen.com/"
lumen_username = null
lumen_password = null
lumen_access_token = "0000-0000-0000-0000"
```
