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

resource "lumen_bare_metal_network" "network" {
  #  Example request data
  name = "testNetwork5"
  location_id = "DNVTCO56LEC"
  network_size_id = "6529723924b8bf31ebd998e2"
  network_type = "DUAL_STACK_INTERNET"
}

resource "lumen_bare_metal_network" "network2" {
  # Create a private network with an existing VRF
  name = "testNetwork6"
  location_id = "DNVTCO56LEC"
  network_type = "PRIVATE"
  vrf = "88/VP12/003434/ASRT"
}

resource "lumen_bare_metal_network" "network3" {
  # Create a private network with a new VRF
  name = "testNetwork7"
  location_id = "DNVTCO56LEC"
  network_type = "PRIVATE"
  vrf_description = "testPrivateNetwork"
}

output "network" {
  value = lumen_bare_metal_network.network
  description = "Lumen bare metal network details"
}

output "network2" {
  value = lumen_bare_metal_network.network2
  description = "Private network with existing VRF details"
}

output "network3" {
  value = lumen_bare_metal_network.network3
  description = "Private network with new VRF details"
}