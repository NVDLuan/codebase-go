# 3. Dùng EKS Pod Identity thay cho IRSA

- Trạng thái: Chấp nhận
- Ngày: 2026-10-07

## Bối cảnh
Pod cần quyền AWS (Secrets Manager...). Có 2 lựa chọn: IRSA và EKS Pod Identity.

## Quyết định
Dùng EKS Pod Identity: association tạo bằng Terraform (`infra/modules/iam`), ServiceAccount không cần annotation.

## Hệ quả
- Trust policy đơn giản (`pods.eks.amazonaws.com`), không phụ thuộc OIDC URL của cluster.
- Cần add-on `eks-pod-identity-agent` (đã bật trong `infra/modules/eks`).
