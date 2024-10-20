import json
from decimal import Decimal

import boto3

client = boto3.client("dynamodb")
dynamodb = boto3.resource("dynamodb")
table = dynamodb.Table("${dynamodb_table_name}")


def lambda_handler(event, context):
    print(event)
    body = {}
    statusCode = 200
    headers = {"Content-Type": "application/json"}

    try:
        if event["routeKey"] == "DELETE /teams":
            scan = table.scan()
            with table.batch_writer() as batch:
                for each in scan["Items"]:
                    batch.delete_item(
                        Key={"team": each["team"], "timestamp": each["timestamp"]}
                    )
            body = "Deleted all items"
        elif event["routeKey"] == "GET /teams":
            body = table.scan()
            body = body["Items"]
            responseBody = []
            for items in body:
                responseItems = {
                    "score": float(items["score"]),
                    "team": items["team"],
                    "timestamp": items["timestamp"],
                }
                responseBody.append(responseItems)
            body = responseBody
        elif event["routeKey"] == "PUT /teams":
            requestJSON = json.loads(event["body"])
            table.put_item(
                Item={
                    "team": requestJSON["team"],
                    "score": Decimal(str(requestJSON["score"])),
                    "timestamp": requestJSON["timestamp"],
                }
            )
            body = "Put item " + requestJSON["team"]
    except KeyError:
        statusCode = 400
        body = "Unsupported route: " + event["routeKey"]
    body = json.dumps(body)
    res = {
        "statusCode": statusCode,
        "headers": {"Content-Type": "application/json"},
        "body": body,
    }
    return res
