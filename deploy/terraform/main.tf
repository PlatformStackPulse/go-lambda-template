terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    key = "go-template/terraform.tfstate"
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Environment = var.environment
      Project     = "go-template"
      ManagedBy   = "Terraform"
    }
  }
}

# Example deployment
resource "aws_lambda_function" "go_template" {
  filename         = data.archive_file.lambda_zip.output_path
  source_code_hash = data.archive_file.lambda_zip.output_base64sha256
  function_name    = "${var.project_name}-${var.environment}"
  role             = aws_iam_role.lambda_role.arn
  handler          = "bootstrap"
  runtime          = "provided.al2"
  timeout          = 60
  memory_size      = 512

  environment {
    variables = {
      DEBUG = var.debug_mode ? "true" : "false"
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

data "archive_file" "lambda_zip" {
  type        = "zip"
  source_file = "bin/go-template"
  output_path = "lambda_function.zip"
}

output "lambda_function_arn" {
  value       = aws_lambda_function.go_template.arn
  description = "ARN of the Lambda function"
}

output "lambda_function_name" {
  value       = aws_lambda_function.go_template.function_name
  description = "Name of the Lambda function"
}
