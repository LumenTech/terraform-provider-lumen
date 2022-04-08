# Variables and user credentials to access lumen resources.
variable "lumen_api_url" {
  description = "Lumen API url"
  type = string
}

variable "lumen_access_token" {
  description = "Lumen user token (user should have API access)"
  type = string
}

variable "lumen_username" {
  description = "Lumen username (user should have API access)"
  type = string
}

variable "lumen_password" {
  description = "Lumen password"
  type = string
}

variable "instance_name" {
  description = "Lumen instance name"
  type = string
}
