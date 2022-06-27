# Terraform core utility
# Retrieve network instance details based on instance id
data "lumen_network_instance_id" "instance" {
    id = var.instance_id
}

output "instance" {
    value = data.lumen_network_instance_id.instance
    description = "Lumen network instance details based on instance id"
}
