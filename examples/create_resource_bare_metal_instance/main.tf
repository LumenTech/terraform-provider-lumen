# Terraform core utility
# The below configuration will create instance and display result.

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
        name = "test-tf-ml4tst"
        value = "tf-demo"
        export = true
        masked = true
    }
}

output "tf_test_instance" {
    value = lumen_bare_metal_instance.tf_test
}

