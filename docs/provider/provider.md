| Page_Title               | Description                         |
|--------------------------|-------------------------------------|
| Lumen Terraform Provider | Details on Lumen Terraform provider |

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
    "lumen_network_instances":        DataSourceNetworkAllInstances(),
		"lumen_network_instance_id":      DataSourceNetworkInstanceId(),
		"lumen_network_instance_name":    DataSourceNetworkInstanceName(),
},
```
Details on data-sources are provided in [docs](../data-sources).

# Resources
```golang
ResourcesMap: {
    "lumen_bare_metal_instance": ResourceBareMetalInstance(),
    "lumen_network_instance":    ResourceNetworkInstance(),
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
      source = "lumen.com/lumentech/lumen"
      version = "0.4.1"
    }
  }
}

# Provider access creds
provider "lumen" {
  username = var.lumen_username
  password = var.lumen_password
  api_access_token = var.lumen_api_access_token
  api_refresh_token = var.lumen_api_refresh_token
}
```

`variables.tf`
```hcl
- "username" : $consumer_key
- "password" : $consumer_secret
- "lumen_api_access_token": $lumen_api_access_token
- "lumen_api_refresh_token": $lumen_api_refresh_token
```

`terraform.tfvars`
```hcl
lumen_username = $consumer_key
lumen_password = $consumer_secret
lumen_api_access_token = $lumen_api_access_token
lumen_api_refresh_token = $lumen_api_refresh_token
```
