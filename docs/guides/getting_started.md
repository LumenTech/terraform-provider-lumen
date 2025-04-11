## Lumen Provider Operations

Lumen's terraform provider supports the following data sources and resources:
- data_sources:
  - `data_source_bare_metal_locations`:  Provides a list of locations for Lumen bare metal
  - `data_source_bare_metal_configurations`: Provides a list of Lumen bare metal configurations at a specific location
  - `data_source_bare_metal_os_images`: Provides a list of available OS images at a specific location
  - `data_source_bare_metal_network_sizes`: Provides a list of Lumen network sizes at a specific location
- resources:
  - `resource_bare_metal_server`: Used for bare metal server creation
  - `resource_bare_metal_network`: Used for bare metal network creation

## Deprecation Notice
Lumen Bare Metal API will be deprecating v2 endpoints which are directly tied with the version 2.X.X of the terraform provider.
There will be a new v3 Lumen Bare Metal API which will be made available, and it will use a new backend authentication.
We will be releasing a 3.X.X version of the terraform provider which will target the new API version.  In order to upgrade
to the 3.X.X version of terraform provider you will need to generate a new consumer_key and consumer_secret see 
[authentication](./authentication.md).  As part of this deprecation process you will see warning messages when running 
the provider once we have finalized on a Sunset date.

Please use the correct provider based on the Consumer Key and Consumer Secret.
Terraform 2.X.X targets Lumen Bare Metal API v2.
Terraform 3.X.X will target Lumen Bare Metal API v3 with different authentication.

## Deprecated - Generate API Key for Lumen Provider
Lumen Bare Metal users on the older backend API need to generate API token for authentication to use Lumen resources. 
Once you have migrated to the new version of Lumen Edge Bare Metal API, an api token is no longer needed.
Details on how to generate tokens is provided in [authentication](./authentication.md) guide.

## Enable logging
Terraform logs can be enabled by setting the TF_LOG environment variable to get detailed logs to appear on stderr. Log level options are: TRACE, DEBUG, INFO, WARN or ERROR to change the verbosity of the logs. To persist logged output, set the TF_LOG_PATH variable to force the log to always be appended to a specific file. The below lines can be added to ~/.bashrc.
```shell
# Terraform log settings
export TF_LOG=TRACE
export TF_LOG_PATH="/<path_to_terraform-provider-lumen>/client/logs/terraform.log"
```