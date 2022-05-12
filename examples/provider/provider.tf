# Provider configuration
terraform {
  required_providers {
    lumen = {
      source = "lumen.com/lumentech/lumen"
      version = "0.3.5"
    }
  }
}

provider "lumen" {
  # Configuration options
  username = var.lumen_username
  password = var.lumen_password
  api_access_token = var.lumen_api_access_token
  api_refresh_token = var.lumen_api_refresh_token
}
