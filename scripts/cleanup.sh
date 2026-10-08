#!/usr/bin/env bash
# Dọn tài nguyên môi trường dev. Mặc định chỉ xóa cluster kind local.
#   ./scripts/cleanup.sh          → xóa cluster kind local
#   ./scripts/cleanup.sh --aws    → terraform destroy môi trường dev (hỏi xác nhận)
set -euo pipefail

ROOT=$(cd "$(dirname "$0")/.." && pwd)
CLUSTER=${CLUSTER:-my-app-local}

if command -v kind >/dev/null && kind get clusters | grep -qx "$CLUSTER"; then
  kind delete cluster --name "$CLUSTER"
fi

if [[ "${1:-}" == "--aws" ]]; then
  read -r -p "Destroy TOÀN BỘ hạ tầng dev trên AWS? Gõ 'dev' để xác nhận: " ans
  [[ "$ans" == "dev" ]] || { echo "Hủy."; exit 1; }
  # Xóa Ingress trước để ALB Controller dọn ALB, tránh VPC bị kẹt khi destroy
  kubectl -n my-app delete ingress --all --ignore-not-found || true
  terraform -chdir="$ROOT/infra/envs/dev" destroy
fi
