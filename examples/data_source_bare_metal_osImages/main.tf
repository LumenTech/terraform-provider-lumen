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

data "lumen_bare_metal_os_images" "os_images" {
  location_id = var.location_id
}

output "os_images" {
  value = data.lumen_bare_metal_os_images.os_images
  description = "Lumen bare metal OS images"
}