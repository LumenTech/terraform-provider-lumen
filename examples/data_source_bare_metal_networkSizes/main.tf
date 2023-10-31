terraform {
  required_providers {
    lumen = {
      source = "LumenTech/lumen"
      version = "0.5.3"
    }
  }
}

provider "lumen" {
  username = var.username
  password = var.password
  account_number = var.accountNumber
}

data "lumen_bare_metal_network_sizes" "network_sizes" {
  location_id = var.location_id
}

output "network_sizes" {
  value = data.lumen_bare_metal_network_sizes.network_sizes
  description = "Lumen bare metal network sizes"
}