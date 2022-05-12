# Variables and user credentials to access lumen resources.
variable "lumen_username" {
  description = "Lumen username"
  type = string
}

variable "lumen_password" {
  description = "Lumen password"
  type = string
}

variable "lumen_api_access_token" {
  description = "Lumen API access token"
  type = string
}

variable "lumen_api_refresh_token" {
  description = "Lumen API refresh token"
  type = string
}

# Only used for updating instance
variable "instance_name_updated" {
  description = "Lumen updated instance name"
  type = string
}

# Only used for updating description
variable "instance_description_updated" {
  description = "Lumen updated instance description"
  type = string
}
