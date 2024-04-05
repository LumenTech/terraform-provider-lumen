# Provider configuration
terraform {
  required_providers {
    lumen = {
      source = "LumenTech/lumen"
      version = "2.2.0"
    }
  }
}

provider "lumen" {
  # Configuration options
  consumer_key = var.consumer_key
  consumer_secret = var.consumer_secret
  account_number = var.account_number
}
