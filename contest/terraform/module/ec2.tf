resource "aws_spot_instance_request" "player" {
  for_each = var.use_spot_instance ? { for idx, player in var.players : idx => player } : {}

  ami           = var.ami_id
  instance_type = var.instance_type

  subnet_id = module.vpc.public_subnets[0]

  vpc_security_group_ids = [
    aws_security_group.player.id
  ]

  user_data = <<-EOF
#!/bin/bash
mkdir /home/ishocon/.ssh
curl https://github.com/${each.value}.keys > /home/ishocon/.ssh/authorized_keys
${join("\n", [for admin in var.admins : "curl https://github.com/${admin}.keys >> /home/ishocon/.ssh/authorized_keys"])}
chown -R ishocon:ishocon /home/ishocon/.ssh
useradd -u 1001 -g 1001 -o -N -d /home/ishocon -s /bin/bash ${each.value}

EOF

  root_block_device {
    volume_size = 8
  }

  wait_for_fulfillment = true

  tags = {
    Name = "ISHOCON1 Player - ${each.value}"
    player_name = each.value
  }
}
