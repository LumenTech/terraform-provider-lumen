## Lumen Provider Operations
Lumen's terraform provider supports the following data sources and resources:
- data_sources:
  - `data_source_bare_metal_locations`:  Provides a list of locations for Lumen bare metal
  - `data_source_bare_metal_configurations`: Provides a list of Lumen bare metal configurations at a specific location
  - `data_source_bare_metal_osImages`: Provides a list of available OS images at a specific location
  - `data_source_bare_metal_networkSizes`: Provides a list of Lumen network sizes at a specific location
- resources:
  - `resource_bare_metal_server`: Used for bare metal server creation
  - `resource_bare_metal_network`: Used for bare metal network creation

## Deprecated Provider Operations
The following data sources and resources use the old backend API. Once you have been migrated to the new version of Lumen Edge Bare Metal API, they will no longer be supported.
- data_sources:
  - `data_source_bare_metal_instances`: To list all bare-metal instances under a specific tenant.
  - `data_source_bare_metal_instance_id`: To get bare-metal instance details based on instance id.
  - `data_source_bare_metal_instance_name`: To get bare-metal instance details based on instance name.
  - `data_source_network_instances`: To list all network instances under a specific tenant.
  - `data_source_bare_metal_instance_id`: To get network instance details based on instance id.
  - `data_source_bare_metal_instance_name`: To get network instance details based on instance name.
- resources:
  - `create_resource_bare_metal_instance`: To build bare metal instance.
  - `update_resource_bare_metal_instance`: To update bare metal instance.
  - `delete_resource_bare_metal_instance`: To delete bare metal instance.
  - `create_resource_network_instance`: To build network instance.
  - `delete_resource_network_instance`: To delete network instance.

## Enable logging
Terraform logs can be enabled by setting the TF_LOG environment variable to get detailed logs to appear on stderr. Log level options are: TRACE, DEBUG, INFO, WARN or ERROR to change the verbosity of the logs. To persist logged output, set the TF_LOG_PATH variable to force the log to always be appended to a specific file. The below lines can be added to ~/.bashrc.
```shell
# Terraform log settings
export TF_LOG=TRACE
export TF_LOG_PATH="/<path_to_terraform-provider-lumen>/client/logs/terraform.log"
```

## Deprecated - Generate API Key for Lumen Provider
Lumen Bare Metal users on the older backend API need to generate API token for authentication to use Lumen resources. 
Once you have migrated to the new version of Lumen Edge Bare Metal API, an api token is no longer needed.
Details on how to generate tokens is provided in [authentication](./authentication.md) guide.
