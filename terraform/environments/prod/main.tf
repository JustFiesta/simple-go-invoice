# Reapply 14
locals {
  project_name = "mbocak-k8"
  tags = {
    env     = "prod"
    owner   = "mbocak"
    project = "kubernetes-learning"
  }
}

module "vpc" {
  source       = "../../modules/vpc"
  vpc_name     = "${local.project_name}-vpc"
  cluster_name = "${local.project_name}-cluster"
  tags         = local.tags
}

module "eks" {
  source          = "../../modules/eks"
  cluster_name    = "${local.project_name}-cluster"
  k8s_version     = "1.34"

  public_endpoint = true

  cluster_role_name = "${local.project_name}-cluster-role"
  node_role_name    = "${local.project_name}-node-role"

  vpc_id          = module.vpc.vpc_id
  private_subnets = module.vpc.private_subnets
  public_subnets = module.vpc.public_subnets

  tags = local.tags
}

module "ecr" {
  source       = "../../modules/ecr"
  project_name = local.project_name
  tags         = local.tags
}
