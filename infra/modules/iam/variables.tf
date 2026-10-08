variable "app_name" {
  type = string
}

variable "app_namespace" {
  type = string
}

variable "environment" {
  type = string
}

variable "cluster_name" {
  type = string
}

variable "create_github_ci_role" {
  description = "Chỉ bật ở 1 môi trường (thường là account chứa ECR)"
  type        = bool
  default     = false
}

variable "github_oidc_provider_arn" {
  type    = string
  default = ""
}

variable "github_repo" {
  description = "Dạng <org>/<repo>"
  type        = string
  default     = ""
}

variable "ecr_repository_arns" {
  type    = list(string)
  default = []
}

variable "tags" {
  type    = map(string)
  default = {}
}
