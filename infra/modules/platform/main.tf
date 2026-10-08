# Các add-on dùng chung cho cluster, cài bằng Helm.

# ---------- AWS Load Balancer Controller ----------
module "alb_controller_pod_identity" {
  source  = "terraform-aws-modules/eks-pod-identity/aws"
  version = "~> 1.7"

  name                            = "${var.cluster_name}-alb-controller"
  attach_aws_lb_controller_policy = true
  associations = {
    main = {
      cluster_name    = var.cluster_name
      namespace       = "kube-system"
      service_account = "aws-load-balancer-controller"
    }
  }
  tags = var.tags
}

resource "helm_release" "alb_controller" {
  name       = "aws-load-balancer-controller"
  namespace  = "kube-system"
  repository = "https://aws.github.io/eks-charts"
  chart      = "aws-load-balancer-controller"
  version    = var.chart_versions.alb_controller

  values = [yamlencode({
    clusterName    = var.cluster_name
    region         = var.region
    vpcId          = var.vpc_id
    serviceAccount = { name = "aws-load-balancer-controller" }
  })]

  depends_on = [module.alb_controller_pod_identity]
}

# ---------- External Secrets Operator ----------
module "external_secrets_pod_identity" {
  source  = "terraform-aws-modules/eks-pod-identity/aws"
  version = "~> 1.7"

  name                                  = "${var.cluster_name}-external-secrets"
  attach_external_secrets_policy        = true
  external_secrets_secrets_manager_arns = var.secrets_manager_arns
  associations = {
    main = {
      cluster_name    = var.cluster_name
      namespace       = "external-secrets"
      service_account = "external-secrets"
    }
  }
  tags = var.tags
}

resource "helm_release" "external_secrets" {
  name             = "external-secrets"
  namespace        = "external-secrets"
  create_namespace = true
  repository       = "https://charts.external-secrets.io"
  chart            = "external-secrets"
  version          = var.chart_versions.external_secrets

  values = [yamlencode({
    installCRDs    = true
    serviceAccount = { name = "external-secrets" }
  })]

  depends_on = [module.external_secrets_pod_identity]
}

resource "kubernetes_manifest" "cluster_secret_store" {
  count = var.create_cluster_secret_store ? 1 : 0

  manifest = {
    apiVersion = "external-secrets.io/v1beta1"
    kind       = "ClusterSecretStore"
    metadata   = { name = "aws-secrets-manager" }
    spec = {
      provider = {
        aws = {
          service = "SecretsManager"
          region  = var.region
        }
      }
    }
  }

  depends_on = [helm_release.external_secrets]
}

# ---------- metrics-server (cho HPA) ----------
resource "helm_release" "metrics_server" {
  name       = "metrics-server"
  namespace  = "kube-system"
  repository = "https://kubernetes-sigs.github.io/metrics-server"
  chart      = "metrics-server"
  version    = var.chart_versions.metrics_server
}

# ---------- Argo CD ----------
resource "helm_release" "argocd" {
  count = var.install_argocd ? 1 : 0

  name             = "argocd"
  namespace        = "argocd"
  create_namespace = true
  repository       = "https://argoproj.github.io/argo-helm"
  chart            = "argo-cd"
  version          = var.chart_versions.argocd

  values = [yamlencode({
    configs = { params = { "server.insecure" = true } }
  })]
}
