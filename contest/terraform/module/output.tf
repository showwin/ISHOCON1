output "ip_addr" {
  value = {
    for team in local.team_list :
    team => try(aws_instance.main[team]["public_ip"] , aws_spot_instance_request.main[team]["public_ip"])
  }
}
