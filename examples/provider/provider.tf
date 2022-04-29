# Provider configuration
terraform {
  required_providers {
    lumen = {
      source = "LumenTech/lumen"
      version = "0.3.3"
    }
  }
}

provider "lumen" {
  # Configuration options
  api_url = var.lumen_api_url
  auth_url = var.lumen_auth_url
  username = var.lumen_username
  password = var.lumen_password
  api_access_token = var.lumen_api_access_token
  api_refresh_token = var.lumen_api_refresh_token
}
