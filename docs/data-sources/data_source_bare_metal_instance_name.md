| Page_Title      | Description                                 |
|-----------------|---------------------------------------------|
| Data_Source_Bare_Metal_Instance_Name  | Details on bare metal instance(s) based on instance(s) name        |

## Introduction
This document provides details on data source of Lumen bare metal instance(s) details based on instance name(s). The API details are provided in Ref [[1]](#1).

## Example Usage
`main.tf`
```hcl
data "lumen_bare_metal_instance_name" "instance" {
    name = var.instance_name
}

output "instance" {
    value = data.lumen_bare_metal_instance_name.instance
}
```

## Schema

### Required
- name (String) "Instance Name"

### Computed
- id (Integer) "Instance ID"
- description (String) "The instance description"
- cloud_id (Integer) "The ID of the cloud associated with the instance"
- group_id (Integer) "The ID of the group associated with the instance"
- instance_type_id (Integer) "The type of instance to provision"
- instance_layout_id (Integer) "The layout to provision the instance from"
- plan_id (Integer) "The service plan associated with the instance"
- status (String) "Instance status"
- instance_location (String) "The instance location"
- instance_ip (String) "The instance ip address"
- network_id (Integer) "The network id associated with the instance"
- network_type (String) "The network type associated with the instance"
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

- instance_name "Instance Name"

Each of the variables are defined in `terraform.tfvars`.

## Example usage
`terraform.tfvars`
```hcl
lumen_username = $consumer_key
lumen_password = $consumer_secret
lumen_api_access_token = $lumen_api_access_token
lumen_api_refresh_token = $lumen_api_refresh_token

# Instance name
instance_name = $instance_name
```

## References
<a id="1">[1]</a> API: GET https://api.lumen.com/EdgeServices/v1/Compute/api/instances/:name
