output "cluster_id" {
    value = module.eks.cluster_id
}

output "cluster_endpoint" {
    value = module.eks.cluster_endpoint
}

output "kubeconfig-certificate-authority-data" {
    value = module.eks.cluster_certificate_authority_data
}

output "cluster_name" {
    value = module.eks.cluster_id
}

output "node_group_role_arns" {
    value = module.eks.node_group_arns
}
