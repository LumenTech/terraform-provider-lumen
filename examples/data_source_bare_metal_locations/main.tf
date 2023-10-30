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

data "lumen_bare_metal_locations" "locations" {}

output "locations" {
  value = data.lumen_bare_metal_locations.locations
  description = "Lumen bare metal locations"
}