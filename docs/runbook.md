# Runbook

## Kết nối cluster
```bash
aws eks update-kubeconfig --name base-go-<env> --region ap-southeast-1
```

## Rollback ứng dụng
1. `git revert <commit đổi tag>` trên `deploy/overlays/<env>/kustomization.yaml`, merge.
2. Argo CD tự sync (dev/staging). Prod: bấm **Sync** trên Argo CD UI hoặc `argocd app sync my-app-prod`.
3. Khẩn cấp (bỏ qua Git, Argo CD sẽ cảnh báo OutOfSync): `kubectl -n my-app rollout undo deploy/my-app`.

## Pod CrashLoopBackOff
```bash
kubectl -n my-app get pods
kubectl -n my-app describe pod <pod>
kubectl -n my-app logs <pod> --previous
```
- `CreateContainerConfigError` → thường do thiếu `my-app-secret`: xem mục ExternalSecret.
- Lỗi ghi file → `readOnlyRootFilesystem: true`, chỉ `/tmp` được ghi.

## ExternalSecret không đồng bộ
```bash
kubectl -n my-app get externalsecret my-app
kubectl -n my-app describe externalsecret my-app
```
- `SecretSyncedError` → kiểm tra secret `my-app/<env>` có tồn tại trong Secrets Manager.
- `AccessDenied` → kiểm tra policy trong `infra/modules/platform` (`secrets_manager_arns`).

## Ingress không tạo ALB
```bash
kubectl -n kube-system logs deploy/aws-load-balancer-controller
```
- Thiếu tag subnet `kubernetes.io/role/elb` → xem `infra/modules/vpc`.
- Sai `certificate-arn` → xem patch trong overlay.

## HPA hiện `<unknown>`
- metrics-server chưa chạy: `kubectl -n kube-system get deploy metrics-server`.

## Terraform state bị khóa
- Kiểm tra không có job `infra-apply` đang chạy, rồi: `terraform force-unlock <LOCK_ID>`.
