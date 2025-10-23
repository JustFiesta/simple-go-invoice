output "vpc_id" {
  value = module.vpc.vpc_id
}

output "cluster_endpoint" {
  value = module.eks.cluster_arn
}

output "repo_url" {
  value = module.ecr.repo_url
}
