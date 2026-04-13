# Workflow Guide

This repository uses one local development path and one infrastructure path.

## Tool Responsibilities

- AWS SAM: local build, local invoke, and local API simulation
- Terraform: default infrastructure provisioning and deployment path
- Docker and `scripts/push-ecr.sh`: default Lambda image packaging and ECR push path
- GitHub Actions: lint, test, security scanning, packaging, and release publishing

## Daily Development Flow

1. make your code changes in `internal/handler`, `internal/usecase`, `internal/domain`, or the AWS adapters
2. run `make test`
3. run `make sam-start-api` or `make sam-invoke`
4. deploy through `make terraform-apply` when ready

## Local Development Commands

```bash
make test
make package
make sam-build
make sam-invoke
make sam-start-api
```

## Deployment Commands

```bash
make terraform-plan
make terraform-apply
make terraform-plan-zip
make terraform-apply-zip
```

Default deployment path:

- `make terraform-apply` creates the ECR repository if needed, builds the Docker image, pushes it to ECR, and applies the Terraform stack

Optional deployment path:

- `make package && make terraform-apply-zip` uploads the zip artifact through S3-backed Terraform resources

## Branch Protection

Recommended required checks on `main`:

1. `CI Pipeline / Lint & Format Check`
2. `CI Pipeline / Test (1.22)`
3. `CI Pipeline / Test (1.23)`
4. `CI Pipeline / Security Scans`
5. `CI Pipeline / Commit Lint`
6. `CI Pipeline / Build`
7. `CodeQL Analysis / Analyze`

Apply via script:

```bash
export GITHUB_TOKEN=ghp_xxx
scripts/apply-branch-protection.sh
```

Optional override:

```bash
GITHUB_OWNER=PlatformStackPulse GITHUB_REPO=go-lambda-template BRANCH=main scripts/apply-branch-protection.sh
```

## CI and Release Behavior

- CI runs lint, tests, security checks, builds the Lambda artifact, and validates the SAM template
- the repository is designed around Docker plus ECR deployment for Terraform-managed environments
- zip packaging remains available as an optional deployment path

## Release Flow

1. merge changes to `main`
2. create a version tag such as `v1.2.3`
3. push the tag
4. let the release workflow publish the Lambda artifacts

## Operational Defaults

- Lambda logs are JSON structured through `slog`
- CloudWatch log retention is managed in Terraform
- API Gateway access logging is enabled in Terraform
- the sample Lambda role includes scoped DynamoDB and SSM app-config access
- the sample Lambda runtime reads platform env vars and an SSM app-config document through dedicated adapters
- optional Aurora PostgreSQL Serverless v2 (Data API) can be enabled through Terraform variables
- the template expects Twelve-Factor practices: env-driven config, stateless handlers, and backing-service resource abstraction
