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

  tags = var.tags


  access_entries = {
    cluster-admin = {
      principal_arn = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:user/your-username"
      policy_associations = {
        admin-policy = {
          policy_arn = "arn:aws:eks::aws:cluster-access-policy/AmazonEKSClusterAdminPolicy"
          access_scope = {
            type = "cluster"
          }
        }
      }
    },

    read-only-user = {
      principal_arn = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:user/readonly-user"
      policy_associations = {
        view-policy = {
          policy_arn = "arn:aws:eks::aws:cluster-access-policy/AmazonEKSViewPolicy"
          access_scope = {
            type = "cluster"
          }
        }
      }
    },

    namespace-admin = {
      principal_arn = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/NamespaceAdmin"
      policy_associations = {
        edit-policy = {
          policy_arn = "arn:aws:eks::aws:cluster-access-policy/AmazonEKSEditPolicy"
          access_scope = {
            type       = "namespace"
            namespaces = ["default", "production"]
          }
        }
      }
    }
  }
}


data "aws_caller_identity" "current" {}
