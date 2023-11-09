# Provider configuration
terraform {
  required_providers {
    lumen = {
      source = "LumenTech/lumen"
      version = "1.0.0"
    }
  }
}

provider "lumen" {
  # Configuration options
  username = var.username
  password = var.password
  account_number = var.account_number
  api_access_token = var.api_access_token
  api_refresh_token = var.api_refresh_token
}
