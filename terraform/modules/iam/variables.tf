variable "cluster_role_name" {
  type = string
  default = "EKSClusterRole-AutoMode"
}

variable "node_role_name" {
  type = string
  default = "EKSNodeRole-AutoMode"
}

variable "tags" {
  type    = map(string)
  default = {}
}