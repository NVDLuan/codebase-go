# 2. Deploy theo GitOps với Argo CD + Kustomize

- Trạng thái: Chấp nhận
- Ngày: 2026-10-07

## Bối cảnh
Cần deploy lên nhiều môi trường EKS, có audit, rollback dễ, không cấp quyền cluster cho CI.

## Quyết định
- Manifest viết bằng Kustomize (`deploy/base` + `deploy/overlays/<env>`).
- CI chỉ build image và commit image tag; Argo CD pull từ Git và sync.
- Prod sync thủ công sau khi PR promote được duyệt.

## Hệ quả
- Git là nguồn sự thật duy nhất, rollback bằng `git revert`.
- CI không cần kubeconfig hay quyền ghi vào cluster.
- Cần vận hành thêm Argo CD.
