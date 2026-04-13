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
  function_name             = "${var.project_name}-${var.environment}"
  table_name                = "${var.project_name}-${var.environment}-requests"
  ssm_parameter_prefix      = var.ssm_parameter_prefix != "" ? trimsuffix(var.ssm_parameter_prefix, "/") : "/${var.project_name}/${var.environment}"
  app_config_parameter      = var.app_config_parameter_name != "" ? var.app_config_parameter_name : "${local.ssm_parameter_prefix}/app-config"
  app_config_parameter_path = trimprefix(local.app_config_parameter, "/")
  api_name                  = "${local.function_name}-http-api"
  lambda_log_group          = "/aws/lambda/${local.function_name}"
  api_access_log_group      = "/aws/apigateway/${local.api_name}"
  ecr_repository_name       = var.ecr_repository_name != "" ? var.ecr_repository_name : "${local.function_name}-lambda"
  deploy_image              = var.deployment_package_type == "image"
  deploy_zip                = var.deployment_package_type == "zip"
  zip_bucket_name           = var.enable_zip_artifact_bucket ? aws_s3_bucket.artifacts[0].bucket : var.artifact_bucket_name
  postgres_cluster_name     = "${local.function_name}-aurora-pg"
  postgres_enabled          = var.enable_postgres
  postgres_cluster_arn      = local.postgres_enabled ? aws_rds_cluster.postgres[0].arn : ""
  postgres_secret_arn       = local.postgres_enabled ? aws_rds_cluster.postgres[0].master_user_secret[0].secret_arn : ""
  app_config_document = jsonencode({
    "api.base_path"                = var.api_base_path
    "dynamodb.requests_table_name" = local.table_name
    "kms.key_arn"                  = var.kms_key_arn
    "platform.environment"         = var.environment
    "platform.name"                = var.project_name
    "platform.version"             = "terraform"
    "s3.source.bucket"             = var.s3_source_bucket_name
    "s3.source.key_prefix"         = var.s3_source_key_prefix
    "s3.target.bucket"             = var.s3_target_bucket_name
    "s3.target.key_prefix"         = var.s3_target_key_prefix
    "sample.greeting.prefix"       = "Hello"
  })
}

data "aws_vpc" "default" {
  count   = local.postgres_enabled ? 1 : 0
  default = true
}

data "aws_subnets" "default" {
  count = local.postgres_enabled ? 1 : 0

  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default[0].id]
  }
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

resource "aws_ssm_parameter" "app_config" {
  name  = local.app_config_parameter
  type  = "String"
  value = local.app_config_document
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

resource "aws_db_subnet_group" "postgres" {
  count       = local.postgres_enabled ? 1 : 0
  name        = "${local.postgres_cluster_name}-subnets"
  subnet_ids  = data.aws_subnets.default[0].ids
  description = "Subnet group for ${local.postgres_cluster_name}"
}

resource "aws_security_group" "postgres" {
  count       = local.postgres_enabled ? 1 : 0
  name        = "${local.postgres_cluster_name}-sg"
  description = "Security group for ${local.postgres_cluster_name}"
  vpc_id      = data.aws_vpc.default[0].id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_rds_cluster" "postgres" {
  count                               = local.postgres_enabled ? 1 : 0
  cluster_identifier                  = local.postgres_cluster_name
  engine                              = "aurora-postgresql"
  engine_version                      = var.postgres_engine_version
  database_name                       = var.postgres_database_name
  master_username                     = var.postgres_master_username
  manage_master_user_password         = true
  storage_encrypted                   = true
  db_subnet_group_name                = aws_db_subnet_group.postgres[0].name
  vpc_security_group_ids              = [aws_security_group.postgres[0].id]
  backup_retention_period             = var.postgres_backup_retention_days
  deletion_protection                 = var.environment == "prod" ? true : var.postgres_deletion_protection
  skip_final_snapshot                 = var.postgres_skip_final_snapshot
  copy_tags_to_snapshot               = true
  iam_database_authentication_enabled = true
  enable_http_endpoint                = true

  serverlessv2_scaling_configuration {
    min_capacity = var.postgres_min_acu
    max_capacity = var.postgres_max_acu
  }
}

resource "aws_rds_cluster_instance" "postgres" {
  count                = local.postgres_enabled ? 1 : 0
  identifier           = "${local.postgres_cluster_name}-instance-1"
  cluster_identifier   = aws_rds_cluster.postgres[0].id
  instance_class       = "db.serverless"
  engine               = aws_rds_cluster.postgres[0].engine
  engine_version       = aws_rds_cluster.postgres[0].engine_version
  db_subnet_group_name = aws_db_subnet_group.postgres[0].name
  publicly_accessible  = false
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
      APP_NAME                     = var.project_name
      APP_ENV                      = var.environment
      APP_VERSION                  = "terraform"
      DEBUG                        = var.debug_mode ? "true" : "false"
      API_BASE_PATH                = var.api_base_path
      API_SOURCE_LABEL             = var.api_source_label
      DYNAMODB_REQUESTS_TABLE_NAME = aws_dynamodb_table.requests.name
      SSM_PARAMETER_PREFIX         = local.ssm_parameter_prefix
      APP_CONFIG_PARAMETER_NAME    = aws_ssm_parameter.app_config.name
      S3_SOURCE_BUCKET_NAME        = var.s3_source_bucket_name
      S3_SOURCE_KEY_PREFIX         = var.s3_source_key_prefix
      S3_TARGET_BUCKET_NAME        = var.s3_target_bucket_name
      S3_TARGET_KEY_PREFIX         = var.s3_target_key_prefix
      KMS_KEY_ARN                  = var.kms_key_arn
      POSTGRES_ENABLED             = local.postgres_enabled ? "true" : "false"
      POSTGRES_DATABASE_NAME       = var.postgres_database_name
      POSTGRES_CLUSTER_ARN         = local.postgres_cluster_arn
      POSTGRES_SECRET_ARN          = local.postgres_secret_arn
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
      APP_NAME                     = var.project_name
      APP_ENV                      = var.environment
      APP_VERSION                  = "terraform"
      DEBUG                        = var.debug_mode ? "true" : "false"
      API_BASE_PATH                = var.api_base_path
      API_SOURCE_LABEL             = var.api_source_label
      DYNAMODB_REQUESTS_TABLE_NAME = aws_dynamodb_table.requests.name
      SSM_PARAMETER_PREFIX         = local.ssm_parameter_prefix
      APP_CONFIG_PARAMETER_NAME    = aws_ssm_parameter.app_config.name
      S3_SOURCE_BUCKET_NAME        = var.s3_source_bucket_name
      S3_SOURCE_KEY_PREFIX         = var.s3_source_key_prefix
      S3_TARGET_BUCKET_NAME        = var.s3_target_bucket_name
      S3_TARGET_KEY_PREFIX         = var.s3_target_key_prefix
      KMS_KEY_ARN                  = var.kms_key_arn
      POSTGRES_ENABLED             = local.postgres_enabled ? "true" : "false"
      POSTGRES_DATABASE_NAME       = var.postgres_database_name
      POSTGRES_CLUSTER_ARN         = local.postgres_cluster_arn
      POSTGRES_SECRET_ARN          = local.postgres_secret_arn
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
        Resource = "arn:aws:ssm:${var.aws_region}:${data.aws_caller_identity.current.account_id}:parameter/${local.app_config_parameter_path}"
      }
    ]
  })
}

resource "aws_iam_role_policy" "lambda_postgres_access" {
  count = local.postgres_enabled ? 1 : 0
  name  = "${local.function_name}-postgres-access"
  role  = aws_iam_role.lambda_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "rds-data:ExecuteStatement",
          "rds-data:BatchExecuteStatement",
          "rds-data:BeginTransaction",
          "rds-data:CommitTransaction",
          "rds-data:RollbackTransaction"
        ]
        Resource = local.postgres_cluster_arn
      },
      {
        Effect = "Allow"
        Action = [
          "secretsmanager:GetSecretValue"
        ]
        Resource = local.postgres_secret_arn
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

resource "aws_apigatewayv2_route" "sample_base" {
  api_id    = aws_apigatewayv2_api.http.id
  route_key = "GET ${var.api_base_path}"
  target    = "integrations/${aws_apigatewayv2_integration.lambda.id}"
}

resource "aws_apigatewayv2_route" "sample_with_name" {
  api_id    = aws_apigatewayv2_api.http.id
  route_key = "GET ${var.api_base_path}/{name}"
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

output "app_config_parameter_name" {
  value       = aws_ssm_parameter.app_config.name
  description = "SSM parameter name that stores the JSON app configuration document"
}

output "api_base_path" {
  value       = var.api_base_path
  description = "API Gateway base path used by the sample routes"
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

output "postgres_enabled" {
  value       = local.postgres_enabled
  description = "Whether Aurora PostgreSQL is enabled"
}

output "postgres_cluster_arn" {
  value       = local.postgres_enabled ? local.postgres_cluster_arn : null
  description = "Aurora PostgreSQL cluster ARN for Data API usage"
}

output "postgres_secret_arn" {
  value       = local.postgres_enabled ? local.postgres_secret_arn : null
  description = "Secrets Manager ARN for Aurora PostgreSQL credentials"
}

output "postgres_database_name" {
  value       = local.postgres_enabled ? var.postgres_database_name : null
  description = "Aurora PostgreSQL default database name"
}
