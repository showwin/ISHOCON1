output "ip_addr" {
  value = module.main.ip_addr
}

output "scoreboard_url" {
  value = module.main.scoreboard_url
}

# For commands in Makefile
output "apigateway_url" {
  value = module.main.apigateway_url
}

output "api_gateway_id" {
  value = module.main.api_gateway_id
}

output "api_gateway_get_route_id" {
  value = module.main.api_gateway_get_route_id
}

output "api_gateway_put_route_id" {
  value = module.main.api_gateway_put_route_id
}
