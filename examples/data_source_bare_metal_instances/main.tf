# Terraform core utility
# Display bare-metal instance(s) details 
data "lumen_bare_metal_instances" "all" {}

output "all_instances" {
  value = data.lumen_bare_metal_instances.all.instances
  description = "Lumen bare metal instance(s) details"
}
