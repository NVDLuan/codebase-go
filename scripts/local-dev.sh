#!/usr/bin/env bash
# Chạy app trên cluster kind local, dùng overlay dev nhưng bỏ phần phụ thuộc AWS.
set -euo pipefail

CLUSTER=${CLUSTER:-my-app-local}
ROOT=$(cd "$(dirname "$0")/.." && pwd)
TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

command -v kind >/dev/null || { echo "Cần cài kind: https://kind.sigs.k8s.io"; exit 1; }

kind get clusters | grep -qx "$CLUSTER" || kind create cluster --name "$CLUSTER"

docker build -t my-app:local "$ROOT/app"
kind load docker-image my-app:local --name "$CLUSTER"

# Overlay tạm: bỏ Ingress/ExternalSecret (cần ALB + ESO), tạo secret giả
cat > "$TMP/kustomization.yaml" <<YAML
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - $ROOT/deploy/base
patches:
  - patch: |-
      \$patch: delete
      apiVersion: networking.k8s.io/v1
      kind: Ingress
      metadata: { name: my-app }
  - patch: |-
      \$patch: delete
      apiVersion: external-secrets.io/v1beta1
      kind: ExternalSecret
      metadata: { name: my-app }
secretGenerator:
  - name: my-app-secret
    options: { disableNameSuffixHash: true }
    literals: [DUMMY=local]
images:
  - name: my-app
    newTag: local
YAML

kubectl kustomize --load-restrictor LoadRestrictionsNone "$TMP" | kubectl apply -f -
kubectl -n my-app rollout status deploy/my-app --timeout=120s
echo "Truy cập: kubectl -n my-app port-forward svc/my-app 8080:80  →  http://localhost:8080"
