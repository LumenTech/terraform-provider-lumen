# Provider configuration
terraform {
  required_providers {
    lumen = {
      source = "lumen.com/edge/lumen-technologies"
      version = "0.3.1"
    }
  }
}

# Provider access creds
provider "lumen" {
  api_url = var.lumen_api_url
  auth_url = var.lumen_auth_url
  username = var.lumen_username
  password = var.lumen_password
  api_access_token = var.lumen_api_access_token
  api_refresh_token = var.lumen_api_refresh_token
}
