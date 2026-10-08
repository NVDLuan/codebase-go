# Giá trị môi trường — KHÔNG chứa secret
environment         = "staging"
vpc_cidr            = "10.20.0.0/16"
az_count            = 3
single_nat_gateway  = true
node_instance_types = ["t3.large"]
node_min_size       = 2
node_max_size       = 4
node_desired_size   = 2
