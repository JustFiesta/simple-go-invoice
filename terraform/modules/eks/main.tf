module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 21.0"

  endpoint_public_access = var.public_endpoint

  name               = var.cluster_name
  kubernetes_version = var.k8s_version
  
  vpc_id             = var.vpc_id
  subnet_ids         = var.private_subnets

  compute_config = {
    enabled    = true
    node_pools = ["general-purpose"]
  }

  tags = var.tags
}
