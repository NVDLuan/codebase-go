# CLAUDE.md — base-go (root)

Dự án 1 service (`my-app`) chạy trên AWS EKS. Terraform (hạ tầng) + Kustomize (manifest) + GitHub Actions (CI/CD) + Argo CD (GitOps).

## Bản đồ thư mục
| Thư mục | Ghi chú | CLAUDE.md riêng |
|---|---|---|
| `app/` | Go 1.26, Clean Architecture (domain/app/transport/adapter), port 8080 | `app/CLAUDE.md` |
| `deploy/` | Kustomize base + overlays + Argo CD | `deploy/CLAUDE.md` (+ base, overlays, argocd) |
| `infra/` | Terraform: bootstrap, modules, envs | `infra/CLAUDE.md` (+ modules, envs, bootstrap) |
| `.github/` | Workflows CI/CD | `.github/CLAUDE.md` |
| `scripts/` | Script tiện ích | `scripts/CLAUDE.md` |
| `docs/` | Kiến trúc, runbook, ADR | `docs/CLAUDE.md` |

Đọc `CLAUDE.md` của thư mục con trước khi sửa file trong đó.

## Quy tắc bắt buộc (áp dụng toàn repo)
1. **Không bao giờ commit secret.** Không tạo `Secret` K8s có giá trị thật, không đưa secret vào `terraform.tfvars`, `.env`, workflow. Secret nằm ở AWS Secrets Manager, vào cluster qua External Secrets.
2. **Không dùng tag `latest`.** Image tag = 7 ký tự đầu git SHA.
3. **Không `kubectl apply` lên EKS từ CI.** Chỉ đổi Git, Argo CD sync.
4. **Tách pipeline:** đổi `app/` không được kích hoạt Terraform và ngược lại.
5. **Mỗi môi trường độc lập:** state Terraform riêng (`infra/envs/<env>`), overlay riêng (`deploy/overlays/<env>`).
6. Đổi `prod` (overlay hoặc infra) luôn cần approval, nên nêu rõ trong PR.

## Placeholder cần thay khi dùng thật
`<account_id>`, `<env-cert-id>`, `my-org`, `example.com`, `mycompany-tfstate`, `203.0.113.0/24`.
Region mặc định: `ap-southeast-1`. Tên app/namespace/ServiceAccount: `my-app` (đổi thì đổi đồng bộ ở app, deploy, infra, workflows).

## Lệnh hay dùng
```bash
make test                 # go vet + go test (lint/integration: xem app/CLAUDE.md)
make render ENV=dev       # kubectl kustomize deploy/overlays/dev
make plan-dev             # terraform plan
make lint                 # pre-commit: fmt, tflint, yamllint, gitleaks
```

## Kiểm tra sau khi sửa
- Sửa `deploy/` → `make render ENV=dev|staging|prod` phải chạy được.
- Sửa `infra/` → `terraform fmt -recursive infra` và `terraform validate` trong env liên quan.
- Sửa `app/` → `make test`.
