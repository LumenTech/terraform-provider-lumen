# Terraform core utility
# Gets instance details based on instance name

data "lumen_bare_metal_instance_name" "instance" {
    name = var.instance_name
}

output "instance" {
    value = data.lumen_bare_metal_instance_name.instance
}
