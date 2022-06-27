# Variables and user credentials to access lumen resources.
# User credentials
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

# Variables for creating nw instance
variable "nw_instance_name" {
  description = "Lumen network instance name"
  type = string
}

variable "nw_instance_description" {
  description = "Lumen network instance description"
  type = string
}

variable "nw_plan_id" {
  description = "Lumen network instance plan id"
  type = number
}

variable "nw_instance_type_code" {
  description = "Lumen network instance type code"
  type = string
}

variable "nw_instance_layout_id" {
  description = "Lumen network instance layout id"
  type = number
}

variable "nw_instance_bandwidth" {
  description = "Lumen edge bandwidth"
  type = number
}

# Variables for creating bare-metal instance
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

variable "instance_type_id" {
  description = "Instance type id"
  type = number
}

variable "instance_type_code" {
  description = "Instance type code"
  type = string
}

variable "instance_type_name" {
  description = "Instance type name"
  type = string
}

variable "instance_layout_id" {
  description = "Instance layout id"
  type = number
}
variable "plan_id" {
  description = "Instance plan id"
  type = number
}

variable "instance_location" {
  type = string
  description = "Lumen edge location"
}

variable "instance_network_type" {
  description = "Lumen instance network type"
  type = string
}
