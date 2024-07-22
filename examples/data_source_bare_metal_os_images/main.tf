terraform {
  required_providers {
    lumen = {
      source = "LumenTech/lumen"
      version = ">= 2.4.0"
    }
  }
}

provider "lumen" {
  consumer_key = var.consumer_key
  consumer_secret = var.consumer_secret
  account_number = var.account_number
}

data "lumen_bare_metal_os_images" "os_images" {
  location_id = var.location_id
}

output "os_images" {
  value = data.lumen_bare_metal_os_images.os_images
  description = "Lumen bare metal OS images"
}