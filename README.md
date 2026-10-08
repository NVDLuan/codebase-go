# base-go

Ứng dụng `my-app` chạy trên AWS EKS: hạ tầng bằng Terraform, manifest bằng Kustomize, CI/CD bằng GitHub Actions, deploy theo GitOps (Argo CD).

## Cấu trúc

| Thư mục | Nội dung |
|---|---|
| `app/` | Source code Go 1.26 (Clean Architecture) + Dockerfile |
| `deploy/` | Manifest Kubernetes (Kustomize) + Argo CD Application |
| `infra/` | Terraform: bootstrap, modules, envs |
| `.github/` | Workflow CI/CD, CODEOWNERS |
| `scripts/` | Chạy local, dọn dẹp |
| `docs/` | Kiến trúc, runbook, ADR |

Mỗi thư mục chính có `CLAUDE.md` mô tả quy ước riêng.

## Bắt đầu nhanh

```bash
make help              # Xem các lệnh
make test              # Unit test
make local             # Chạy trên kind local
make render ENV=prod   # Xem manifest đã render
make plan-dev          # terraform plan môi trường dev
```

## Dựng hạ tầng lần đầu

1. `cd infra/bootstrap && terraform init && terraform apply` → tạo S3 state + GitHub OIDC provider.
2. Điền `github_oidc_provider_arn` và `<account_id>` vào `infra/envs/dev/terraform.tfvars`.
3. `make plan-dev && make apply-dev` (lần đầu đặt `create_cluster_secret_store = false`, apply lại với `true` sau khi CRD đã cài).
4. Tạo secret `my-app/dev` trong AWS Secrets Manager.
5. `kubectl apply -f deploy/argocd/app-dev.yaml` → Argo CD sync `deploy/overlays/dev`.
6. Cấu hình GitHub: biến `AWS_CI_ROLE_ARN`, `ECR_REGISTRY`, `AWS_TF_PLAN_ROLE_ARN`, `AWS_TF_APPLY_ROLE_ARN`; Environment `prod` có required reviewers.

## Luồng thay đổi

- **Code:** PR → `app-ci` (lint, test, build, Trivy) → merge → push ECR → `app-cd` bump tag dev → Argo CD sync.
- **Promote:** PR đổi `newTag` ở `deploy/overlays/staging` rồi `prod` (cần approval).
- **Hạ tầng:** PR → `infra-plan` comment plan → merge → `infra-apply` (prod cần approval).

Chi tiết: [docs/architecture.md](docs/architecture.md), [docs/runbook.md](docs/runbook.md).
