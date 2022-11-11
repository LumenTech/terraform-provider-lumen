
# Terraform Provider Lumen Technologies
This is the Terraform Provider for Lumen Bare Metal Management Platform. It interfaces with the Lumen API client. Developed in Go.

## Version
Current: 0.5.0

## Requirements
- [Terraform](https://www.terraform.io/downloads) >= 1.2.3
- [Golang](https://go.dev/doc/install) >= 1.18.3

## Getting Started
To get started follow the guidelines provided in [getting started](./docs/guides/getting_started.md)

## Supported Data Sources by Lumen Terraform Provider.

The following list of data sources are supported by Lumen Terraform Provider:

| Data Source Name | Description |
|------------------|-------------|
| [data_source_bare_metal_instances](./docs/data-sources/data_source_bare_metal_instances.md) | Lumen Data Source for listing all bare metal instances |
| [data_source_bare_metal_instance_id](./docs/data-sources/data_source_bare_metal_instance_id.md) | Lumen Data Source for listing bare metal instance(s) based on instance id |
| [data_source_bare_metal_instance_name](./docs/data-sources/data_source_bare_metal_instance_name.md) | Lumen Data Source for listing bare metal instance(s) based on instance name |
| [data_source_network_instances](./docs/data-sources/data_source_network_instances.md) | Lumen Data Source for listing network instance(s) |
| [data_source_network_instance_id](./docs/data-sources/data_source_network_instance_id.md) | Lumen Data Source for listing network instance(s) based on instance id |
| [data_source_network_instance_name](./docs/data-sources/data_source_network_instance_name.md) | Lumen Data Source for listing network instance(s) based on instance name |

## Supported Resources by Lumen Terraform Provider.

The following list of resources are supported by Lumen Terraform Provider:

| Resource Name | Description |
|---------------|-------------|
| [create_resource_bare_metal_instance](./docs/resources/resource_order_create_bare_metal_instance.md) | Create Lumen resource bare metal instance |
| [delete_resource_bare_metal_instance](./docs/resources/resource_order_delete_bare_metal_instance.md) | Delete Lumen resource bare metal instance |
| [update_resource_bare_metal_instance](./docs/resources/resource_order_update_bare_metal_instance.md) | Update Lumen resource bare metal instance |
| [create_resource_network_instance](./docs/resources/resource_order_create_network_instance.md) | Create Lumen resource network instance |
| [delete_resource_network_instance](./docs/resources/resource_order_delete_network_instance.md) | Delete Lumen resource network instance |
