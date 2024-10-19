resource "aws_spot_instance_request" "main" {
  for_each = var.use_spot_instance ? { for idx, team in local.team_list : team => team } : {}

  ami           = var.ami_id
  instance_type = var.instance_type

  subnet_id                   = module.vpc.public_subnets[0]
  associate_public_ip_address = true

  vpc_security_group_ids = [
    aws_security_group.main.id
  ]

  user_data = <<-EOF
#!/bin/bash
mkdir /home/ishocon/.ssh
${join("\n", [for player_in_team in var.teams[each.value] : "curl https://github.com/${player_in_team}.keys >> /home/ishocon/.ssh/authorized_keys"])}
${join("\n", [for admin in var.admins : "curl https://github.com/${admin}.keys >> /home/ishocon/.ssh/authorized_keys"])}
chown -R ishocon:ishocon /home/ishocon/.ssh
useradd -u 1001 -g 1001 -o -N -d /home/ishocon -s /bin/bash ${each.value}

EOF

  root_block_device {
    volume_size = 8
  }

  wait_for_fulfillment = true
}

resource "aws_ec2_tag" "spot_instance_name" {
  for_each    = aws_spot_instance_request.main

  resource_id = each.value.spot_instance_id
  key         = "Name"
  value       = "ISHOCON1 - ${each.key}"
}

resource "aws_ec2_tag" "spot_instance_team_name" {
  for_each    = aws_spot_instance_request.main

  resource_id = each.value.spot_instance_id
  key         = "team_name"
  value       = each.key
}

resource "aws_ec2_tag" "spot_instance_players" {
  for_each    = aws_spot_instance_request.main

  resource_id = each.value.spot_instance_id
  key         = "players"
  value       = join(",", var.teams[each.key])
}
