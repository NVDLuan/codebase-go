# CLAUDE.md — infra/bootstrap/

Chạy **1 lần duy nhất** cho mỗi AWS account, trước mọi env. Dùng **state local** (không có backend) vì bucket state chưa tồn tại.

## Tạo gì
- S3 bucket `mycompany-tfstate`: versioning, mã hoá KMS, chặn public, `prevent_destroy`.
- IAM OIDC provider cho GitHub Actions (`token.actions.githubusercontent.com`).

## Chạy
```bash
cd infra/bootstrap
terraform init
terraform apply
terraform output github_oidc_provider_arn   # điền vào infra/envs/dev/terraform.tfvars
```

## Lưu ý
- File `terraform.tfstate` local sinh ra ở đây **không commit** (đã ignore). Cất giữ an toàn hoặc migrate vào chính bucket vừa tạo (thêm backend s3 với key `base-go/bootstrap/terraform.tfstate` rồi `terraform init -migrate-state`).
- Đổi `state_bucket_name` → sửa `bucket` trong `infra/envs/*/backend.tf`.
- Không thêm tài nguyên ứng dụng vào đây; chỉ những thứ cần có trước khi Terraform chính chạy được.
