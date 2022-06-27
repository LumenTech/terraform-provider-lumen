# Terraform core utility
# Retrieve network instance details based on instance name
data "lumen_network_instance_name" "instance" {
    name = var.instance_name
}

output "instance" {
    value = data.lumen_network_instance_name.instance
    description = "Lumen network instance details based on instance name"
}
