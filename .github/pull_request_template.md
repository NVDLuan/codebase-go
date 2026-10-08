## Mô tả
<!-- Thay đổi gì, vì sao -->

## Loại thay đổi
- [ ] App (`app/`)
- [ ] Manifest K8s (`deploy/`)
- [ ] Hạ tầng (`infra/`)
- [ ] CI/CD (`.github/`)

## Checklist
- [ ] Không commit secret (đã chạy `pre-commit` / gitleaks)
- [ ] Đã test local (`make test`, `make render ENV=...`)
- [ ] Nếu đổi `infra/`: đã xem `terraform plan` trên PR
- [ ] Nếu đổi `deploy/overlays/prod`: có kế hoạch rollback
