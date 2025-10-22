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

variable "min_size" {
    type = number
    default = 1
}

variable "max_size" {
    type = number
    default = 1
}

variable "desired_size" {
    type = number
    default = 1
}

variable ami_type {
    type = string
    default = "AL2023_x86_64_STANDARD"
}

variable instance_types {
    type = list(string)
    default = ["t3.small"]
}