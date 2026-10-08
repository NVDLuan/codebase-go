output "cluster_name" {
  value = module.eks.cluster_name
}

output "cluster_endpoint" {
  value = module.eks.cluster_endpoint
}

output "kubeconfig_command" {
  value = "aws eks update-kubeconfig --name ${module.eks.cluster_name} --region ${var.region}"
}

output "app_role_arn" {
  value = module.iam.app_role_arn
}

output "ecr_repository_urls" {
  value = try(module.ecr[0].repository_urls, {})
}

output "github_ci_role_arn" {
  value = module.iam.github_ci_role_arn
}
