| Page_Title      | Description                                 |
|-----------------|---------------------------------------------|
| Data_Source_Bare_Metal_Instances  | Details on bare metal instances for a particular tenant        |

## Introduction
This document provides Lumen bare metal instance(s) detail for a specific tenant. The API details are provided in Ref [[1]](#1).

## Example Usage
`main.tf`
```hcl
data "lumen_bare_metal_instances" "all" {}

output "all_instances" {
  value = data.lumen_bare_metal_instances.all.instances
}
```

## Schema

### Computed
- id (Integer) "The IDs of instances created"
- name (String) "The names of instances created"
- description (String) "The instance description"
- cloud_id (Integer) "The ID of the cloud associated with the instance"
- group_id (Integer) "The ID of the group associated with the instance"
- instance_type_id (Integer) "The type of instance to provision"
- instance_layout_id (Integer) "The layout to provision the instance from"
- plan_id (Integer) "The service plan associated with the instance"
- resource_pool_id (Integer) "The ID of the resource pool to provision the instance to"
- environment (String) "The environment to assign the instance to"
- version (String)
- status (String) "Instance status"
- instance_location (String) "The instance location"
- instance_ip (String) "The instance ip address"

## Terraform Input Variables
`variables.tf`
### Required
- lumen_api_url "Lumen API endpoint"
- lumen_access_token "Lumen user token (user should have API access)"
### Optional
- lumen_username "Lumen username (user should have API access)"
- lumen_password "Lumen password"

Each of the variables are defined in `terraform.tfvars`.

### Example usage
`terraform.tfvars` 
```hcl
lumen_api_url = "https://api.lumen.com/EdgeServices/v1/Compute/api/instances"
lumen_username = null
lumen_password = null
lumen_access_token = "0000-0000-0000-0000"
```


## References
<a id="1">[1]</a> API doc: http://apidocs.edge.lumen.com/#get-all-instances
