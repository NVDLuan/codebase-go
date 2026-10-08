terraform {
  backend "s3" {
    bucket       = "mycompany-tfstate"
    key          = "base-go/staging/terraform.tfstate"
    region       = "ap-southeast-1"
    encrypt      = true
    use_lockfile = true # S3 native locking (Terraform >= 1.10)
  }
}
