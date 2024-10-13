terraform {
  required_version = "1.9.7"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.69.0"
    }
  }
}

provider "aws" {
  region  = "ap-northeast-1"

  default_tags {
    tags = {
      Service     = "ISHOCON1"
      ManagedBy   = "Terraform"
    }
  }
}
