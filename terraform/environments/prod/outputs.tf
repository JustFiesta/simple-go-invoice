output "vpc_id" {
  value = module.vpc.vpc_id
}

output "eks_cluster_name" {
  value = module.eks.cluster_name
}

output "frontend_repo_url" {
  value = module.ecr.frontend_repo_url
}

output "backend_repo_url" {
  value = module.ecr.backend_repo_url
}
