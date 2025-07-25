terraform {
  required_version = ">= 1.12.0"

  backend "s3" {
    bucket       = "s3-state-bucket73579"
    key          = "golang-api/terraform.tfstate"
    region       = "us-east-1"
    encrypt      = true
    # use_lockfile = true
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.38.0"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

# Read cluster info from infra state
data "terraform_remote_state" "infra" {
  backend = "s3"
  config = {
    bucket = "s3-state-bucket73579"
    key    = "infrastructure/terraform.tfstate"
    region = "us-east-1"
  }
}

data "aws_eks_cluster_auth" "eks" {
  name = data.terraform_remote_state.infra.outputs.cluster_name
}

# Configure Kubernetes provider using remote state
provider "kubernetes" {
  host                   = data.terraform_remote_state.infra.outputs.cluster_endpoint
  cluster_ca_certificate = base64decode(data.terraform_remote_state.infra.outputs.cluster_ca_certificate)
  token                  = data.aws_eks_cluster_auth.eks.token
}

module "go_api_k8s" {
  source = "git::https://github.com/Rippyblogger/Base-AWS-Infrastructure.git//modules/kubernetes?ref=main"

  cluster_name           = data.terraform_remote_state.infra.outputs.cluster_name
  cluster_endpoint       = data.terraform_remote_state.infra.outputs.cluster_endpoint
  cluster_ca_certificate = data.terraform_remote_state.infra.outputs.cluster_certificate_authority_data

  deployment_name      = "go-api-deployment"
  namespace            = "default"
  env                  = "production"
  replicas_count       = 1
  app_name             = "go-api"
  image_name           = var.image_name
  service_account_name = "my-go-api-irsa-sa"

  oidc_arn          = data.terraform_remote_state.infra.outputs.oidc_provider_arn
  oidc_provider_url = data.terraform_remote_state.infra.outputs.oidc_provider_url
}
