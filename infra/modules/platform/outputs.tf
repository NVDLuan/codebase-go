output "alb_controller_role_arn" {
  value = module.alb_controller_pod_identity.iam_role_arn
}

output "external_secrets_role_arn" {
  value = module.external_secrets_pod_identity.iam_role_arn
}
