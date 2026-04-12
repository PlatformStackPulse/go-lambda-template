# Go Lambda Template

![Go Version](https://img.shields.io/badge/Go-1.23+-blue?style=flat-square&logo=go)
![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)
![DevContainer](https://img.shields.io/static/v1?label=DevContainer&message=Ready&color=blue&style=flat-square&logo=visual-studio-code)

Slim, production-ready GitHub template for building Go AWS Lambda APIs.

## Overview

This template is opinionated by default:

- Lambda behind API Gateway HTTP API
- DynamoDB table provisioned for application data
- SSM Parameter Store used for runtime configuration
- Lambda environment variables exposed through a dedicated runtime adapter
- CloudWatch logs configured with retention
- Terraform as the default deployment path
- Docker plus ECR for the default Lambda package
- optional zip deployment through S3
- AWS SAM for local build and local API invocation

The sample implementation returns a greeting from API Gateway, reads the greeting prefix from SSM, records each request in DynamoDB, and logs structured JSON to CloudWatch.

## What Developers Edit

Most teams only need to touch these files first:

- `internal/handler/api.go`
- `internal/usecase/greeting.go`
- `internal/domain/greeting.go`
- `test/fixtures/events/apigw-request.json`

Everything else is there to keep packaging, infra, and local workflows out of the way.

## Project Layout

```text
cmd/lambda/main.go                 # Lambda entrypoint
internal/app/app.go                # Dependency wiring
internal/handler/api.go            # API Gateway adapter
internal/usecase/greeting.go       # Business orchestration
internal/domain/greeting.go        # Pure domain logic
internal/adapter/dynamodb/         # DynamoDB integration
internal/adapter/ssm/              # SSM integration
deploy/terraform/                  # Production infrastructure
deploy/sam/template.yaml           # Local SAM workflow and quick-start deploy
scripts/push-ecr.sh                # Push Docker image to ECR
test/fixtures/events/              # Sample API Gateway events
test/integration/lambda/           # Handler integration tests
```

## Quick Start

### 1. Create a repository from the template

```bash
gh repo create my-lambda-api --template PlatformStackPulse/go-lambda-template
cd my-lambda-api
```

### 2. Install local tools

```bash
make dev-setup
```

Also install these separately if they are not already available:

- AWS SAM CLI
- Terraform
- Docker
- AWS CLI

### 3. Run tests

```bash
make test
```

### 4. Run the default deployment flow with Terraform

```bash
make terraform-apply
```

This default flow does all of the following:

- creates the ECR repository with Terraform if needed
- builds the Lambda container image with Docker
- pushes the image to ECR with the provided bash script
- applies the full Lambda, API Gateway, DynamoDB, SSM, and CloudWatch stack with Terraform

### 5. Run the API locally with SAM

```bash
make sam-start-api
```

Then call the sample endpoint:

```bash
curl http://127.0.0.1:3000/hello/Ada
```

### 6. Optional zip deployment path

```bash
make package
make terraform-apply-zip
```

Use the zip path only if you explicitly want S3-backed Lambda artifacts instead of the default ECR image flow.

### 7. Call the deployed API

After apply, call the deployed API using the Terraform output:

```bash
terraform -chdir=deploy/terraform output -raw api_gateway_invoke_url
curl "$(terraform -chdir=deploy/terraform output -raw api_gateway_invoke_url)/hello/Ada"
```

## Sample Behavior

The sample flow is:

1. API Gateway invokes Lambda on `GET /hello` or `GET /hello/{name}`.
2. Lambda reads the greeting prefix from SSM Parameter Store.
3. Lambda reads optional runtime overrides from environment variables.
4. Lambda writes a request record to DynamoDB.
5. Lambda returns JSON like this:

```json
{
  "message": "Hello, Ada!",
  "request_id": "...",
  "source": "$default",
  "timestamp": "2026-04-12T12:00:00Z"
}
```

## Example Customization

If you want to replace the sample greeting logic with your own business flow:

1. Change the API contract in `internal/handler/api.go`.
2. Replace the use case in `internal/usecase/greeting.go`.
3. Update the environment adapter in `internal/adapter/env/runtime_settings.go` if you want different env-driven behavior.
4. Update the DynamoDB and SSM adapters if your data model changes.
5. Update the fixture in `test/fixtures/events/apigw-request.json`.

Example: keeping the same route but returning a product-specific message.

```go
func (uc *GreetingUseCase) Execute(ctx context.Context, input GreetingInput) (GreetingOutput, error) {
	prefix, err := uc.configProvider.GreetingPrefix(ctx)
	if err != nil {
		return GreetingOutput{}, err
	}

	message := fmt.Sprintf("%s from %s", prefix, domain.NormalizeName(input.Name))
	return GreetingOutput{Message: message, RequestID: input.RequestID, Source: input.Source}, nil
}
```

## Environment Variables

The template uses these Lambda environment variables:

| Variable | Purpose | Default |
| --- | --- | --- |
| `APP_NAME` | Logical application name | `go-lambda-template` |
| `APP_ENV` | Environment name | `dev` |
| `APP_VERSION` | Build or release version | `dev` |
| `AWS_REGION` | AWS region injected by Lambda runtime (reserved) | runtime-provided |
| `DEBUG` | Debug logging flag | `false` |
| `DYNAMODB_TABLE_NAME` | DynamoDB table used by the sample adapter | `${APP_NAME}-${APP_ENV}-requests` |
| `GREETING_PARAMETER_NAME` | SSM parameter that stores the greeting prefix | `/${APP_NAME}/${APP_ENV}/greeting-prefix` |
| `GREETING_PREFIX_OVERRIDE` | Optional env var override for the greeting prefix | `""` |
| `GREETING_SOURCE_LABEL` | Optional env var override for the response source field | `""` |

## Commands

```bash
make build           # Build the bootstrap binary
make package         # Build and zip the Lambda artifact
make docker-build    # Build the Lambda container image
make terraform-bootstrap # Create the ECR repository before the first push
make ecr-push        # Push the image to ECR
make test            # Run all tests
make lint            # Run golangci-lint
make security        # Run gosec and govulncheck
make sam-build       # Prepare the SAM application
make sam-invoke      # Invoke the Lambda locally with a sample event
make sam-start-api   # Start the local API Gateway emulator
make terraform-plan  # Plan the default image-based deployment
make terraform-apply # Build, push, and deploy the default image-based stack
make terraform-plan-zip  # Plan the optional zip deployment
make terraform-apply-zip # Deploy the optional zip + S3 path
```

## Infrastructure Defaults

Terraform provisions these resources by default:

- ECR repository for the Lambda container image
- Lambda function deployed from ECR by default
- API Gateway HTTP API with `GET /hello` and `GET /hello/{name}`
- DynamoDB table for request records
- SSM parameter for the greeting prefix
- CloudWatch log groups for Lambda and API Gateway access logs
- IAM role with Lambda basic execution plus scoped DynamoDB and SSM access plus Lambda environment variables

SAM provides a matching local workflow. Terraform is the first and default deployment path. The zip + S3 deployment path is supported, but optional.

## Testing

The repository includes:

- unit tests for config, domain logic, use case orchestration, handler mapping, and AWS adapters
- integration tests for the API Gateway handler using a fixture event

Run everything with:

```bash
make test
```

## How To Use This Template

Use this sequence when turning the template into your own Lambda service:

1. Install AWS SAM CLI locally and run the sample API end to end:

```bash
make sam-start-api
curl http://127.0.0.1:3000/hello/Ada
```

2. Edit the sample business logic in `internal/usecase/greeting.go` and the request mapping in `internal/handler/api.go` so the Lambda matches your real API behavior.

3. Set your AWS values in `deploy/terraform/terraform.dev.tfvars` and deploy the default image-based stack:

```bash
make terraform-apply
```

This gives you a practical flow:

- verify the template locally with SAM
- replace the sample logic with your own logic
- deploy the real stack with Terraform

## Next Steps After Creating Your Repo

1. Rename the module in `go.mod`.
2. Rename the AWS resource defaults in Terraform and SAM.
3. Replace the sample greeting use case with your own business logic.
4. Adjust DynamoDB schema and SSM parameter names for your project.
5. Update the example event fixture and integration tests.

More detailed customization notes are in `TEMPLATE_GUIDE.md` and `WORKFLOW.md`.
