terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.0"
    }
  }
  required_version = ">= 1.2"

  backend "s3" {
    bucket = ""
    key    = ""
    region = "eu-west-1"
  }
}

provider "aws" {
  region = "eu-west-1"
}
