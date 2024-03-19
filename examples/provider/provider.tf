# Provider configuration
terraform {
  required_providers {
    lumen = {
      source = "LumenTech/lumen"
      version = "2.1.1"
    }
  }
}

provider "lumen" {
  # Configuration options
  username = var.username
  password = var.password
  account_number = var.account_number
}
