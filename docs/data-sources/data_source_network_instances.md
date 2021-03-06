| Page_Title      | Description                                 |
|-----------------|---------------------------------------------|
| Data_Source_Network_Instances  | Details on network instances for a particular tenant        |

## Introduction
This document provides Lumen network instance(s) detail for a specific tenant. The API details are provided in Ref [[1]](#1).

## Example Usage
`main.tf`
```hcl
data "lumen_network_instances" "all" {}

output "all_instances" {
  value = data.lumen_network_instances.all.instances
}
```

## Schema

### Computed
- id (Integer) "The IDs of instance(s) created"
- name (String) "The names of instance(s) created"
- description (String) "The instance description"
- cloud_id (Integer) "The ID of the cloud associated with the instance"
- group_id (Integer) "The ID of the group associated with the instance"
- instance_type_id (Integer) "The type of instance to provision"
- instance_layout_id (Integer) "The layout to provision the instance from"
- plan_id (Integer) "The service plan associated with the instance"
- status (String) "Instance status"
- instance_location (String) "The instance location"
- instance_type (String) "The network instance type"
- instance_bandwidth (String) "The network instance type"
- instance_cidr (String) "CIDR associated with network instance"
- network_id (Integer) "The network id associated with the instance"
- transaction_id (String) "The network id associated with the instance"
- date_created (String) "Timestamp on instance creation"
- last_updated (String) "Timestamp on last instance update"
- instance_created_by (String) "User who created the instance"
- instance_owner (String) "Instance owner"

## Terraform Input Variables
`variables.tf`
### Required
- lumen_username "Lumen username"
- lumen_password "Lumen password"
- lumen_api_access_token "Lumen Api access token"
- lumen_api_refresh_token "Lumen Api refresh token"

Each of the variables are defined in `terraform.tfvars`.

### Example usage
`terraform.tfvars` 
```hcl
lumen_username = $consumer_key
lumen_password = $consumer_secret
lumen_api_access_token = $lumen_api_access_token
lumen_api_refresh_token = $lumen_api_refresh_token
```

## References
<a id="1">[1]</a> API doc: http://apidocs.edge.lumen.com/#get-all-instances
