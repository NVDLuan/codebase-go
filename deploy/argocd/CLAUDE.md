# CLAUDE.md — deploy/argocd/

Khai báo Argo CD `Application` cho từng môi trường. Apply **1 lần** khi dựng cluster (`kubectl apply -f deploy/argocd/app-<env>.yaml`), sau đó Argo CD tự theo dõi Git.

## Quy ước
- `source.path: deploy/overlays/<env>`, `targetRevision: main`.
- `destination.name: eks-<env>` — mỗi env 1 cluster. Nếu Argo CD chạy trong chính cluster đó, đổi thành `server: https://kubernetes.default.svc`.
- dev/staging: `automated` + `prune` + `selfHeal`.
- **prod: không có `automated`** → sync thủ công sau khi PR promote được merge. Đừng bật auto-sync cho prod khi chưa có quyết định (ghi ADR).
- `CreateNamespace=false` vì namespace được quản lý trong `deploy/base/namespace.yaml`.
- Giữ finalizer `resources-finalizer.argocd.argoproj.io` để xoá Application thì dọn luôn resource.

## Khi sửa
- Đổi `repoURL` (`my-org/base-go`) → sửa cả 3 file.
- Repo private → cần khai báo repository credential trong Argo CD (không commit token vào đây).
