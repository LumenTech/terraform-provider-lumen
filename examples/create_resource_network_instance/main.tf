# Terraform core utility
# The below configuration will create network instance and display result.
resource "lumen_network_instance" "tf_nw_test" {
    name = var.nw_instance_name
    description = var.nw_instance_description

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
    value = lumen_network_instance.tf_nw_test
    description = "Lumen network instance details"
}
