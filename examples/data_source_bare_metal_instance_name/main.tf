# Terraform core utility
# Retrieve bare-metal instance details based on instance name
data "lumen_bare_metal_instance_name" "instance" {
    name = var.instance_name
}

output "instance" {
    value = data.lumen_bare_metal_instance_name.instance
    description = "Lumen bare metal instance details based on instance name"
}
