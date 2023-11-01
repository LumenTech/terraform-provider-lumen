variable "username" {
  description = "Lumen username (used for authentication)"
  type = string
}

variable "password" {
  description = "Lumen password (used for authentication)"
  type = string
}

variable "accountNumber" {
  description = "Account Number that will be acted upon"
  type = string
}

variable "location_id" {
  description = "The id of a location"
  type = string
}