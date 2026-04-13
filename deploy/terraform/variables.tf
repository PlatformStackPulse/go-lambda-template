variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment name"
  type        = string
  validation {
    condition     = contains(["dev", "staging", "prod"], var.environment)
    error_message = "Environment must be one of: dev, staging, prod"
  }
}

variable "project_name" {
  description = "Project name used for AWS resources"
  type        = string
  default     = "go-lambda-template"
}

variable "debug_mode" {
  description = "Enable debug mode"
  type        = bool
  default     = false
}

variable "deployment_package_type" {
  description = "Default Lambda deployment package type"
  type        = string
  default     = "image"

  validation {
    condition     = contains(["image", "zip"], var.deployment_package_type)
    error_message = "deployment_package_type must be one of: image, zip"
  }
}

variable "image_tag" {
  description = "Container image tag used for the ECR deployment path"
  type        = string
  default     = "latest"
}

variable "ecr_repository_name" {
  description = "Optional explicit ECR repository name"
  type        = string
  default     = ""
}

variable "lambda_zip_path" {
  description = "Path to the packaged Lambda zip artifact"
  type        = string
  default     = "../../dist/lambda.zip"
}

variable "lambda_architecture" {
  description = "Lambda architecture"
  type        = string
  default     = "arm64"

  validation {
    condition     = contains(["arm64", "x86_64"], var.lambda_architecture)
    error_message = "Lambda architecture must be one of: arm64, x86_64"
  }
}

variable "lambda_timeout" {
  description = "Lambda timeout in seconds"
  type        = number
  default     = 10
}

variable "lambda_memory_size" {
  description = "Lambda memory size in MB"
  type        = number
  default     = 512
}

variable "log_retention_in_days" {
  description = "CloudWatch log retention in days"
  type        = number
  default     = 14
}

variable "api_base_path" {
  description = "Base path exposed by API Gateway for the sample HTTP routes"
  type        = string
  default     = "/hello"

  validation {
    condition     = startswith(var.api_base_path, "/")
    error_message = "api_base_path must start with /."
  }
}

variable "api_source_label" {
  description = "Optional response source label surfaced through Lambda environment variables"
  type        = string
  default     = ""
}

variable "ssm_parameter_prefix" {
  description = "Optional explicit SSM parameter prefix for platform configuration"
  type        = string
  default     = ""
}

variable "app_config_parameter_name" {
  description = "Optional explicit SSM parameter name that stores the JSON app configuration document"
  type        = string
  default     = ""
}

variable "s3_source_bucket_name" {
  description = "Optional source S3 bucket name exposed to the runtime"
  type        = string
  default     = ""
}

variable "s3_source_key_prefix" {
  description = "Optional source S3 key prefix exposed to the runtime"
  type        = string
  default     = ""
}

variable "s3_target_bucket_name" {
  description = "Optional target S3 bucket name exposed to the runtime"
  type        = string
  default     = ""
}

variable "s3_target_key_prefix" {
  description = "Optional target S3 key prefix exposed to the runtime"
  type        = string
  default     = ""
}

variable "kms_key_arn" {
  description = "Optional KMS key ARN exposed to the runtime"
  type        = string
  default     = ""
}

variable "enable_zip_artifact_bucket" {
  description = "Create and manage an S3 bucket for optional zip deployment artifacts"
  type        = bool
  default     = false
}

variable "artifact_bucket_name" {
  description = "Existing or managed S3 bucket name for the optional zip deployment path"
  type        = string
  default     = ""
}

variable "lambda_s3_key" {
  description = "S3 object key used for the optional zip deployment path"
  type        = string
  default     = "lambda/lambda.zip"
}

variable "enable_postgres" {
  description = "Enable Aurora PostgreSQL Serverless v2 (Data API) as an additional backing service"
  type        = bool
  default     = false
}

variable "postgres_engine_version" {
  description = "Aurora PostgreSQL engine version"
  type        = string
  default     = "16.4"
}

variable "postgres_database_name" {
  description = "Default database name for Aurora PostgreSQL"
  type        = string
  default     = "app"
}

variable "postgres_master_username" {
  description = "Master username for Aurora PostgreSQL"
  type        = string
  default     = "appadmin"
}

variable "postgres_min_acu" {
  description = "Minimum Aurora Serverless v2 capacity units"
  type        = number
  default     = 0.5
}

variable "postgres_max_acu" {
  description = "Maximum Aurora Serverless v2 capacity units"
  type        = number
  default     = 2
}

variable "postgres_backup_retention_days" {
  description = "Backup retention for Aurora PostgreSQL"
  type        = number
  default     = 7
}

variable "postgres_deletion_protection" {
  description = "Enable deletion protection for Aurora PostgreSQL"
  type        = bool
  default     = false
}

variable "postgres_skip_final_snapshot" {
  description = "Skip final snapshot when deleting Aurora PostgreSQL"
  type        = bool
  default     = true
}
