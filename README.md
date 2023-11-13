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

| Data Source Name                                                                                      | Description                                                               |
|-------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------|
| [data_source_bare_metal_locations](./docs/data-sources/data_source_bare_metal_locations.md)           | Provides a list of locations for Lumen bare metal                         |
| [data_source_bare_metal_configurations](./docs/data-sources/data_source_bare_metal_configurations.md) | Provides a list of Lumen bare metal configurations at a specific location |
| [data_source_bare_metal_osImages](./docs/data-sources/data_source_bare_metal_os_images.md)            | Provides a list of available OS images at a specific location             |
| [data_source_bare_metal_networkSizes](./docs/data-sources/data_source_bare_metal_network_sizes.md)    | Provides a list of Lumen network sizes at a specific location             |

## Supported Resources

The following list of resources are supported by Lumen's Terraform Provider:

| Resource Name                                                                  | Description                          |
|--------------------------------------------------------------------------------|--------------------------------------|
| [resource_bare_metal_server](./docs/resources/resource_bare_metal_server.md)   | Used for bare metal server creation  |
| [resource_bare_metal_network](./docs/resources/resource_bare_metal_network.md) | Used for bare metal network creation |
