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

    # Instance custom configs
    location = lumen_network_instance.tf_nw_test.instance_location
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
