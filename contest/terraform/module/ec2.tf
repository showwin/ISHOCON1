resource "aws_instance" "main" {
  for_each = var.use_spot_instance ? {} : { for idx, team in local.team_list : team => team }

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
echo "export BENCH_TEAM_NAME=${each.value}" >> /home/ishocon/.bashrc
echo "export BENCH_SCOREBOARD_APIGW_URL=${aws_apigatewayv2_stage.scoreboard.invoke_url}" >> /home/ishocon/.bashrc

EOF

  root_block_device {
    volume_size = 8
  }
}

resource "aws_ec2_tag" "instance_name" {
  for_each = aws_instance.main

  resource_id = each.value.id
  key         = "Name"
  value       = "${var.name} - ${each.key}"
}

resource "aws_ec2_tag" "instance_team_name" {
  for_each = aws_instance.main

  resource_id = each.value.id
  key         = "team_name"
  value       = each.key
}

resource "aws_ec2_tag" "instance_players" {
  for_each = aws_instance.main

  resource_id = each.value.id
  key         = "players"
  value       = join(",", var.teams[each.key])
}
