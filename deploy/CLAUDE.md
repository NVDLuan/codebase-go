# CLAUDE.md — deploy/

Toàn bộ manifest Kubernetes, quản lý bằng **Kustomize**, triển khai bằng **Argo CD** (GitOps).

## Cấu trúc
```
base/        # Phần chung mọi môi trường        → base/CLAUDE.md
overlays/    # dev | staging | prod: chỉ phần khác → overlays/CLAUDE.md
argocd/      # Argo CD Application cho từng env  → argocd/CLAUDE.md
```

## Nguyên tắc
- Git là nguồn sự thật: mọi thay đổi trên cluster đi qua commit ở đây. Không `kubectl edit/apply` tay trên dev/staging/prod (Argo CD selfHeal sẽ ghi đè).
- Thứ gì giống nhau ở mọi env → `base/`. Thứ gì khác → `overlays/<env>/`. Không copy cả resource sang overlay.
- Không có `Secret` chứa giá trị thật. Chỉ dùng `ExternalSecret`.
- Label chuẩn `app.kubernetes.io/name: my-app` do `base/kustomization.yaml` gắn; selector của Deployment/Service/PDB/NetworkPolicy dựa vào label này.

## Kiểm tra trước khi commit
```bash
kubectl kustomize deploy/overlays/dev
kubectl kustomize deploy/overlays/staging
kubectl kustomize deploy/overlays/prod
# (tuỳ chọn) validate schema
kubectl kustomize deploy/overlays/prod | kubeconform -strict -ignore-missing-schemas -summary
```
Cả 3 overlay phải render được. Sửa `base/` thì render cả 3.

## Phụ thuộc cluster (cài bởi `infra/modules/platform`)
- AWS Load Balancer Controller → `IngressClass alb`
- External Secrets Operator + `ClusterSecretStore aws-secrets-manager`
- metrics-server → HPA
- VPC CNI bật network policy → `NetworkPolicy`
- EKS Pod Identity association cho ServiceAccount `my-app` (`infra/modules/iam`)
