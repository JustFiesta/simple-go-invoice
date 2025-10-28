data "aws_caller_identity" "current" {}

resource "aws_iam_role" "cluster" {
  name                 = var.cluster_role_name
  permissions_boundary = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:policy/DefaultBoundaryPolicy"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "eks.amazonaws.com"
      }
    }]
  })

  tags = var.tags
}

resource "aws_iam_role_policy_attachment" "cluster_policies" {
  for_each = toset([
    "arn:aws:iam::aws:policy/AmazonEKSBlockStoragePolicy",
    "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy",
    "arn:aws:iam::aws:policy/AmazonEKSComputePolicy",
    "arn:aws:iam::aws:policy/AmazonEKSLoadBalancingPolicy",
    "arn:aws:iam::aws:policy/AmazonEKSNetworkingPolicy"
  ])

  role       = aws_iam_role.cluster.name
  policy_arn = each.value
}

resource "aws_iam_role" "node" {
  name                 = var.node_role_name
  permissions_boundary = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:policy/DefaultBoundaryPolicy"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ec2.amazonaws.com"
      }
    }]
  })

  tags = var.tags
}

resource "aws_iam_role_policy_attachment" "node_policies" {
  for_each = toset([
    "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly",
    "arn:aws:iam::aws:policy/AmazonEKSWorkerNodeMinimalPolicy"
  ])

  role       = aws_iam_role.node.name
  policy_arn = each.value
}

module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 21.0"

  endpoint_public_access                   = var.public_endpoint
  enable_cluster_creator_admin_permissions = true

  name               = var.cluster_name
  kubernetes_version = var.k8s_version

  vpc_id     = var.vpc_id
  subnet_ids = concat(var.public_subnets, var.private_subnets)

  iam_role_permissions_boundary      = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:policy/DefaultBoundaryPolicy"
  node_iam_role_permissions_boundary = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:policy/DefaultBoundaryPolicy"

  enable_irsa = false
#   create_iam_role = false
#   create_node_iam_role = false
#   create_auto_mode_iam_resources = false

  iam_role_arn    = aws_iam_role.cluster.arn

  compute_config = {
    enabled       = true
    node_pools    = ["general-purpose"]
    node_role_arn = aws_iam_role.node.arn
  }

  tags = var.tags
}