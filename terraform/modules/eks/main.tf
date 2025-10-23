module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 21.0"

  endpoint_public_access                   = var.public_endpoint
  enable_cluster_creator_admin_permissions = true

  name               = var.cluster_name
  kubernetes_version = var.k8s_version

  vpc_id     = var.vpc_id
  subnet_ids = var.private_subnets

  iam_role_permissions_boundary      = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:policy/DefaultBoundaryPolicy"
  node_iam_role_permissions_boundary = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:policy/DefaultBoundaryPolicy"

#   enable_irsa = false

  eks_managed_node_groups = {
    group = {
      instance_types = var.instance_types

      min_size     = var.min_size
      max_size     = var.max_size
      desired_size = var.desired_size

      iam_role_permissions_boundary = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:policy/DefaultBoundaryPolicy"

      tags = merge(
        var.tags,
        {
          Name = "${var.cluster_name}-eks-node-group"
        }
      )
    }
  }

  tags = var.tags
}

data "aws_caller_identity" "current" {}