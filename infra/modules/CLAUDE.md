# CLAUDE.md — infra/modules/

Module Terraform tái sử dụng. **Không** có backend, provider block, hay giá trị riêng môi trường — mọi thứ đi qua `variables.tf`.

| Module | Tạo gì | Ghi chú |
|---|---|---|
| `vpc/` | VPC, subnet public/private, NAT | Wrap `terraform-aws-modules/vpc`. Gắn tag `kubernetes.io/role/elb` / `internal-elb` cho ALB Controller. |
| `eks/` | Cluster, managed node group, add-ons, StorageClass `gp3` | Add-on: coredns, kube-proxy, vpc-cni (**bật network policy**), pod-identity-agent, ebs-csi. Cần provider `kubernetes`. |
| `ecr/` | ECR repository + lifecycle | Tag `IMMUTABLE`, scan on push. Dùng chung cho mọi env. |
| `iam/` | Role cho pod app (Pod Identity) + role GitHub CI (OIDC) | Least privilege: app chỉ đọc `secret:my-app/<env>*`. CI chỉ push ECR. |
| `platform/` | Helm: ALB Controller, External Secrets (+ `ClusterSecretStore`), metrics-server, Argo CD | Cần provider `helm` + `kubernetes`. |

## Quy ước
- Mỗi module có tối thiểu `main.tf`, `variables.tf`, `outputs.tf`; thêm `versions.tf` nếu dùng provider ngoài `aws`.
- Mọi biến có `type`; biến không hiển nhiên có `description`.
- Không hard-code region/account; lấy từ biến hoặc `data "aws_caller_identity"` / `aws_partition`.
- Quyền AWS cho workload dùng **EKS Pod Identity** (`terraform-aws-modules/eks-pod-identity` hoặc `aws_eks_pod_identity_association`), không dùng node role.
- Helm chart ghim version qua `var.chart_versions` (module `platform`).
- Thay đổi interface (thêm/đổi biến bắt buộc) → cập nhật `envs/*/main.tf` cả 3 env trong cùng PR.

## Lưu ý đã biết
- `platform`: `kubernetes_manifest` cho `ClusterSecretStore` cần CRD đã tồn tại lúc plan → lần apply đầu đặt `create_cluster_secret_store = false`.
- `eks`: nâng `cluster_version` từng minor 1 lần, kiểm tra tương thích add-on trước.
