terraform {
  required_providers {
    lumen = {
      source = "LumenTech/lumen"
      version = ">= 2.5.0"
    }
  }
}

provider "lumen" {
  consumer_key = var.consumer_key
  consumer_secret = var.consumer_secret
  account_number = var.account_number
}

data "lumen_bare_metal_network_sizes" "network_sizes" {
  location_id = var.location_id
}

output "network_sizes" {
  value = data.lumen_bare_metal_network_sizes.network_sizes
  description = "Lumen bare metal network sizes"
}