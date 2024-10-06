resource "aws_security_group" "player" {
  name   = "player_instance"
  vpc_id = aws_vpc.vpc.id

  ingress {
    description = "HTTP"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = [
      "0.0.0.0/0"
    ]
  }

  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = [
      "0.0.0.0/0"
    ]
  }

  egress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"
    cidr_blocks = [
      "0.0.0.0/0"
    ]
  }
}

resource "aws_spot_instance_request" "player" {
  for_each = local.players

  ami           = "ami-06cda439fc5c0da1b"
  instance_type = "c5.xlarge"

  subnet_id = aws_subnet.player_subnet.id

  vpc_security_group_ids = [
    aws_security_group.player.id
  ]

  user_data = <<-EOF
#!/bin/bash
mkdir /home/ishocon/.ssh
curl https://github.com/${each.value}.keys > /home/ishocon/.ssh/authorized_keys
${join("\n", [for admin in local.admins : "curl https://github.com/${admin}.keys >> /home/ishocon/.ssh/authorized_keys"])}
chown -R ishocon:ishocon /home/ishocon/.ssh
useradd -u 1001 -g 1001 -o -N -d /home/ishocon -s /bin/bash ${each.value}

EOF

  root_block_device {
    volume_size = 30
  }

  wait_for_fulfillment = true

  tags = {
    player_name = each.value
  }
}

resource "aws_ec2_tag" "spot_instance" {
  for_each    = aws_spot_instance_request.player
  resource_id = each.value.spot_instance_id
  key         = "player_name"
  value       = each.key
}
