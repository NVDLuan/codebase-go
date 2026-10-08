data "aws_partition" "current" {}

# ---------- Role cho pod của app (EKS Pod Identity) ----------
data "aws_iam_policy_document" "pod_identity_trust" {
  statement {
    actions = ["sts:AssumeRole", "sts:TagSession"]
    principals {
      type        = "Service"
      identifiers = ["pods.eks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "app" {
  name               = "${var.app_name}-${var.environment}"
  assume_role_policy = data.aws_iam_policy_document.pod_identity_trust.json
  tags               = var.tags
}

data "aws_iam_policy_document" "app" {
  # Quyền tối thiểu: chỉ đọc secret của chính app ở môi trường này
  statement {
    actions   = ["secretsmanager:GetSecretValue", "secretsmanager:DescribeSecret"]
    resources = ["arn:${data.aws_partition.current.partition}:secretsmanager:*:*:secret:${var.app_name}/${var.environment}*"]
  }
}

resource "aws_iam_role_policy" "app" {
  name   = "app"
  role   = aws_iam_role.app.id
  policy = data.aws_iam_policy_document.app.json
}

resource "aws_eks_pod_identity_association" "app" {
  cluster_name    = var.cluster_name
  namespace       = var.app_namespace
  service_account = var.app_name
  role_arn        = aws_iam_role.app.arn
}

# ---------- Role cho GitHub Actions (OIDC) ----------
data "aws_iam_policy_document" "github_trust" {
  count = var.create_github_ci_role ? 1 : 0

  statement {
    actions = ["sts:AssumeRoleWithWebIdentity"]
    principals {
      type        = "Federated"
      identifiers = [var.github_oidc_provider_arn]
    }
    condition {
      test     = "StringEquals"
      variable = "token.actions.githubusercontent.com:aud"
      values   = ["sts.amazonaws.com"]
    }
    condition {
      test     = "StringLike"
      variable = "token.actions.githubusercontent.com:sub"
      values   = ["repo:${var.github_repo}:ref:refs/heads/main"]
    }
  }
}

resource "aws_iam_role" "github_ci" {
  count              = var.create_github_ci_role ? 1 : 0
  name               = "${var.app_name}-github-ci"
  assume_role_policy = data.aws_iam_policy_document.github_trust[0].json
  tags               = var.tags
}

data "aws_iam_policy_document" "github_ci" {
  count = var.create_github_ci_role ? 1 : 0

  statement {
    actions   = ["ecr:GetAuthorizationToken"]
    resources = ["*"]
  }
  statement {
    actions = [
      "ecr:BatchCheckLayerAvailability", "ecr:InitiateLayerUpload", "ecr:UploadLayerPart",
      "ecr:CompleteLayerUpload", "ecr:PutImage", "ecr:BatchGetImage",
    ]
    resources = var.ecr_repository_arns
  }
}

resource "aws_iam_role_policy" "github_ci" {
  count  = var.create_github_ci_role ? 1 : 0
  name   = "ecr-push"
  role   = aws_iam_role.github_ci[0].id
  policy = data.aws_iam_policy_document.github_ci[0].json
}
