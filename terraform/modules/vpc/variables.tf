variable "vpc_name" {
  type = string
  default = "mbocak-test-vpc"
}

variable "tags" {
  type = map(string)
  default = {}
}