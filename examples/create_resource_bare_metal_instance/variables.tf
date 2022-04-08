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

variable "instance_description" {
  description = "Lumen instance description"
  type = string
}

variable "group_id" {
  description = "Instance group id"
  type = number
}

variable "cloud_id" {
  description = "Instance cloud id"
  type = number
}

variable "plan_id" {
  description = "Instance plan id"
  type = number
}

variable "instance_type" {
  description = "Instance type"
  type = string
}

variable "instance_type_code" {
  description = "Instance type code"
  type = string
}

variable "instance_layout_id" {
  description = "Instance layout id"
  type = number
}

variable "instance_resource_pool_id" {
  description = "Instance resource pool id"
  type = number
}

variable "create_user" {
  description = "Create user id"
  type = bool
}

variable "instance_location" {
  type = string
  description = "Lumen edge location"
}

variable "instance_bandwidth" {
  description = "Lumen edge bandwidth"
  type = string
}

variable "instance_network_type" {
  description = "Lumen instance network type"
  type = string
}

