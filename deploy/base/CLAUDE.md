# CLAUDE.md — deploy/base/

Manifest dùng chung cho mọi môi trường. Namespace `my-app` được gán ở `kustomization.yaml`.

## File
| File | Ghi chú |
|---|---|
| `kustomization.yaml` | Danh sách resource + label chung. Thêm file mới phải khai báo ở đây. |
| `namespace.yaml` | Pod Security Admission `restricted` → pod phải đạt chuẩn restricted. |
| `serviceaccount.yaml` | Tên `my-app`; quyền AWS cấp qua EKS Pod Identity (không cần annotation). |
| `deployment.yaml` | Probes, resources, securityContext, topology spread theo AZ, `/tmp` emptyDir. |
| `service.yaml` | ClusterIP, port 80 → `http` (8080). |
| `ingress.yaml` | ALB, `target-type: ip`; host & `certificate-arn` do overlay patch. |
| `hpa.yaml` | CPU 70%; overlay chỉnh min/max. |
| `pdb.yaml` | `minAvailable: 1` (prod patch lên 2). |
| `networkpolicy.yaml` | Ingress 8080 từ VPC; egress DNS + 443. |
| `configmap.yaml` | `my-app-config`, chỉ giá trị không nhạy cảm. |
| `externalsecret.yaml` | Tạo Secret `my-app-secret` từ Secrets Manager key `my-app/<env>`. |

## Bất biến (không được bỏ)
- `runAsNonRoot`, `allowPrivilegeEscalation: false`, `readOnlyRootFilesystem: true`, `capabilities.drop: [ALL]`, `seccompProfile: RuntimeDefault` — bỏ đi thì pod bị namespace `restricted` từ chối.
- `resources.requests` + `limits.memory` (không đặt limit CPU để tránh throttling).
- readiness + liveness probe.
- `image: my-app` là **tên logic**; tag/registry thật do `images:` trong overlay quyết định. Đừng ghi registry/tag vào base.

## Khi thêm resource
- Thêm vào `resources:` trong `kustomization.yaml`.
- Nếu pod cần gọi dịch vụ mới (DB port 5432, Redis 6379...) → mở thêm rule egress trong `networkpolicy.yaml`.
- Cần lưu trữ → dùng PVC với StorageClass `gp3` (mặc định), không tạo PV tĩnh. Lưu ý Deployment nhiều replica + EBS (RWO) không dùng chung được, cân nhắc StatefulSet.
- Cần quyền AWS mới → sửa policy trong `infra/modules/iam`, không dùng node role.
