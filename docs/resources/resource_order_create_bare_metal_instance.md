| Page_Title      | Description                                 |
|-----------------|---------------------------------------------|
| Resource_Order_Create_Bare_Metal_Instance_Id  | Details on bare metal instance creation |

## Introduction
This document provides details on resource order to create Lumen bare metal instance(s). The API details are provided in Ref [[1]](#1).

## Example Usage
`main.tf`
```hcl
resource "lumen_bare_metal_instance" "tf_test" {
    name = var.instance_name
    description = var.instance_description
    group_id = var.group_id
    cloud_id = var.cloud_id
    plan_id = var.plan_id
    
    instance_type = var.instance_type
    instance_type_code = var.instance_type_code
    instance_layout_id = var.instance_layout_id
    instance_resource_pool_id = var.instance_resource_pool_id

    create_user = var.create_user

    location =  var.instance_location
    bandwidth = var.instance_bandwidth
    network_type = var.instance_network_type 

    tags = {
        name = "tf-test"
    }

    labels = ["tf_resource_test", "demo"]
    evar {
        name = "test-tf"
        value = "tf-demo"
        export = true
        masked = true
    }
}

output "tf_test_instance" {
    value = lumen_bare_metal_instance.tf_test
}
```

## Schema

### Required
- name (String) "Instance name"
- description (String) "Instance description"
- group_id (Integer) "Instance group id"
- cloud_id (Integer) "Instance cloud id"
- plan_id (Integer) "Instance plan id"
- instance_type (String) "Instance type"
- instance_type_code (String) "Instance type code"
- instance_layout_id (Integer) "Instance layout id"
- instance_resource_pool_id (Integer) "Instance resource pool id"
- create_user (bool) "Create user"
- location (String) "Edge location"
- bandwidth (String) "Instance bandwidth"
- network_type (String) "Instance network type"
- tags (JSON formatted, [name value] pair) "Instance tags"
- labels (List of strings) "Instance labels"
- evar (JSON formatted, [name value] pair) "Instance evars" 

### Computed
- id (Integer) Instance ID

## Terraform Input Variables
`variables.tf`
### Required
- "lumen_api_url" : "Lumen API url"
- "lumen_access_token" : "Lumen user token (user should have API access)"
- "lumen_username" : "Lumen username (user should have API access, optional)"
- "lumen_password" : "Lumen password" (optional)
- "instance_name" : "Instance name" (String)
- "instance_description" : "Lumen instance description" (String)
- "group_id" : "Instance group id" (Integer)
- "cloud_id" : "Instance cloud id" (Integer)
- "plan_id" : "Instance plan id" (Integer)
- "instance_type" : "Instance type" (String)
- "instance_type_code" : "Instance type code" (String)
- "instance_layout_id" : "Instance layout id" (Integer)
- "instance_resource_pool_id" : "Instance resource pool id" (Integer)
- "create_user" : "Create user id" (boolean)
- "instance_location" : "Lumen edge location" (String)
- "instance_bandwidth" : "Lumen edge bandwidth" (String)
- "instance_network_type" : "Lumen instance network type" (String)

Each of the variables are defined in `terraform.tfvars`. Details related to `group_id`, `cloud_id`, `plan_id`, `instance_location`, `instance_bandwidth`, `instance_network_type` are provided in Ref [[2]](#2).

## Example usage
`terraform.tfvars`
```hcl
# Url and credentials
lumen_api_url = "https://api.lumen.com/EdgeServices/v1/Compute/api/instances"
lumen_username = null
lumen_password = null
lumen_access_token = "0000-0000-0000-0000"

# Instance name
instance_name = "tf-test"

# Instance description
instance_description = "Terraform test instance"

# Instance group id
group_id = $group_id

# Instance cloud id
cloud_id = $cloud_id

# Instance plan id
plan_id = $plan_id

# Instance type
instance_type = "$os"

# Instance type code
instance_type_code = "$os"

# Instance layout id
instance_layout_id = $instance_layout_id

# Instance resource pool id
instance_resource_pool_id = $instance_resource_pool_id

# Instance create user
create_user = $y/n

# Edge location for resource creation
instance_location = $edge_location

# Edge bandwidth for resource creation
instance_bandwidth = $edge_bandwidth

# Instance network type for resource creation
instance_network_type = $instance_network_type
```

## References
<a id="1">[1]</a> API doc: http://apidocs.edge.lumen.com/#create-an-instance

<a id="2">[2]</a> API doc: http://apidocs.edge.lumen.com/#id-code-tables
