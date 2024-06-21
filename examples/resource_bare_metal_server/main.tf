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
  attach_networks {
    network_id = "65c283e988f85707cc53b308"
    assign_ipv6_address = true
  }
  attach_networks {
    network_id = "6553dc0d861724132e81ce42"
  }
  attach_networks {
    network_id = "6553dc22861724132e81ce43"
  }
  username = "admin"
  password = "**********"
}

resource "lumen_bare_metal_server" "server3" {
  #  Example server provisioning request with hyperthreading true/false
  name = "BASTION02"
  location_id = "DNVTCO56LEC"
  configuration_name = "small_plus"
  os_image_name = "Ubuntu 20.04"
  network_size_id = "64ef800284b3f203bc27f9e2"
  enable_hyperthreading = true # or false
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

output "server3" {
  sensitive = true
  value = lumen_bare_metal_server.server3
  description = "Server with hyperthreading enabled/disabled"
}