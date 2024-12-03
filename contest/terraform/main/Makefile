.PHONY: delete-scoreboard-data
delete-scoreboard-data:
	terraform output -raw apigateway_url | xargs -I {} curl -X "DELETE" "{}teams"

.PHONY: close-scoreboard
close-scoreboard:
	aws apigatewayv2 update-route --api-id `terraform output -raw api_gateway_id` --route-id `terraform output -raw api_gateway_get_route_id` --authorization-type AWS_IAM > /dev/null
	@echo "Scoreboard is closed."

.PHONY: reopen-scoreboard
reopen-scoreboard:
	aws apigatewayv2 update-route --api-id `terraform output -raw api_gateway_id` --route-id `terraform output -raw api_gateway_get_route_id` --authorization-type NONE > /dev/null
	@echo "Scoreboard is reopened."

.PHONY: freeze-scoreboard
freeze-scoreboard:
	aws apigatewayv2 update-route --api-id `terraform output -raw api_gateway_id` --route-id `terraform output -raw api_gateway_put_route_id` --authorization-type AWS_IAM > /dev/null
	@echo "Scoreboard is freezed."

.PHONY: unfreeze-scoreboard
unfreeze-scoreboard:
	aws apigatewayv2 update-route --api-id `terraform output -raw api_gateway_id` --route-id `terraform output -raw api_gateway_put_route_id` --authorization-type NONE > /dev/null
	@echo "Scoreboard can accepts new score."
