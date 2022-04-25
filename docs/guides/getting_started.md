## Prerequisites
Download and Install:
- [Golang](https://go.dev/doc/install) >= 1.18.1
- [Terraform](https://www.terraform.io/downloads) >= 1.1.8

## Installation
Clone terraform-provider-lumen-technologies from github:
```shell
$ git clone https://github.com/LumenTech/terraform-provider-lumen-technologies.git
```

## Setup and install Lumen Terraform Provider
Next, setup and install terraform provider lumen technologies.
```shell
$ cd terrform-provider-lumen-technologies
$ make setup
$ make install
```
This will install Morpheus SDKv2 and other SDKs for lumen terraform provider to run and execute opertations.

## Enable logging
Terraform logs can be enabled by setting the TF_LOG environment variable to any value, to get detailed logs to appear on stderr. TF_LOG to one of the log levels TRACE, DEBUG, INFO, WARN or ERROR to change the verbosity of the logs. To persist logged output set the TF_LOG_PATH variable to force the log to always be appended to a specific file when logging is enabled. For logging add the below lines in ~/.bashrc, close the working sessions and open a new session.
```shell
# Terraform log settings
export TF_LOG=TRACE
export TF_LOG_PATH="/$path_to_terraform-provider-lumen-technologies/client/logs/terraform.log"
```

## Generate API key for Lumen Provider
Lumen Bare Metal users need to generate API token for authentication to use Lumen resources. Details on how to generate tokens is provided in [authentication](./authentication.md) guide.

## Lumen provider operations
At this point you are ready to perform provider operations. Lumen Terraform provider currently supports the following data-sources and resource-order operations.
- data_sources:
    - `data_source_bare_metal_instances`: To list all instances under a specific tenant
    - `data_source_bare_metal_instance_id`: To get instance details based on instance id
    - `data_source_bare_metal_instance_name`: To get instance details based on instance name
- resource_order:
    - `create_resource_bare_metal_instance`: To build bare metal instance.
    - `update_resource_bare_metal_instance`: To update bare metal instance.
    - `delete_resource_bare_metal_instance`: To delete bare metal instance.

Additional details for each operation, related to `schema`, `input_variables`, `output` are provided in [Lumen_Provider_Operations](../indexes.md). Depending upon use case, example scripts are provided in the repo under `examples` directory. Copy over terraform utility scripts `main.tf`, `variables.tf`, `provider.tf` and `terraform.tfvars` to client. For example to get list of instances copy over terraform files from `exmaples/data_source_bare_metal_instances` as mentioned below:
```shell
$ cd ~/terrform-provider-lumen-technologies/client/
# Copy terraform configuration files
$ cp ../examples/data_source_bare_metal_instances/*.tf* .

# Copy lumen terraform provider
$ cp ../examples/provider/*.tf* .
```

Add lumen API end point, api token obtained from above or username and password in `terraform.tfvars` file. Populate variables as required for terraform operation. Additional details are provided in docs directory for each terraform operation. Initialize workspace and execute CLI based on use case as mentioned below. 

### 5.2: To create / read / update instance(s)
```shell
$ terraform init # initializes terraform-provider-lumen workspace
$ terraform plan # displays list of resource(s) to be created / read / updated
$ terraform apply --auto-approve # to create resource(s) as per plan
```
Curently terraform-provider-lumen-technologies supports updating instance name, description, labels. Backend dev work for updation of tags (add/remove/replace) is in place. For updating these instance(s) details, example script is provided. 

### 5.3: To delete instance(s)
```shell
terraform init # initializes lec-terraform-provider and workspace
terraform destroy --auto-approve # to delete the resource shown in terraform plan
```

## Go code format check
Excute:
```shell
make fmtcheck
```

To fix format:
```shell
make fmt
```

## Clean-up
To clean up terraform workspace, excute:
```shell
make clean
```
This will clean up existing terraform lock files and existing terraform plugins. This will not clean the terraform state files or utility files. Execute `make install` for building the provider again and execute `terraform init` to re-enable the workspace.
