| Page_Title      | Description                                 |
|-----------------|---------------------------------------------|
| Data_Source_Network_Instance_Id  | Provides details on network instance based on instance id|

## Introduction
This document provides Lumen network instance(s) detail based on instance id for a specific tenant. The API details are provided in Ref [[1]](#1).

## Example Usage
`main.tf`
```hcl
data "lumen_network_instance_id" "instance" {
    id = var.instance_id
}

output "all_instances" {
  value = data.lumen_network_instance_id.instance
}
```

## Schema

### Required
- id (Integer) "Instance ID"

### Computed
- name (String) "The names of instance(s) created"
- description (String) "The instance description"
- cloud_id (Integer) "The ID of the cloud associated with the instance"
- group_id (Integer) "The ID of the group associated with the instance"
- instance_type_id (Integer) "The type of instance to provision"
- instance_type_name (Integer) "The type of instance to provision"
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

- instance_id "Instance ID"

Each of the variables are defined in `terraform.tfvars`.

### Example usage
`terraform.tfvars` 
```hcl
lumen_username = $consumer_key
lumen_password = $consumer_secret
lumen_api_access_token = $lumen_api_access_token
lumen_api_refresh_token = $lumen_api_refresh_token

# Instance id
instance_id = $instance_id
```

## References
<a id="1">[1]</a> API doc: http://apidocs.edge.lumen.com/#get-all-instances
