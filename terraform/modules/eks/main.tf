module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 21.0"

  endpoint_public_access                   = var.public_endpoint
  enable_cluster_creator_admin_permissions = true

  name               = var.cluster_name
  kubernetes_version = var.k8s_version

  vpc_id     = var.vpc_id
  subnet_ids = var.private_subnets

  iam_role_permissions_boundary = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:policy/DefaultBoundaryPolicy"

  create_iam_role = false
  iam_role_arn    = var.cluster_iam_role_arn

  compute_config = {
    enabled       = true
    node_pools    = ["general-purpose"]
    node_role_arn = var.node_iam_role_arn
  }

  tags = var.tags
}

data "aws_caller_identity" "current" {}