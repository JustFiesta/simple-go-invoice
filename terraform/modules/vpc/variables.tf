variable "vpc_name" {
  type    = string
  default = "mbocak-test-vpc"
}

variable "cluster_name" {
  type        = string
  description = "Name of the EKS cluster for subnet tagging"
}

variable "tags" {
  type    = map(string)
  default = {}
}