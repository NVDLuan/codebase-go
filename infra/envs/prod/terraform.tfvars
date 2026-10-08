# Giá trị môi trường — KHÔNG chứa secret
environment                  = "prod"
vpc_cidr                     = "10.30.0.0/16"
az_count                     = 3
single_nat_gateway           = false
endpoint_public_access_cidrs = ["203.0.113.0/24"] # TODO: thay bằng IP VPN/văn phòng/CI
node_instance_types          = ["m6i.large", "m6a.large"]
node_min_size                = 3
node_max_size                = 10
node_desired_size            = 3
