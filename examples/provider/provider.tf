# Provider configuration
terraform {
  required_providers {
    lumen = {
      source = "lumen.com/edge/lumen-technologies"
      version = "0.3.0"
    }
  }
}

# Provider access creds
provider "lumen" {
  url = var.lumen_api_url
  username = null
  password = null
  access_token = var.lumen_access_token
}

