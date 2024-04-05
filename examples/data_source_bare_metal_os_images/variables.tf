variable "consumer_key" {
  description = "Lumen consumer key (used for authentication)"
  type = string
}

variable "consumer_secret" {
  description = "Lumen consumer secret (used for authentication)"
  type = string
}

variable "account_number" {
  description = "Account Number that will be acted upon"
  type = string
}

variable "location_id" {
  description = "The id of a location"
  type = string
}