output "app_role_arn" {
  value = aws_iam_role.app.arn
}

output "github_ci_role_arn" {
  value = try(aws_iam_role.github_ci[0].arn, null)
}
