variable "cluster_name" {
  type = string
}

variable "cluster_version" {
  type    = string
  default = "1.33"
}

variable "vpc_id" {
  type = string
}

variable "private_subnet_ids" {
  type = list(string)
}

variable "endpoint_public_access_cidrs" {
  description = "Giới hạn IP được gọi API server (VPN/văn phòng/CI)"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}

variable "access_entries" {
  description = "EKS access entries cho IAM role (CI, admin...)"
  type        = any
  default     = {}
}

variable "node_instance_types" {
  type    = list(string)
  default = ["t3.large"]
}

variable "node_capacity_type" {
  type    = string
  default = "ON_DEMAND"
}

variable "node_min_size" {
  type    = number
  default = 2
}

variable "node_max_size" {
  type    = number
  default = 5
}

variable "node_desired_size" {
  type    = number
  default = 2
}

variable "tags" {
  type    = map(string)
  default = {}
}
