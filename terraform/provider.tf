terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "3.37.0"
    }
  }
}

provider "aws" {
  profile = "ishocon1"
  region  = "ap-northeast-1"
}
