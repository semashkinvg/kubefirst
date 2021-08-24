# todo discuss this file being a module... flat forever?
# todo mgmt_ === module. if we pull it up a layer, pass in the data object from remote state ? 1 var 
#! todo fix
data "terraform_remote_state" "eks" {
  backend = "s3"
  config = {
    bucket = "kubefirst-mgmt-state-store"
    key    = "spike/us-west-2/kubefirst/terraform/tfstate.tf"
    aws_region =var.aws_region
  }
}

variable "aws_region" {
  type = string 
}

data "aws_eks_cluster" "cluster" {
  name = data.terraform_remote_state.eks.outputs.eks_module.eks_cluster_id
}

data "aws_eks_cluster_auth" "cluster" {
  name = data.terraform_remote_state.eks.outputs.eks_module.eks_cluster_id
}

provider "kubernetes" {
  host                   = data.aws_eks_cluster.cluster.endpoint
  cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority.0.data)
  token                  = data.aws_eks_cluster_auth.cluster.token
}

resource "vault_auth_backend" "k8s_mgmt" {
  type = "kubernetes"
  path = "kubernetes/${var.aws_account_name}-${var.aws_region}"
}

data "kubernetes_service_account" "mgmt_external_secrets" {
  metadata {
    name = "external-secrets"
    namespace = "external-secrets"
  }
}

data "kubernetes_secret" "mgmt_external_secrets_token_secret" {
  metadata {
    name = data.kubernetes_service_account.mgmt_external_secrets.default_secret_name
    namespace = "external-secrets"
  }
}

resource "vault_kubernetes_auth_backend_config" "vault_k8s_auth" {
  backend            = vault_auth_backend.k8s_mgmt.path
  kubernetes_host    = data.aws_eks_cluster.cluster.endpoint
  kubernetes_ca_cert = base64decode(data.aws_eks_cluster.cluster.certificate_authority.0.data)
  token_reviewer_jwt = data.kubernetes_secret.mgmt_external_secrets_token_secret.data.token
}

resource "vault_kubernetes_auth_backend_role" "k8s_external_secrets_mgmt" {
  backend                          = vault_auth_backend.k8s_mgmt.path
  role_name                        = "external-secrets"
  bound_service_account_names      = ["external-secrets"]
  bound_service_account_namespaces = ["*"]
  token_ttl                        = 86400
  token_policies                   = ["admin", "default"]
}
