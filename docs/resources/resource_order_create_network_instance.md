| Page_Title      | Description                                 |
|-----------------|---------------------------------------------|
| Resource_Order_Create_Network_Instance  | Details on network instance creation |

## Introduction
This document provides deatils on resource order to create Lumen network instance(s).

## Example Usage
`main.tf`
```hcl
# Terraform core utility
# The below configuration will create network instance and display result.
resource "lumen_network_instance" "tf_nw_test" {
    name = var.nw_instance_name
    description = var.nw_instance_name

    group_id = var.nw_group_id
    cloud_id = var.nw_cloud_id
    plan_id = var.nw_plan_id

    instance_type_code = var.nw_instance_type_code
    instance_layout_id = var.nw_instance_layout_id

    location = var.nw_instance_location
    bandwidth = var.nw_instance_bandwidth
    network_type = var.nw_instance_network_type

    tags = {
        name = "nw-tf-test"
    }

    labels = ["nw_tf_resource_test", "demo"]
}

output "tf_nw_test_instance" {
    value = lumen_network_instance.tf_test
    description = "Lumen network instance details"
}
```

## Schema

### Terraform Input Variables
`variables.tf`
### Required
- lumen_username "Lumen username"
- lumen_password "Lumen password"
- lumen_api_access_token "Lumen Api access token"
- lumen_api_refresh_token "Lumen Api refresh token"
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

Each of these variables are defined in `terraform.tfvars`. Details related to `group_id`, `cloud_id`, `plan_id`, `instance_location`, `instance_bandwidth`, `instance_network_type` are provided in the API doc.

### Computed
- id (Integer) Instance ID
- network_id (Integer) Network ID

### Example usage
`terraform.tfvars`
```hcl
# Credentials and input parameters
lumen_username = "Lumen username"
lumen_password = "Lumen password"
lumen_api_access_token = "Lumen API access token"
lumen_api_refresh_token = "Lumen API refresh token"

# Instance name
nw_instance_name = $nw_instance_name

# Instance description
nw_instance_description = $nw_instance_description

# Instance cloud id
nw_cloud_id = $nw_cloud_id

# Instance group id
nw_group_id = $nw_group_id

# Instance type code
nw_instance_type_code = $nw_instance_type_code

# Instance layout id
nw_instance_layout_id = $nw_instance_layout_id

# Instance plan id
nw_plan_id = $nw_plan_id

# Edge location for resource creation
nw_instance_location = $nw_instance_location

# Instance network type for resource creation
nw_instance_network_type = $nw_instance_network_type

# Edge bandwidth for resource creation
nw_instance_bandwidth = $nw_instance_bandwidth
```

## References
<a id="1">[1]</a> API doc: http://apidocs.edge.lumen.com/#create-an-instance

<a id="2">[2]</a> API doc: http://apidocs.edge.lumen.com/#id-code-tables
