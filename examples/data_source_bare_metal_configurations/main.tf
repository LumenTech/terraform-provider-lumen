terraform {
  required_providers {
    lumen = {
      source = "LumenTech/lumen"
      version = "2.3.0"
    }
  }
}

provider "lumen" {
  consumer_key = var.consumer_key
  consumer_secret = var.consumer_secret
  account_number = var.account_number
}

data "lumen_bare_metal_configurations" "configurations" {
  location_id = var.location_id
}

output "configurations" {
  value = data.lumen_bare_metal_configurations.configurations
  description = "Lumen bare metal configurations"
}