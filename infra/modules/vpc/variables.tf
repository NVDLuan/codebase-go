variable "name" {
  type = string
}

variable "cidr" {
  type    = string
  default = "10.0.0.0/16"
}

variable "az_count" {
  type    = number
  default = 3
}

variable "single_nat_gateway" {
  description = "true cho dev để tiết kiệm chi phí, false cho prod (mỗi AZ 1 NAT)"
  type        = bool
  default     = true
}

variable "tags" {
  type    = map(string)
  default = {}
}
