variable "username" {
  description = "Lumen username (used for authentication) required"
  type = string
}

variable "password" {
  description = "Lumen password (used for authentication) required"
  type = string
}

variable "account_number" {
  description = "Account Number that will be acted upon (optional)"
  type = string
}

variable "api_access_token" {
  description = "Lumen api access token (option) used with (0.X.X) release bare metal commands"
  type = string
}

variable "api_refresh_token" {
  description = "Lumen api refresh token (optional) used with (0.X.X) release bare metal command"
}
