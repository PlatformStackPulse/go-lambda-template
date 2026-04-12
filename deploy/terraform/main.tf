terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    key = "go-lambda-template/terraform.tfstate"
  }
}

data "aws_caller_identity" "current" {}

locals {
  function_name           = "${var.project_name}-${var.environment}"
  table_name              = "${var.project_name}-${var.environment}-requests"
  greeting_parameter      = var.ssm_parameter_name != "" ? var.ssm_parameter_name : "/${var.project_name}/${var.environment}/greeting-prefix"
  greeting_parameter_path = trimprefix(local.greeting_parameter, "/")
  api_name                = "${local.function_name}-http-api"
  lambda_log_group        = "/aws/lambda/${local.function_name}"
  api_access_log_group    = "/aws/apigateway/${local.api_name}"
  ecr_repository_name     = var.ecr_repository_name != "" ? var.ecr_repository_name : "${local.function_name}-lambda"
  deploy_image            = var.deployment_package_type == "image"
  deploy_zip              = var.deployment_package_type == "zip"
  zip_bucket_name         = var.enable_zip_artifact_bucket ? aws_s3_bucket.artifacts[0].bucket : var.artifact_bucket_name
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Environment = var.environment
      Project     = var.project_name
      ManagedBy   = "Terraform"
    }
  }
}

resource "aws_cloudwatch_log_group" "lambda" {
  name              = local.lambda_log_group
  retention_in_days = var.log_retention_in_days
}

resource "aws_cloudwatch_log_group" "api_gateway" {
  name              = local.api_access_log_group
  retention_in_days = var.log_retention_in_days
}

resource "aws_ecr_repository" "lambda" {
  name                 = local.ecr_repository_name
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_ecr_lifecycle_policy" "lambda" {
  repository = aws_ecr_repository.lambda.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Keep the latest 10 container images"
        selection = {
          tagStatus   = "any"
          countType   = "imageCountMoreThan"
          countNumber = 10
        }
        action = {
          type = "expire"
        }
      }
    ]
  })
}

resource "aws_dynamodb_table" "requests" {
  name         = local.table_name
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "request_id"

  attribute {
    name = "request_id"
    type = "S"
  }
}

resource "aws_ssm_parameter" "greeting_prefix" {
  name  = local.greeting_parameter
  type  = "String"
  value = var.greeting_prefix
}

resource "aws_s3_bucket" "artifacts" {
  count  = local.deploy_zip && var.enable_zip_artifact_bucket ? 1 : 0
  bucket = var.artifact_bucket_name != "" ? var.artifact_bucket_name : "${local.function_name}-artifacts"
}

resource "aws_s3_bucket_versioning" "artifacts" {
  count  = local.deploy_zip && var.enable_zip_artifact_bucket ? 1 : 0
  bucket = aws_s3_bucket.artifacts[0].id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_object" "lambda_zip" {
  count  = local.deploy_zip ? 1 : 0
  bucket = local.zip_bucket_name
  key    = var.lambda_s3_key
  source = var.lambda_zip_path
  etag   = filemd5(var.lambda_zip_path)
}

resource "aws_lambda_function" "api_image" {
  count         = local.deploy_image ? 1 : 0
  function_name = local.function_name
  role          = aws_iam_role.lambda_role.arn
  package_type  = "Image"
  image_uri     = "${aws_ecr_repository.lambda.repository_url}:${var.image_tag}"
  timeout       = var.lambda_timeout
  memory_size   = var.lambda_memory_size
  architectures = [var.lambda_architecture]

  depends_on = [aws_cloudwatch_log_group.lambda]

  environment {
    variables = {
      APP_NAME                 = var.project_name
      APP_ENV                  = var.environment
      APP_VERSION              = "terraform"
      AWS_REGION               = var.aws_region
      DEBUG                    = var.debug_mode ? "true" : "false"
      DYNAMODB_TABLE_NAME      = aws_dynamodb_table.requests.name
      GREETING_PARAMETER_NAME  = aws_ssm_parameter.greeting_prefix.name
      GREETING_PREFIX_OVERRIDE = var.greeting_prefix_override
      GREETING_SOURCE_LABEL    = var.greeting_source_label
    }
  }

  tags = {
    Name = "${var.project_name}-lambda"
  }
}

resource "aws_lambda_function" "api_zip" {
  count            = local.deploy_zip ? 1 : 0
  function_name    = local.function_name
  role             = aws_iam_role.lambda_role.arn
  package_type     = "Zip"
  handler          = "bootstrap"
  runtime          = "provided.al2023"
  timeout          = var.lambda_timeout
  memory_size      = var.lambda_memory_size
  architectures    = [var.lambda_architecture]
  s3_bucket        = local.zip_bucket_name
  s3_key           = aws_s3_object.lambda_zip[0].key
  source_code_hash = filebase64sha256(var.lambda_zip_path)

  depends_on = [aws_cloudwatch_log_group.lambda]

  environment {
    variables = {
      APP_NAME                 = var.project_name
      APP_ENV                  = var.environment
      APP_VERSION              = "terraform"
      AWS_REGION               = var.aws_region
      DEBUG                    = var.debug_mode ? "true" : "false"
      DYNAMODB_TABLE_NAME      = aws_dynamodb_table.requests.name
      GREETING_PARAMETER_NAME  = aws_ssm_parameter.greeting_prefix.name
      GREETING_PREFIX_OVERRIDE = var.greeting_prefix_override
      GREETING_SOURCE_LABEL    = var.greeting_source_label
    }
  }

  lifecycle {
    precondition {
      condition     = local.zip_bucket_name != ""
      error_message = "artifact_bucket_name must be set or enable_zip_artifact_bucket must be true when deployment_package_type is zip."
    }
  }

  tags = {
    Name = "${var.project_name}-lambda"
  }
}

resource "aws_iam_role" "lambda_role" {
  name = "${var.project_name}-lambda-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_basic" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy" "lambda_app_access" {
  name = "${local.function_name}-app-access"
  role = aws_iam_role.lambda_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:DeleteItem",
          "dynamodb:Query",
          "dynamodb:Scan"
        ]
        Resource = aws_dynamodb_table.requests.arn
      },
      {
        Effect = "Allow"
        Action = [
          "ssm:GetParameter",
          "ssm:GetParameters",
          "ssm:GetParametersByPath"
        ]
        Resource = "arn:aws:ssm:${var.aws_region}:${data.aws_caller_identity.current.account_id}:parameter/${local.greeting_parameter_path}"
      }
    ]
  })
}

resource "aws_apigatewayv2_api" "http" {
  name          = local.api_name
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_integration" "lambda" {
  api_id                 = aws_apigatewayv2_api.http.id
  integration_type       = "AWS_PROXY"
  integration_uri        = local.deploy_image ? aws_lambda_function.api_image[0].invoke_arn : aws_lambda_function.api_zip[0].invoke_arn
  integration_method     = "POST"
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "hello_root" {
  api_id    = aws_apigatewayv2_api.http.id
  route_key = "GET /hello"
  target    = "integrations/${aws_apigatewayv2_integration.lambda.id}"
}

resource "aws_apigatewayv2_route" "hello_name" {
  api_id    = aws_apigatewayv2_api.http.id
  route_key = "GET /hello/{name}"
  target    = "integrations/${aws_apigatewayv2_integration.lambda.id}"
}

resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.http.id
  name        = "$default"
  auto_deploy = true

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api_gateway.arn
    format = jsonencode({
      requestId      = "$context.requestId"
      sourceIp       = "$context.identity.sourceIp"
      requestTime    = "$context.requestTime"
      httpMethod     = "$context.httpMethod"
      routeKey       = "$context.routeKey"
      status         = "$context.status"
      responseLength = "$context.responseLength"
    })
  }
}

resource "aws_lambda_permission" "allow_apigw" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = local.deploy_image ? aws_lambda_function.api_image[0].function_name : aws_lambda_function.api_zip[0].function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.http.execution_arn}/*/*"
}

output "lambda_function_arn" {
  value       = local.deploy_image ? aws_lambda_function.api_image[0].arn : aws_lambda_function.api_zip[0].arn
  description = "ARN of the Lambda function"
}

output "lambda_function_name" {
  value       = local.deploy_image ? aws_lambda_function.api_image[0].function_name : aws_lambda_function.api_zip[0].function_name
  description = "Name of the Lambda function"
}

output "api_gateway_invoke_url" {
  value       = aws_apigatewayv2_stage.default.invoke_url
  description = "Invoke URL for the default HTTP API stage"
}

output "dynamodb_table_name" {
  value       = aws_dynamodb_table.requests.name
  description = "DynamoDB table name used by the Lambda function"
}

output "greeting_parameter_name" {
  value       = aws_ssm_parameter.greeting_prefix.name
  description = "SSM parameter name that stores the greeting prefix"
}

output "ecr_repository_url" {
  value       = aws_ecr_repository.lambda.repository_url
  description = "ECR repository URL used for the default image deployment path"
}

output "deployment_package_type" {
  value       = var.deployment_package_type
  description = "Active Lambda deployment package type"
}

output "zip_artifact_bucket_name" {
  value       = local.deploy_zip ? local.zip_bucket_name : null
  description = "S3 bucket name used for the optional zip deployment path"
}
