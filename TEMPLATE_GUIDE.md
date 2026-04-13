# Go Lambda Template Guide

This template is built for teams that want to start with a working Go Lambda API and then replace the sample business logic with their own.

## What the Template Gives You

- Lambda runtime entrypoint
- API Gateway HTTP API trigger
- DynamoDB and SSM app-config adapters ready for extension
- environment variable adapter for runtime overrides
- optional Aurora PostgreSQL Serverless v2 (Data API) infrastructure
- Terraform as the default deployment path
- Docker and ECR as the default Lambda packaging path
- optional zip plus S3 deployment path
- SAM for local development
- tests, linting, security checks, and release packaging

## Customization Checklist

After creating a repository from this template, make these changes first:

1. Update the Go module path in `go.mod`.
2. Change `project_name` in `deploy/terraform/terraform.dev.tfvars`.
3. Review `api_base_path`, SSM parameter naming, DynamoDB table defaults, and any S3 or KMS values you want to expose.
4. Review the runtime environment overrides in `internal/adapter/env/runtime_settings.go`.
5. Replace the sample greeting flow in `internal/usecase/greeting.go`.
6. Update the API contract in `internal/handler/api.go` if your route or payload changes.
7. Replace the event fixture in `test/fixtures/events/apigw-request.json`.
8. If you need relational data, set `enable_postgres = true` in `deploy/terraform/terraform.dev.tfvars` and use the generated Data API environment variables.

## Extension Points

### Handler layer

`internal/handler/api.go` should stay thin. Its job is:

- read API Gateway input
- map request fields into a use case input
- map domain and application errors to HTTP responses

### Use case layer

`internal/usecase/greeting.go` is where orchestration belongs.

- call adapters
- apply domain rules
- build the response model

### Domain layer

`internal/domain/greeting.go` should contain pure logic with no AWS dependencies.

### AWS adapters

The template ships with:

- `internal/adapter/dynamodb/greeting_recorder.go`
- `internal/adapter/env/runtime_settings.go`
- `internal/adapter/postgres/greeting_recorder.go`
- `internal/adapter/ssm/parameter_store.go`

Replace or extend these adapters for your own repositories, data access patterns, or configuration needs.

The PostgreSQL adapter intentionally includes both insert and select examples so teams can copy a complete RDS Data API repository pattern instead of starting from a write-only sample.

## Suggested Customization Flow

### If your API shape stays HTTP-based

Keep API Gateway as-is and only replace the sample route behavior.

### If your payload changes

Update:

- the handler request parsing
- the fixture event
- the integration test assertions

### If your persistence model changes

Update:

- the DynamoDB adapter
- the Terraform table definition
- the IAM policy scope if you add more tables

### If your configuration model changes

Update:

- the SSM app-config document schema
- the parameter resource in Terraform and SAM
- the environment variables surfaced in Lambda

## Local Workflow

```bash
make package
make sam-start-api
curl http://127.0.0.1:3000/hello/Ada
```

If you prefer direct event invocation:

```bash
make sam-invoke
```

## Deployment Workflow

```bash
make terraform-apply
```

That is the default path. It creates or reuses the ECR repository, pushes the Docker image, and applies the Terraform stack.

For the optional zip deployment path:

```bash
make package
make terraform-apply-zip
```

## Sample Design Principles

- keep the Lambda handler small
- keep business logic outside AWS SDK code
- keep resource names and IAM scopes explicit
- prefer fixture-driven integration tests over ad hoc manual verification

## Twelve-Factor Expectations

- keep config in environment variables and SSM, not hardcoded constants
- treat DynamoDB, SSM, and optional PostgreSQL as replaceable backing resources
- keep handler processes stateless and push mutable data to backing services
- separate build (`make package` or image build), release (Terraform apply), and run (Lambda)
- emit logs to stdout as structured events and avoid file-based logging

## When to Add More Services

Only add more infrastructure when the project really needs it. This template starts with API Gateway, DynamoDB, SSM, and CloudWatch because they are common for Lambda-backed APIs. Do not grow it into a kitchen sink template.
