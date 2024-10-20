# Dynamo DB
resource "aws_dynamodb_table" "portal" {
  name         = "${var.name}-portal"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "team"
  range_key    = "timestamp"

  attribute {
    name = "team"
    type = "S"
  }

  attribute {
    name = "timestamp"
    type = "S"
  }

  tags = {
    Name = "${var.name}-portal"
  }
}

# Lambda
data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}
data "aws_iam_policy_document" "lambda" {
  statement {
    effect    = "Allow"
    resources = [aws_dynamodb_table.portal.arn]
    actions = [
      "dynamodb:DeleteItem",
      "dynamodb:GetItem",
      "dynamodb:PutItem",
      "dynamodb:Scan",
      "dynamodb:UpdateItem",
      "dynamodb:BatchWriteItem",
    ]
  }
  statement {
    effect    = "Allow"
    resources = ["*"]
    actions = [
      "logs:CreateLogGroup",
    ]
  }
  statement {
    effect    = "Allow"
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/*:*"]
    actions = [
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
  }
}
resource "aws_iam_role" "iam_for_lambda" {
  name               = "${var.name}-portal-lambda"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}
resource "aws_iam_policy" "lambda" {
  name        = "${var.name}-portal-lambda"
  description = "Policy for the portal Lambda function"
  policy      = data.aws_iam_policy_document.lambda.json
}
resource "aws_iam_role_policy_attachment" "lambda" {
  role       = aws_iam_role.iam_for_lambda.name
  policy_arn = aws_iam_policy.lambda.arn
}
data "archive_file" "lambda" {
  type        = "zip"
  output_path = "lambda_function.zip"
  source_content = templatefile("${path.module}/portal_lambda.py.tpl", {
    dynamodb_table_name = aws_dynamodb_table.portal.name
  })
  source_content_filename = "lambda.py"
}
resource "aws_lambda_function" "portal" {
  function_name = "${var.name}-portal"
  description   = "Lambda function to get/put/delete scores for each team"
  filename      = "lambda_function.zip"
  role          = aws_iam_role.iam_for_lambda.arn
  handler       = "lambda.lambda_handler"

  source_code_hash = data.archive_file.lambda.output_base64sha256

  runtime     = "python3.12"
  timeout     = 10
  memory_size = 256
}
resource "aws_lambda_permission" "apigw_lambda" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.portal.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_apigatewayv2_api.portal.id}/*/*/teams"
}
resource "aws_cloudwatch_log_group" "lambda" {
  name = "/aws/lambda/${aws_lambda_function.portal.function_name}"
}

# API Gateway
resource "aws_apigatewayv2_api" "portal" {
  name          = "${var.name}-portal"
  description   = "API Gateway for the portal Lambda function"
  protocol_type = "HTTP"
  cors_configuration {
    allow_origins = ["http://${aws_s3_bucket_website_configuration.portal.website_endpoint}"]
  }
}
resource "aws_apigatewayv2_integration" "portal" {
  api_id           = aws_apigatewayv2_api.portal.id
  integration_type = "AWS_PROXY"

  connection_type        = "INTERNET"
  integration_method     = "POST"
  integration_uri        = aws_lambda_function.portal.invoke_arn
  passthrough_behavior   = "WHEN_NO_MATCH"
  payload_format_version = "2.0"
}
resource "aws_apigatewayv2_route" "portal_teams_get" {
  api_id    = aws_apigatewayv2_api.portal.id
  route_key = "GET /teams"
  target    = "integrations/${aws_apigatewayv2_integration.portal.id}"
}
resource "aws_apigatewayv2_route" "portal_teams_put" {
  api_id    = aws_apigatewayv2_api.portal.id
  route_key = "PUT /teams"
  target    = "integrations/${aws_apigatewayv2_integration.portal.id}"
}
resource "aws_apigatewayv2_route" "portal_teams_delete" {
  api_id    = aws_apigatewayv2_api.portal.id
  route_key = "DELETE /teams"
  target    = "integrations/${aws_apigatewayv2_integration.portal.id}"
}
resource "aws_apigatewayv2_deployment" "portal" {
  api_id = aws_apigatewayv2_api.portal.id
  triggers = {
    redeployment = sha1(join(",", tolist([
      jsonencode(aws_apigatewayv2_integration.portal),
      jsonencode(aws_apigatewayv2_route.portal_teams_get),
      jsonencode(aws_apigatewayv2_route.portal_teams_put),
      jsonencode(aws_apigatewayv2_route.portal_teams_delete),
    ])))
  }
  lifecycle {
    create_before_destroy = true
  }
}
resource "aws_apigatewayv2_stage" "portal" {
  api_id      = aws_apigatewayv2_api.portal.id
  name        = "$default"
  auto_deploy = true
}

# S3
resource "aws_s3_bucket" "portal" {
  bucket_prefix = "${var.name}-portal"
  tags = {
    Name = "${var.name}-portal"
  }
}
resource "aws_s3_bucket_public_access_block" "portal" {
  bucket = aws_s3_bucket.portal.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}
resource "aws_s3_bucket_website_configuration" "portal" {
  bucket = aws_s3_bucket.portal.id
  index_document {
    suffix = "index.html"
  }
}
resource "aws_s3_bucket_policy" "portal" {
  bucket = aws_s3_bucket.portal.bucket

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect    = "Allow"
        Principal = "*"
        Action    = "s3:GetObject"
        Resource  = "${aws_s3_bucket.portal.arn}/*"
      }
    ]
  })

  depends_on = [aws_s3_bucket_public_access_block.portal]
}
resource "aws_s3_object" "portal_index_html" {
  bucket       = aws_s3_bucket.portal.bucket
  key          = "index.html"
  source       = "${path.module}/../../portal/index.html"
  content_type = "text/html"

  depends_on = [
    aws_s3_bucket_policy.portal,
    aws_s3_bucket_public_access_block.portal
  ]
}
resource "local_file" "main_js_with_replaced_variables" {
  filename = "${path.module}/dist/main_with_variables.js"
  content  = replace(file("${path.module}/../../portal/dist/main.js"), "<<API_GATEWAY_DOMAIN_URL>>", aws_apigatewayv2_stage.portal.invoke_url)
}
resource "aws_s3_object" "portal_main_js" {
  bucket       = aws_s3_bucket.portal.bucket
  key          = "dist/main.js"
  content      = local_file.main_js_with_replaced_variables.content
  content_type = "application/javascript"

  depends_on = [
    aws_s3_bucket_policy.portal,
    aws_s3_bucket_public_access_block.portal
  ]
}
