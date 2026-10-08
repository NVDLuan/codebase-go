locals {
  name = "${var.project}-${var.environment}"
  tags = {
    Project     = var.project
    Environment = var.environment
    ManagedBy   = "terraform"
  }
}

data "aws_caller_identity" "current" {}

module "vpc" {
  source = "../../modules/vpc"

  name               = local.name
  cidr               = var.vpc_cidr
  az_count           = var.az_count
  single_nat_gateway = var.single_nat_gateway
}

module "eks" {
  source = "../../modules/eks"

  cluster_name                 = local.name
  cluster_version              = var.cluster_version
  vpc_id                       = module.vpc.vpc_id
  private_subnet_ids           = module.vpc.private_subnet_ids
  endpoint_public_access_cidrs = var.endpoint_public_access_cidrs
  node_instance_types          = var.node_instance_types
  node_capacity_type           = var.node_capacity_type
  node_min_size                = var.node_min_size
  node_max_size                = var.node_max_size
  node_desired_size            = var.node_desired_size
}

# ECR dùng chung cho mọi môi trường: chỉ tạo ở 1 env (manage_shared_resources = true)
module "ecr" {
  source = "../../modules/ecr"
  count  = var.manage_shared_resources ? 1 : 0

  repositories = [var.app_name]
}

module "iam" {
  source = "../../modules/iam"

  app_name      = var.app_name
  app_namespace = var.app_namespace
  environment   = var.environment
  cluster_name  = module.eks.cluster_name

  create_github_ci_role    = var.manage_shared_resources
  github_oidc_provider_arn = var.github_oidc_provider_arn
  github_repo              = var.github_repo
  ecr_repository_arns      = var.manage_shared_resources ? values(module.ecr[0].repository_arns) : []
}

module "platform" {
  source = "../../modules/platform"

  cluster_name                = module.eks.cluster_name
  region                      = var.region
  vpc_id                      = module.vpc.vpc_id
  install_argocd              = var.install_argocd
  create_cluster_secret_store = var.create_cluster_secret_store
  secrets_manager_arns = [
    "arn:aws:secretsmanager:${var.region}:${data.aws_caller_identity.current.account_id}:secret:${var.app_name}/${var.environment}*",
  ]

  depends_on = [module.eks]
}
