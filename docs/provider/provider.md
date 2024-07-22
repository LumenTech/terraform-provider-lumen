| Page_Title               | Description                           |
|--------------------------|---------------------------------------|
| Lumen Terraform Provider | Details on Lumen's Terraform provider |

## Introduction
This document outlines details on Lumen's Terraform provider. It walks through details on provider schema, and lumen provider data_source_order and resource_order that is offered by lumen terraform provider. Also details related to how to use lumen provider is mentioned.

## Schema

### Required
- consumer_key (String) "Lumen API consumer key for authentication"
- consumer_secret (String) "Lumen API consumer secret of Lumen API user for authentication"

### Optional
- account_number (String) "Account number for this Lumen account"

## Data Sources
```golang
DataSourcesMap: {
    "lumen_bare_metal_configurations": DataSourceBareMetalConfigurations(),
    "lumen_bare_metal_locations":      DataSourceBareMetalLocations(),
    "lumen_bare_metal_network_sizes":  DataSourceBareMetalNetworkSizes(),
    "lumen_bare_metal_os_images":      DataSourceBareMetalOsImages(),
},
```
Details on data-sources are provided in [docs](../data-sources).

# Resources
```golang
ResourcesMap: {
    "lumen_bare_metal_server":  ResourceBareMetalServer(), 
    "lumen_bare_metal_network": ResourceBareMetalNetwork(),
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
      version = ">= 2.4.0"
    }
  }
}

provider "lumen" {
  # Configuration options
  consumer_key = var.consumer_key
  consumer_secret = var.consumer_secret
  account_number = var.account_number
}
```

`variables.tf`
```hcl
- "consumer_key" : $consumer_key
- "consumer_secret" : $consumer_secret
- "account_number": $lumen_account_number
```
