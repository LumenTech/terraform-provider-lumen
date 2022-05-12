| Page_Title      | Description                                 |
|-----------------|---------------------------------------------|
| Resource_Order_Update_Bare_Metal_Instance_Id  | Details on bare metal instance update |

## Introduction
This document provides details on resource order to update Lumen bare metal instance(s). The API details are provided in Ref [[1]](#1).

## Example Usage
`main.tf`
```hcl
resource "lumen_bare_metal_instance" "tf_test" {
    # updated name
    name = var.instance_name_updated
    # updated description
    description = var.instance_description_updated
    
    tags = {
        name = "tf-test"
    }

    # updated labels
    labels = ["tf_resource_test_updated", "demo_updated"]

    evar {
        name = "test-tf-ml4tst"
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
- updated name (String) "Updated instance name"
- updated description (String) "Updated instance description"
- updated tags (JSON formatted, [name value] pair) "updated instance tags"
- updated labels (List of strings) "updated instance labels"
- evar (JSON formatted, [name value] pair) "Instance evars"

### Computed
Updated instance details in `output`.

## Terraform Input Variables
`variables.tf`
### Required
- lumen_username "Lumen username"
- lumen_password "Lumen password"
- lumen_api_access_token "Lumen Api access token"
- lumen_api_refresh_token "Lumen Api refresh token"

- "instance_name_updated" : "Instance name" (String)
- "instance_description_updated" : "Lumen instance description" (String)

Each of the variables are defined in `terraform.tfvars`.

## Example usage
`terraform.tfvars`
```hcl
# Url and credentials
lumen_username = $consumer_key
lumen_password = $consumer_secret
lumen_api_access_token = $lumen_api_access_token
lumen_api_refresh_token = $lumen_api_refresh_token

# Instance name
#instance_name = "tf-demo"

# Instance description
#instance_description = "Terraform demo instance"

# Instance name - updated
instance_name_updated = "updated-tf-demo"

# Instance description - updated
instance_description_updated = "Updated - Terraform demo instance"
```

## References
<a id="1">[1]</a> API doc: http://apidocs.edge.lumen.com/#updating-an-instance
