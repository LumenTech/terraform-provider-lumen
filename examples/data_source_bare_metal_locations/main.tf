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

data "lumen_bare_metal_locations" "locations" {}

output "locations" {
  value = data.lumen_bare_metal_locations.locations
  description = "Lumen bare metal locations"
}