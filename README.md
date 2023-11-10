# Lumen Technologies Terraform Provider
This is the Terraform Provider for Lumen's Bare Metal Management Platform. It interfaces with the Lumen API client. Developed in Go.

## Version
Current: 1.0.0

## Requirements
- [Golang](https://go.dev/doc/install) >= 1.19
- [Terraform](https://www.terraform.io/downloads) >= 1.2.3

## Clone Repo
Clone terraform-provider-lumen from github:
```shell
git clone https://github.com/CenturyLink/terraform-provider-lumen.git
```

## Install Terraform Provider
First, the [developer overrides](https://developer.hashicorp.com/terraform/cli/config/config-file) will need to be set. Add the following to your terraform rc file:
```
provider_installation {
  dev_overrides {
    "LumenTech/lumen" = "<Insert path to Go's bin>"
  }
}
```
Next, install the terraform provider:
```shell
cd terraform-provider-lumen
go install
```

This will install the libraries needed for Lumen's terraform provider to execute API operations.


## Supported Data Sources

The following list of data sources are supported by Lumen's Terraform Provider:

| Data Source Name                                                                                      | Description |
|-------------------------------------------------------------------------------------------------------|-------------|
| [data_source_bare_metal_locations](./docs/data-sources/data_source_bare_metal_locations.md)           | Provides a list of locations for Lumen bare metal |
| [data_source_bare_metal_configurations](./docs/data-sources/data_source_bare_metal_configurations.md) | Provides a list of Lumen bare metal configurations at a specific location |
| [data_source_bare_metal_osImages](./docs/data-sources/data_source_bare_metal_os_images.md)             | Provides a list of available OS images at a specific location |
| [data_source_bare_metal_networkSizes](./docs/data-sources/data_source_bare_metal_network_sizes.md)     | Provides a list of Lumen network sizes at a specific location |

## Deprecated Data Sources
These data sources use the older backend API. Once you have migrated to the new version of Lumen Edge Bare Metal API, the below data sources are no longer supported.

| Data Source Name | Description |
|------------------|-------------|
| [data_source_bare_metal_instances](./docs/data-sources/data_source_bare_metal_instances.md) | Lumen Data Source for listing all bare metal instances |
| [data_source_bare_metal_instance_id](./docs/data-sources/data_source_bare_metal_instance_id.md) | Lumen Data Source for listing bare metal instance(s) based on instance id |
| [data_source_bare_metal_instance_name](./docs/data-sources/data_source_bare_metal_instance_name.md) | Lumen Data Source for listing bare metal instance(s) based on instance name |
| [data_source_network_instances](./docs/data-sources/data_source_network_instances.md) | Lumen Data Source for listing network instance(s) |
| [data_source_network_instance_id](./docs/data-sources/data_source_network_instance_id.md) | Lumen Data Source for listing network instance(s) based on instance id |
| [data_source_network_instance_name](./docs/data-sources/data_source_network_instance_name.md) | Lumen Data Source for listing network instance(s) based on instance name |

## Supported Resources

The following list of resources are supported by Lumen's Terraform Provider:

| Resource Name                                                                  | Description                          |
|--------------------------------------------------------------------------------|--------------------------------------|
| [resource_bare_metal_server](./docs/resources/resource_bare_metal_server.md)   | Used for bare metal server creation  |
| [resource_bare_metal_network](./docs/resources/resource_bare_metal_network.md) | Used for bare metal network creation |

## Deprecated Resources
These resources use the older backend API. Once you have migrated to the new version of Lumen Edge Bare Metal API, the below resources are no longer supported.

| Resource Name | Description |
|---------------|-------------|
| [create_resource_bare_metal_instance](./docs/resources/resource_order_create_bare_metal_instance.md) | Create Lumen resource bare metal instance |
| [delete_resource_bare_metal_instance](./docs/resources/resource_order_delete_bare_metal_instance.md) | Delete Lumen resource bare metal instance |
| [update_resource_bare_metal_instance](./docs/resources/resource_order_update_bare_metal_instance.md) | Update Lumen resource bare metal instance |
| [create_resource_network_instance](./docs/resources/resource_order_create_network_instance.md) | Create Lumen resource network instance |
| [delete_resource_network_instance](./docs/resources/resource_order_delete_network_instance.md) | Delete Lumen resource network instance |
