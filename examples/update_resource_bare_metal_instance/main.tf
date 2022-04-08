# Terraform core utility
# The below configuration will update below instance details
# name
# description
# labels

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

