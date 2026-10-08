variable "project" {
  type    = string
  default = "base-go"
}

variable "environment" {
  type = string
}

variable "region" {
  type    = string
  default = "ap-southeast-1"
}

variable "app_name" {
  type    = string
  default = "my-app"
}

variable "app_namespace" {
  type    = string
  default = "my-app"
}

variable "vpc_cidr" {
  type = string
}

variable "az_count" {
  type    = number
  default = 3
}

variable "single_nat_gateway" {
  type    = bool
  default = true
}

variable "cluster_version" {
  type    = string
  default = "1.33"
}

variable "endpoint_public_access_cidrs" {
  type    = list(string)
  default = ["0.0.0.0/0"]
}

variable "node_instance_types" {
  type = list(string)
}

variable "node_capacity_type" {
  type    = string
  default = "ON_DEMAND"
}

variable "node_min_size" {
  type = number
}

variable "node_max_size" {
  type = number
}

variable "node_desired_size" {
  type = number
}

variable "manage_shared_resources" {
  description = "Tạo ECR + role GitHub CI (chỉ bật ở 1 môi trường)"
  type        = bool
  default     = false
}

variable "github_oidc_provider_arn" {
  description = "Output của infra/bootstrap"
  type        = string
  default     = ""
}

variable "github_repo" {
  type    = string
  default = "my-org/base-go"
}

variable "install_argocd" {
  type    = bool
  default = true
}

variable "create_cluster_secret_store" {
  type    = bool
  default = true
}
