terraform {
  required_version = "1.9.7"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.69"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "~> 2.6.0"
    }
    local = {
      source  = "hashicorp/local"
      version = "~> 2.5.0"
    }
  }
}

provider "aws" {
  region = "ap-northeast-1"
}
