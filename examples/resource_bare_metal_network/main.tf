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
  account_number = var.account_number
}

resource "lumen_bare_metal_network" "network" {
  #  Example request data
  name = "testNetwork5"
  location_id = "DNVTCO56LEC"
  network_size_id = "6529723924b8bf31ebd998e2"
  network_type = "DUAL_STACK_INTERNET"
}

output "network" {
  value = lumen_bare_metal_network.network
  description = "Lumen bare metal network details"
}