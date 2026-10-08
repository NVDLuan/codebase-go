# Giá trị môi trường — KHÔNG chứa secret
environment              = "dev"
vpc_cidr                 = "10.10.0.0/16"
az_count                 = 2
single_nat_gateway       = true
node_instance_types      = ["t3.large"]
node_capacity_type       = "SPOT"
node_min_size            = 1
node_max_size            = 3
node_desired_size        = 2
manage_shared_resources  = true
github_oidc_provider_arn = "arn:aws:iam::<account_id>:oidc-provider/token.actions.githubusercontent.com"
