| Page_Title      | Description                                 |
|-----------------|---------------------------------------------|
| Resource_Order_Create_Bare_Metal_Instance  | Details on bare metal instance creation |

## Introduction
This document provides details on resource order to create Lumen bare metal instance(s). The API details are provided in Ref [[1]](#1). In order to create a bare metal resource, a network resource needs to be created first, and then the id from the created network resource will be used in creating the bare metal resource. Example payload for creating bare-metal instance is provide in Ref [[2]](#2).

## Example Usage
`main.tf`
```hcl
# Terraform core utility
# The below configuration will create
# network instance and display result.
resource "lumen_network_instance" "tf_nw_test" {
    name = var.nw_instance_name
    description = var.nw_instance_name

    group_id = var.group_id
    cloud_id = var.cloud_id
    plan_id = var.nw_plan_id

    instance_type_code = var.nw_instance_type_code
    instance_layout_id = var.nw_instance_layout_id

    location = var.instance_location
    bandwidth = var.nw_instance_bandwidth
    network_type = var.instance_network_type

    tags = {
        name = "nw-tf-test"
    }

    labels = ["nw_tf_resource_test", "demo"]
}

output "tf_nw_test_instance" {
    value = lumen_network_instance.tf_nw_test
    description = "Lumen network instance details"
}

output "nw_id" {
    value = lumen_network_instance.tf_nw_test.network_id
    description = "Lumen network id"
    sensitive = true
}

# Below configuration will create bare metal
# instance with network instance id.
resource "lumen_bare_metal_instance" "tf_bm_test" {
    name = var.instance_name
    description = var.instance_description
    group_id = var.group_id
    cloud_id = var.cloud_id
    plan_id = var.plan_id
    # Instance Type
    instance_type_id = var.instance_type_id
    instance_type_code = var.instance_type_code
    instance_type_name = var.instance_type_name
    instance_layout_id = var.instance_layout_id

    #create_user = var.create_user
    # Instance custom configs
    location =  var.instance_location
    network_id = lumen_network_instance.tf_nw_test.network_id
    network_type = var.instance_network_type
    # Instance tags
    tags = {
        name = "tf-test"
    }
    # Instance labels and evars
    labels = ["tf_resource_test", "demo"]

    evar {
        name = "test-tf-ml4tst"
        value = "tf-demo"
        export = true
        masked = true
    }
}

output "tf_bm_test_instance" {
    value = lumen_bare_metal_instance.tf_bm_test
    description = "Lumen bare metal instance details"
}
```

## Schema

### Required Variables for Network resource creation
- name (String) "Instance name"
- description (String) "Instance description"
- group_id (Integer) "Instance group id"
- cloud_id (Integer) "Instance cloud id"
- plan_id (Integer) "Instance plan id"
- instance_type_code (String) "Instance type code"
- instance_layout_id (Integer) "Instance layout id"
- location (String) "Edge location"
- bandwidth (String) "Instance bandwidth"
- network_type (String) "Instance network type"
- tags (JSON formatted, [name value] pair) "Instance tags"
- labels (List of strings) "Instance labels"

### Required Variables for BM resource creation
- name (String) "Instance name"
- description (String) "Instance description"
- group_id (Integer) "Instance group id"
- cloud_id (Integer) "Instance cloud id"
- plan_id (Integer) "Instance plan id"
- instance_type_id (Integer) "Instance type"
- instance_type_code (String) "Instance type code"
- instance_type_name (String) "Instance type name"
- instance_layout_id (Integer) "Instance layout id"
- location (String) "Edge location"
- network_id (String) "Instance network id"
- network_type (String) "Instance network type"
- tags (JSON formatted, [name value] pair) "Instance tags"
- labels (List of strings) "Instance labels"
- evar (JSON formatted, [name value] pair) "Instance evars" 

### Computed
- id (Integer) Instance ID

## Terraform Input Variables
`variables.tf`
### Required
- lumen_username "Lumen username"
- lumen_password "Lumen password"
- lumen_api_access_token "Lumen Api access token"
- lumen_api_refresh_token "Lumen Api refresh token"
- "instance_name" : "Instance name" (String)
- "instance_description" : "Lumen instance description" (String)
- "group_id" : "Instance group id" (Integer)
- "cloud_id" : "Instance cloud id" (Integer)
- "plan_id" : "Instance plan id" (Integer)
- "instance_type_id" : "Instance type id" (Integer)
- "instance_type_code" : "Instance type code" (String)
- "instance_type_name" : "Instance type name" (String)
- "instance_layout_id" : "Instance layout id" (Integer)
- "instance_resource_pool_id" : "Instance resource pool id" (Integer)
- "instance_location" : "Lumen edge location" (String)
- "instance_bandwidth" : "Lumen edge bandwidth" (String)
- "instance_network_id" : "Lumen instance network id" (Integer)
- "instance_network_type" : "Lumen instance network type" (String)

Each of these variables are defined in `terraform.tfvars`. Details related to `group_id`, `cloud_id`, `plan_id`, `instance_location`, `instance_bandwidth`, `instance_network_type` are provided in Ref [[3]](#3).

## Example usage
`terraform.tfvars`
```hcl
# Credentials
lumen_username = "Lumen username"
lumen_password = "Lumen password"
lumen_api_access_token = "Lumen API access token"
lumen_api_refresh_token = "Lumen API refresh token"

# User input parameters for creating network instance
# Instance name
nw_instance_name = $nw_instance_name

# Instance description
nw_instance_description = $nw_instance_description

# Instance type code
nw_instance_type_code = $nw_instance_type_code

# Instance layout id
nw_instance_layout_id = $nw_instance_layout_id

# Instance plan id
nw_plan_id = $nw_plan_id

# Edge bandwidth for resource creation
nw_instance_bandwidth = $nw_instance_bandwidth

# User input parameters for creating bare-metal instance
# Instance name
instance_name = $instance_name

# Instance description
instance_description = $instance_description

# Instance group id
# This attribute is same for both bare-metal and network instances
group_id = $group_id

# Instance cloud id
# This attribute is same for both bare-metal and network instances
cloud_id = $cloud_id

# Instance type id
instance_type_id = $instance_type

# Instance type code
instance_type_code = $instance_type_code

# Instance type name
instance_type_name = $instance_type_name

# Instance layout id
instance_layout_id = $instance_layout_id

# Instance plan id
plan_id = $plan_id

# Instance create user
# create_user = $boolean(true/false)

# Edge location for resource creation
# This attribute is same for both bare-metal and network instances 
instance_location = $instance_location

# Instance network type for resource creation
# This attribute is same for both bare-metal and network instances
instance_network_type = $instance_network_type
```

## References
<a id="1">[1]</a> Swagger API doc: https://apimarketplace.lumen.com/api/edge-bare-metal?tab=document-specification#/Instances/addInstance

<a id="2">[2]</a> API doc: https://apimarketplace.lumen.com/api/edge-bare-metal?tab=code-samples

<a id="3">[3]</a> API doc: https://apimarketplace.lumen.com/api/edge-bare-metal?tab=status-error-codes
