
# Terraform Provider Lumen Technologies
This is the Terraform Provider for Lumen Bare Metal Management Platform. It interfaces with the Lumen API client. Developed in Go.

- Lumen API doc: (http://apidocs.edge.lumen.com/)


## Requirements
- [Terraform](https://www.terraform.io/downloads) >= 1.1.8
- [Golang](https://go.dev/doc/install) >= 1.18.1

## Getting Started
To get started follow the guidelines provided in [getting started](./docs/guides/getting_started.md)

## Supported Data Sources by Lumen Terraform Provider for Bare Metal

The following list of data sources are supported by Lumen Terraform Provider

| Data Source Name | Description |
|------------------|-------------|
| [data_source_bare_metal_instances](./docs/data-sources/data_source_bare_metal_instances.md) | Lumen Data Source for listing all bare metal instances |
| [data_source_bare_metal_instance_id](./docs/data-sources/data_source_bare_metal_instance_id.md) | Lumen Data Source for listing bare metal instance(s) based on instance id |
| [data_source_bare_metal_instance_name](./docs/data-sources/data_source_bare_metal_instance_name.md) | Lumen Data Source for listing bare metal instance(s) based on instance name |

## Supported Resources by Lumen Terraform Provider for Bare Metal

The following list of resources are supported by Lumen Terraform Provider.

| Resource Name | Description |
|---------------|-------------|
| [create_resource_bare_metal_instance](./docs/resources/create_resource_bare_metal_instance.md) | Create Lumen resource bare metal instance |
| [delete_resource_bare_metal_instance](./docs/resources/delete_resource_bare_metal_instance.md) | Delete Lumen resource bare metal instance |
| [update_resource_bare_metal_instance](./docs/resources/update_resource_bare_metal_instance.md) | Update Lumen resource bare metal instance |

