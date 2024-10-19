output "ip_addr" {
  value = {
    for team in local.team_list :
    team => try(aws_instance.main[team]["public_ip"], aws_spot_instance_request.main[team]["public_ip"])
  }
}

output "portal_url" {
  value = "http://${aws_s3_bucket_website_configuration.portal.website_endpoint}"
}

output "apigateway_url" {
  value = aws_apigatewayv2_stage.portal.invoke_url
}
