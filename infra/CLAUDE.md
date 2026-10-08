# CLAUDE.md — infra/

Toàn bộ hạ tầng AWS bằng **Terraform** (>= 1.10, provider AWS ~> 5.80).

## Cấu trúc
```
bootstrap/   # Chạy 1 lần, state local: S3 bucket state + GitHub OIDC provider → bootstrap/CLAUDE.md
modules/     # Module tái sử dụng, KHÔNG chứa giá trị môi trường            → modules/CLAUDE.md
envs/        # Root module cho từng env, state riêng                       → envs/CLAUDE.md
```

## Thứ tự dựng lần đầu
`bootstrap` → `envs/dev` (vpc → eks → ecr → iam → platform) → tạo secret trong Secrets Manager → apply Argo CD Application.

## Nguyên tắc
- Chạy Terraform **trong `envs/<env>/`**, không chạy trong `modules/`.
- Mỗi env 1 state: `s3://mycompany-tfstate/base-go/<env>/terraform.tfstate`, khóa bằng S3 native lockfile.
- Không commit: `.terraform/`, `*.tfstate*`, `*.tfplan` (đã có trong `.gitignore`). `.terraform.lock.hcl` **nên** commit.
- Không có secret trong `.tf`/`.tfvars`. Secret ứng dụng tạo trong Secrets Manager ngoài Terraform (hoặc chỉ tạo container secret, không đặt value).
- Ưu tiên module chính thức `terraform-aws-modules/*`, ghim `~> major.minor`.
- Trên CI: plan khi PR (`infra-plan.yml`), apply sau merge (`infra-apply.yml`); prod cần approval. Không apply prod từ máy cá nhân trừ khi sự cố.

## Kiểm tra sau khi sửa
```bash
terraform fmt -recursive infra
cd infra/envs/dev && terraform init -backend=false && terraform validate
```
