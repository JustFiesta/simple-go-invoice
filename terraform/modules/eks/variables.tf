variable "cluster_name" {
    type = string
}

variable "k8s_version" {
    type = string
}

variable "vpc_id" {
    type = string
}

variable "private_subnets" {
    type = list(string)
}

variable "tags" {
    type = map(string)
    default = {}
}

variable "public_endpoint" {
    type = bool
    default = false
}

variable "cluster_iam_role_arn" {
  type = string
}

variable "node_iam_role_arn" {
  type = string
}