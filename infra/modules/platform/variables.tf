variable "cluster_name" {
  type = string
}

variable "region" {
  type = string
}

variable "vpc_id" {
  type = string
}

variable "secrets_manager_arns" {
  description = "Secret mà External Secrets được phép đọc"
  type        = list(string)
  default     = ["arn:aws:secretsmanager:*:*:secret:*"]
}

variable "install_argocd" {
  type    = bool
  default = true
}

variable "create_cluster_secret_store" {
  description = "Để false ở lần apply đầu (CRD chưa có), bật true ở lần apply sau"
  type        = bool
  default     = true
}

variable "chart_versions" {
  type = object({
    alb_controller   = string
    external_secrets = string
    metrics_server   = string
    argocd           = string
  })
  default = {
    alb_controller   = "1.11.0"
    external_secrets = "0.12.1"
    metrics_server   = "3.12.2"
    argocd           = "7.7.11"
  }
}

variable "tags" {
  type    = map(string)
  default = {}
}
