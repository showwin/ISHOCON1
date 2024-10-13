output "ip_addr" {
  value = {
    for player in local.players :
    player => aws_spot_instance_request.player[player]["public_ip"]
  }
}
