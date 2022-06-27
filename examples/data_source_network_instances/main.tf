# Terraform core utility
# Display network instance(s) details
data "lumen_network_instances" "nw" {}

output "nw_instances" {
    value = data.lumen_network_instances.nw.instances
    description = "Lumen all network instances"
}
