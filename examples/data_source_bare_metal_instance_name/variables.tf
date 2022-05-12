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

variable "instance_name" {
  description = "Lumen instance name"
  type = string
}
