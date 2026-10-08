# CLAUDE.md — infra/envs/

Root module cho từng môi trường: `dev/`, `staging/`, `prod/`. Mỗi thư mục là 1 state độc lập.

## File trong mỗi env (giống nhau, chỉ khác `backend.tf` key và `terraform.tfvars`)
| File | Nội dung |
|---|---|
| `backend.tf` | S3 backend, key `base-go/<env>/terraform.tfstate`, `use_lockfile = true` |
| `versions.tf` | Ghim Terraform + provider (aws, kubernetes, helm v2) |
| `providers.tf` | AWS (`default_tags`), kubernetes/helm xác thực qua `aws eks get-token` |
| `main.tf` | Gọi module: vpc → eks → ecr (có điều kiện) → iam → platform |
| `variables.tf` / `outputs.tf` | Input/Output của env |
| `terraform.tfvars` | Giá trị env — **KHÔNG chứa secret** |

## Khác biệt giữa env (đặt trong `terraform.tfvars`)
| | dev | staging | prod |
|---|---|---|---|
| VPC CIDR | 10.10.0.0/16 | 10.20.0.0/16 | 10.30.0.0/16 |
| AZ / NAT | 2 / 1 | 3 / 1 | 3 / mỗi AZ 1 |
| Node | SPOT, 1–3 | 2–4 | m6i/m6a, 3–10 |
| `manage_shared_resources` | **true** (ECR + role CI) | false | false |

## Quy tắc
- `main.tf`, `variables.tf`, `outputs.tf`, `providers.tf`, `versions.tf` **phải giống nhau** giữa 3 env. Sửa ở 1 env thì sửa cả 3; khác biệt chỉ nằm trong `terraform.tfvars`.
- CIDR các env không được trùng (để có thể peering sau này).
- `manage_shared_resources = true` chỉ ở **1** env.
- Không đổi `key` trong `backend.tf` khi đã có state (sẽ mất liên kết state) — nếu cần, dùng `terraform init -migrate-state`.

## Lệnh
```bash
make plan-dev      # = terraform -chdir=infra/envs/dev init && plan -out=tfplan
make apply-dev
terraform -chdir=infra/envs/dev output kubeconfig_command
```
