output "ip_addr" {
  value = {
    for team in local.team_list :
    team => try(aws_instance.main[team]["public_ip"], aws_spot_instance_request.main[team]["public_ip"])
  }
}

output "scoreboard_url" {
  value = "http://${aws_s3_bucket_website_configuration.scoreboard.website_endpoint}"
}

output "apigateway_url" {
  value = aws_apigatewayv2_stage.scoreboard.invoke_url
}

output "api_gateway_id" {
  value = aws_apigatewayv2_api.scoreboard.id
}

output "api_gateway_get_route_id" {
  value = aws_apigatewayv2_route.scoreboard_teams_get.id
}

output "api_gateway_put_route_id" {
  value = aws_apigatewayv2_route.scoreboard_teams_put.id
}
