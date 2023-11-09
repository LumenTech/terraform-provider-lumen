| Page_Title      | Description                           |
|-----------------|---------------------------------------|
| Lumen Terraform Provider  | Details on Lumen's Terraform provider |

## Introduction
This document outlines details on Lumen's Terraform provider. It walks through details on provider schema, and lumen provider data_source_order and resource_order that is offered by lumen terraform provider. Also details related to how to use lumen provider is mentioned.

## Schema

### Required
- username (String) "Lumen API username for authentication"
- password (String) "Password of Lumen API user for authentication"

### Optional
- account_number (String) "Account number for this Lumen account"
- api_access_token (String) "Deprecated - Access Token of Lumen API user, instead of authenticating with username and password"
- api_refresh_token (String) "Deprecated - Refresh Token of Lumen API user"

## Data Sources
```golang
DataSourcesMap: {
    "lumen_bare_metal_configurations": DataSourceBareMetalConfigurations(),
    "lumen_bare_metal_locations":      DataSourceBareMetalLocations(),
    "lumen_bare_metal_network_sizes":  DataSourceBareMetalNetworkSizes(),
    "lumen_bare_metal_os_images":      DataSourceBareMetalOsImages(),
    // Deprecated Data Sources (below)
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
    "lumen_bare_metal_server":  ResourceBareMetalServer(), 
    "lumen_bare_metal_network": ResourceBareMetalNetwork(), 
    // Deprecated Resources (below) 
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
      source = "LumenTech/lumen"
      version = "1.0.0"
    }
  }
}

provider "lumen" {
  # Configuration options
  username = var.username
  password = var.password
  account_number = var.account_number
  api_access_token = var.api_access_token
  api_refresh_token = var.api_refresh_token
}
```

`variables.tf`
```hcl
- "username" : $consumer_key
- "password" : $consumer_secret
- "account_number": $lumen_account_number
- "api_access_token": $lumen_api_access_token
- "api_refresh_token": $lumen_api_refresh_token
```
