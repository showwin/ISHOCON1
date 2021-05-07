terraform {
  backend "s3" {
    region  = "ap-northeast-1"
    bucket  = "" // FIXME
    key     = "default.tfstate"
    encrypt = true
  }
}
