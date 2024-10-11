packer {
  required_plugins {
    amazon = {
      version = ">= 1.2.7"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

variable "source_ami" {
  description = "The source AMI to use for the builder"
  type        = string
}

variable "region" {
  description = "The region to build the image in"
  type        = string
  default     = "ap-northeast-1"
}

variable "security_group_id" {
  description = "The ID of the security group Packer will associate with the builder to enable access"
  type        = string
  default     = null
}

variable "ssh_keypair_name" {
  description = "The name of the SSH keypair to associate with the builder instance"
  type        = string
  default     = null
}

variable "ssh_private_key_file" {
  description = "The path to the private key file to use for SSH connection to the builder instance"
  type        = string
  default     = null
}

variable "subnet_id" {
  description = "If using VPC, the ID of the subnet, such as subnet-12345def, where Packer will launch the EC2 instance. This field is required if you are using an non-default VPC"
  type        = string
  default     = null
}

variable "instance_type" {
  description = "The instance type Packer will use for the builder"
  type        = string
  default     = "c7i.xlarge"
}

variable "iam_instance_profile" {
  description = "The IAM instance profile to associate with the builder"
  type        = string
  default     = null
}

source "amazon-ebs" "ishocon1" {
  ami_name          = "ishocon1-${formatdate("YYYYMMDD-hhmm", timestamp())}"
  instance_type     = var.instance_type
  region            = var.region
  security_group_id = var.security_group_id
  subnet_id         = var.subnet_id
  source_ami        = var.source_ami

  associate_public_ip_address = true
  ssh_username                = "ubuntu"
  ssh_interface               = "public_ip"
  communicator                = "ssh"
  ssh_keypair_name            = var.ssh_keypair_name
  ssh_private_key_file        = var.ssh_private_key_file

  run_tags = {
    Name = "ISHOCON1 AMI builder"
  }

  tags = {
    Name     = "ISHOCON1"
    base_ami = "{{ .SourceAMI }}"
    built_by = "https://github.com/showwin/ISHOCON1/tree/master/ami"
  }

  launch_block_device_mappings {
    device_name           = "/dev/sda1"
    volume_size           = 8
    volume_type           = "gp3"
    delete_on_termination = true
  }
}

build {
  name = "ishocon1"
  sources = [
    "source.amazon-ebs.ishocon1"
  ]

  # Init
  provisioner "file" {
    source      = "./scripts/init.sh"
    destination = "/tmp/init.sh"
  }

  provisioner "shell" {
    inline = [
      "sudo bash /tmp/init.sh"
    ]
  }

  # Base
  provisioner "file" {
    source      = "./scripts/base.sh"
    destination = "/tmp/base.sh"
  }

  provisioner "shell" {
    inline = [
      "sudo -u ishocon -H sh -c 'bash /tmp/base.sh'"
    ]
  }

  # WebApp
  provisioner "file" {
    source      = "../admin/.bashrc"
    destination = "/tmp/.bashrc"
  }

  provisioner "file" {
    source      = "../admin/nginx.conf"
    destination = "/tmp/nginx.conf"
  }

  provisioner "file" {
    source      = "../admin/ishocon1.dump.tar.gz"
    destination = "/tmp/ishocon1.dump.tar.gz"
  }

  provisioner "file" {
    source      = "./webapp.tar.gz"
    destination = "/tmp/webapp.tar.gz"
  }

  provisioner "file" {
    source      = "./scripts/webapp.sh"
    destination = "/tmp/webapp.sh"
  }

  provisioner "shell" {
    inline = [
      "sudo -u ishocon -H sh -c 'bash /tmp/webapp.sh'"
    ]
  }

  # Benchmarker
  provisioner "file" {
    source      = "./benchmarker.tar.gz"
    destination = "/tmp/benchmarker.tar.gz"
  }

  provisioner "file" {
    source      = "./scripts/benchmarker.sh"
    destination = "/tmp/benchmarker.sh"
  }

  provisioner "shell" {
    inline = [
      "sudo -u ishocon -H sh -c 'bash /tmp/benchmarker.sh'"
    ]
  }
}
