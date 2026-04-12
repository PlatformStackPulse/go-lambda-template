# Contributing to Go Lambda Template

## Development Principles

- keep the template Lambda-first
- do not reintroduce generic CLI or container deployment scaffolding
- keep the handler thin and push logic into use cases and domain packages
- add AWS services only when they are broadly useful for Lambda-backed APIs
- keep Terraform as the default deployment path

## Local Setup

```bash
make dev-setup
make test
make terraform-plan
```

Install AWS SAM CLI, Terraform, Docker, and the AWS CLI separately.

## Contribution Workflow

1. create a branch
2. make focused changes
3. run `make fmt lint test security`
4. if infrastructure or local development changed, run `make sam-build`
5. open a pull request with a Conventional Commit title

## Areas Most Likely to Change

- `internal/handler/`
- `internal/usecase/`
- `internal/adapter/`
- `deploy/terraform/`
- `deploy/sam/`
- `test/`

## What to Avoid

- adding second-class deployment paths that compete with Terraform
- adding extra event sources into the default scaffold without a strong reason
- widening IAM permissions unnecessarily
- coupling business logic directly to AWS SDK clients
- making zip deployment the primary path again

## Documentation Expectations

If you change runtime behavior, local workflow, or infrastructure defaults, update:

- `README.md`
- `TEMPLATE_GUIDE.md`
- `WORKFLOW.md`

## Security and Review

- prefer scoped IAM policies over broad wildcards
- keep SSM access explicit
- keep DynamoDB access tied to the provisioned table where possible
- ensure sample code does not encourage unsafe secret handling

## Release Notes

Tagged releases publish the Lambda artifact and SBOM through GitHub Actions. If your change affects packaging or deployment, mention that in the pull request description.
