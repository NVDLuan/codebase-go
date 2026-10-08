# CLAUDE.md — deploy/overlays/

Mỗi thư mục = 1 môi trường, **chỉ chứa phần khác** so với `../base`.

| Env | Replicas / HPA | Đặc thù | Sync |
|---|---|---|---|
| `dev` | 1 / 1–3 | `LOG_LEVEL=debug`, host `my-app.dev.example.com` | Tự động |
| `staging` | 2 / 2–3 | host `my-app.staging.example.com` | Tự động |
| `prod` | 3 / 3–20 | resources lớn hơn, PDB `minAvailable: 2`, host lấy từ base | Thủ công |

## Thành phần trong `kustomization.yaml` của mỗi overlay
- `resources: [../../base]`
- `patches:` — file patch (strategic merge) hoặc JSON6902 inline cho Ingress host/cert, ExternalSecret key.
- ConfigMap theo env: `configmap-env.yaml` (dev/staging, patch) hoặc `configMapGenerator` (prod, `behavior: merge`).
- `images:` — `newName` (ECR) + `newTag` (git SHA).

## Image tag — quan trọng
- `deploy/overlays/dev`: `newTag` do workflow `app-cd.yml` **tự cập nhật**. Không sửa tay trừ khi rollback.
- `staging`, `prod`: promote bằng PR copy `newTag` từ env trước (dev → staging → prod). Prod cần approval (CODEOWNERS).
- Rollback = `git revert` commit đổi tag.
- Không bao giờ dùng `latest`.

## Khi thêm môi trường mới
1. Copy `staging/` → `<env>/`, sửa host, cert ARN, secret key `my-app/<env>`, `APP_ENV`.
2. Tạo `deploy/argocd/app-<env>.yaml`.
3. Tạo `infra/envs/<env>/` tương ứng.

## Kiểm tra
`kubectl kustomize deploy/overlays/<env>` phải render được; soát lại Ingress host, cert, ExternalSecret key, image.
