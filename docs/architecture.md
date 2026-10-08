# Kiến trúc

```
                    ┌──────────────── AWS (ap-southeast-1) ─────────────────┐
 Người dùng ──HTTPS──▶ ALB (AWS LB Controller)                              │
                    │      │                                                │
                    │  ┌───▼──────────── EKS cluster (private subnets) ───┐ │
                    │  │  ns my-app: Deployment ◀─ HPA, PDB, NetPol       │ │
                    │  │     │  └─ ServiceAccount ─▶ IAM role (Pod Identity)│
                    │  │     └─ Secret ◀─ External Secrets ◀─ Secrets Mgr │ │
                    │  │  ns argocd: Argo CD ──pull── GitHub (deploy/)    │ │
                    │  └──────────────────────────────────────────────────┘ │
                    │  ECR (image)   S3 (Terraform state)                    │
                    └────────────────────────────────────────────────────────┘

 GitHub Actions ──OIDC──▶ IAM role ──push──▶ ECR
               └─commit tag──▶ deploy/overlays/dev ──▶ Argo CD sync
```

## Thành phần

| Lớp | Công nghệ | Nơi định nghĩa |
|---|---|---|
| Mạng | VPC 2–3 AZ, NAT | `infra/modules/vpc` |
| Compute | EKS managed node group (AL2023) | `infra/modules/eks` |
| Lưu trữ | EBS CSI, StorageClass `gp3` mặc định | `infra/modules/eks` |
| Ingress | AWS Load Balancer Controller (ALB) | `infra/modules/platform`, `deploy/base/ingress.yaml` |
| Secret | AWS Secrets Manager + External Secrets | `infra/modules/platform`, `deploy/base/externalsecret.yaml` |
| Danh tính pod | EKS Pod Identity | `infra/modules/iam` |
| GitOps | Argo CD | `infra/modules/platform`, `deploy/argocd/` |
| CI/CD | GitHub Actions + OIDC | `.github/workflows/` |

## Môi trường

| Env | VPC CIDR | Node | NAT | Sync Argo CD |
|---|---|---|---|---|
| dev | 10.10.0.0/16 | SPOT, 1–3 | 1 | Tự động |
| staging | 10.20.0.0/16 | 2–4 | 1 | Tự động |
| prod | 10.30.0.0/16 | 3–10 | Mỗi AZ 1 | Thủ công |
