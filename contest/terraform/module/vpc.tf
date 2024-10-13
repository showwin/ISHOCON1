data "aws_availability_zones" "available" {
  state = "available"
}

locals {
  # AZによって差が出ないように1つだけを使う
  first_az = slice(data.aws_availability_zones.available.names, 0, 1)
}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "5.13.0"

  name = var.name
  cidr = var.vpc_cidr_block

  azs             = local.first_az
  public_subnets  = [for k, v in local.first_az : cidrsubnet(var.vpc_cidr_block, 8, k + 48)]
}
