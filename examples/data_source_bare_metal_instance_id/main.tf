# Terraform core utility
# Retrieve bare-metal instance details based on instance id
data "lumen_bare_metal_instance_id" "instance" {
    id = var.instance_id
}

output "instance" {
    value = data.lumen_bare_metal_instance_id.instance
    description = "Lumen bare metal instance details"
}
