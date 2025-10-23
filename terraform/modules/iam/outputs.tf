output "cluster_iam_role_arn" {
  value = aws_iam_role.cluster.arn
}

output "node_iam_role_arn" {
  value = aws_iam_role.node.arn
}
