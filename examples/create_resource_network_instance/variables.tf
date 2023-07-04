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

variable "nw_instance_name" {
  description = "Lumen network instance name"
  type = string
}

variable "nw_instance_description" {
  description = "Lumen network instance description"
  type = string
}

variable "nw_group_id" {
  description = "Lumen network instance group id"
  type = number
}

variable "nw_cloud_id" {
  description = "Lumen network instance cloud id"
  type = number
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

variable "nw_instance_location" {
  type = string
  description = "Lumen edge location"
}

variable "nw_instance_bandwidth" {
  description = "Lumen edge bandwidth"
  type = number
}
