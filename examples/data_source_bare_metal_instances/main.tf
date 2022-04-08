# Terraform core utility
# Displays for instance(s) details 
data "lumen_bare_metal_instances" "all" {}

output "all_instances" {
  value = data.lumen_bare_metal_instances.all.instances
}

