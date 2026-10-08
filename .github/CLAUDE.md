# CLAUDE.md — .github/

Pipeline CI/CD (GitHub Actions). Thay đổi ở đây do DevOps duyệt (xem `CODEOWNERS`).

## Workflows
| File | Trigger | Việc làm |
|---|---|---|
| `app-ci.yml` | PR / push `main` có đổi `app/**` | gitleaks → gofmt → go vet → golangci-lint (depguard) → sqlc diff → unit/integration/e2e test → docker build → Trivy → (push main) push ECR, tag = SHA 7 ký tự |
| `app-cd.yml` | `workflow_run` của `app-ci` thành công trên `main` | `kustomize edit set image` trong `deploy/overlays/dev`, commit `[skip ci]` |
| `infra-plan.yml` | PR có đổi `infra/**` | fmt check, validate, plan cho 3 env, comment plan lên PR |
| `infra-apply.yml` | push `main` có đổi `infra/**` | apply tuần tự dev → staging → prod; job gắn GitHub Environment theo env |

## Quy ước
- Xác thực AWS **chỉ qua OIDC** (`aws-actions/configure-aws-credentials` + `role-to-assume`). Không dùng access key trong secrets.
- Role ARN, registry lấy từ `vars.*` (repo/environment variables), không hard-code account.
  - `AWS_CI_ROLE_ARN` (output `github_ci_role_arn` của `infra/envs/dev`)
  - `ECR_REGISTRY`, `AWS_TF_PLAN_ROLE_ARN`, `AWS_TF_APPLY_ROLE_ARN` (2 role TF chưa có trong Terraform, cần tạo riêng)
- Giữ `paths:` filter để pipeline app và infra không kích hoạt lẫn nhau.
- `permissions:` khai báo tối thiểu cho từng workflow.
- Ghim version action theo tag major (`@v4`); Trivy ghim tag cụ thể.
- Approval prod: cấu hình *required reviewers* cho GitHub Environment `prod`, không làm trong YAML.
- Không thêm bước `kubectl apply` vào workflow (GitOps, xem root `CLAUDE.md`).
- `app-cd.yml` chỉ bump **dev**. Promote staging/prod là PR thủ công đổi `newTag`.

## Khi sửa
- Đổi tên app/ECR repo → sửa `env.ECR_REPOSITORY` (app-ci) và tên image trong `app-cd`.
- Đổi version Terraform → đồng bộ `terraform_version` ở cả 2 workflow infra với `required_version` trong `infra/envs/*/versions.tf`.
