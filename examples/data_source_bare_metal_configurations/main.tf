terraform {
  required_providers {
    lumen = {
      source = "LumenTech/lumen"
      version = "1.0.0"
    }
  }
}

provider "lumen" {
  username = var.username
  password = var.password
  account_number = var.accountNumber
}

data "lumen_bare_metal_configurations" "configurations" {
  location_id = var.location_id
}

output "configurations" {
  value = data.lumen_bare_metal_configurations.configurations
  description = "Lumen bare metal configurations"
}