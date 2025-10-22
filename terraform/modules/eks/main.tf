module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 21.0"

  endpoint_public_access = var.public_endpoint

  name               = var.cluster_name
  kubernetes_version = var.k8s_version
  
  vpc_id             = var.vpc_id
  subnet_ids         = var.private_subnets

  iam_role_permissions_boundary = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:policy/DefaultBoundaryPolicy"
  node_iam_role_permissions_boundary = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:policy/DefaultBoundaryPolicy"

  enable_irsa = false

  compute_config = {
    enabled    = true
    node_pools = ["general-purpose"]
  }

  eks_managed_node_groups = {
    example = {
      ami_type       = var.ami_type
      instance_types = var.instance_types

      min_size     = var.min_size
      max_size     = var.max_size
      desired_size = var.desired_size
    }
  }

  tags = var.tags

  access_entries = {
    cluster-admin = {
      principal_arn = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:user/mbocak"
      policy_associations = {
        admin-policy = {
          policy_arn = "arn:aws:eks::aws:cluster-access-policy/AmazonEKSClusterAdminPolicy"
          access_scope = {
            type = "cluster"
          }
        }
      }
    }
  }
}


data "aws_caller_identity" "current" {}
