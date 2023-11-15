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

resource "lumen_bare_metal_server" "server" {
  #  Example request data with a new created network
  name = "BASTION01"
  location_id = "DNVTCO56LEC"
  configuration_name = "small_plus"
  os_image_name = "Ubuntu 20.04"
  network_name = "bastion-network"
  network_size_id = "64ef800284b3f203bc27f9e2"
  username = "admin"
  password = "**********"
}

resource "lumen_bare_metal_server" "server2" {
  #  Example request data using multiple existing networks
  name = "BASTION02"
  location_id = "DNVTCO56LEC"
  configuration_name = "small_plus"
  os_image_name = "Ubuntu 20.04"
  network_ids = [
  "65526f96861724132e81b952",
  "6553dc0d861724132e81ce42",
  "6553dc22861724132e81ce43"
  ]
  username = "admin"
  password = "**********"
}

output "server" {
  sensitive = true
  value = lumen_bare_metal_server.server
  description = "Lumen bare metal server details"
}

output "server2" {
  sensitive = true
  value = lumen_bare_metal_server.server2
  description = "Server with multiple attached networks"
}