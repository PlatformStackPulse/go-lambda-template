# AWS Lambda Template Conversion Plan

## Objective

Convert this repository from a general-purpose Go CLI/API starter into a Lambda-first GitHub template that lets developers focus on business logic only.

The resulting template should:

- start as an AWS Lambda project by default
- expose the function through API Gateway by default
- include first-class access patterns for DynamoDB and AWS Systems Manager Parameter Store
- include CloudWatch logging as part of the default runtime and IAM baseline
- provide clear seams for handler logic, business logic, config, logging, and error handling
- ship with Lambda-specific build, test, and deployment workflows
- remove generic template baggage that creates noise for Lambda users

## Current Repository Digestion

The repository is currently a solid Go foundation, but it is not yet a Lambda template.

What is already useful:

- `internal/domain`, `internal/usecase`, `internal/errors`, `internal/config`, `internal/logger`
- unit and integration test structure
- Terraform already includes an `aws_lambda_function`
- CI, security scanning, versioning, and release automation are in place

What is currently misaligned with the goal:

- the app starts as a Cobra CLI from `cmd/app/main.go`
- `internal/cli` and the `hello` example command are CLI-only scaffolding
- the Makefile, release script, and Dockerfile target generic binaries and container publishing
- Kubernetes manifests and container-centric docs are not relevant for a Lambda-first template
- README and template guidance still describe a general Go template, not a Lambda developer experience

## Recommendation

Make this repository Lambda-first, not dual-mode.

Keeping both CLI mode and Lambda mode would preserve flexibility, but it would also preserve unnecessary complexity in the template. A dedicated Lambda template should be opinionated. The developer experience should begin with:

1. define the event input
2. implement handler/usecase logic
3. run tests
4. build a Lambda zip
5. deploy with Terraform

## AWS SAM Position

AWS SAM should be part of this template, but not as a second full infrastructure system competing with Terraform.

Recommended role for SAM:

- local build and local invocation
- fast feedback for handler development
- optional guided deploy path for developers who want a quick start

Recommended role for Terraform:

- canonical infrastructure definition
- IAM, log groups, function configuration, and environment-specific deployment
- production-oriented deployment path

Reasoning:

- SAM is excellent for developer ergonomics around Lambda packaging, local invoke, and event testing
- Terraform is better for broader infrastructure ownership and teams already managing AWS through IaC
- treating both as equal deployment systems would duplicate configuration and drift over time

So the template should support SAM, but it should be explicit that:

1. SAM is the local developer workflow layer
2. Terraform is the primary infrastructure layer

## Default Platform Assumptions

This template should ship with one opinionated default stack.

Default baseline:

- Lambda behind API Gateway
- DynamoDB available as the default persistence option
- SSM Parameter Store available for runtime configuration and secrets references
- CloudWatch Logs enabled with explicit log group management and retention

What this means in practice:

- the default handler should use API Gateway request and response types
- the template should include a DynamoDB repository seam by default
- the template should include an SSM client seam by default
- the default IAM role should grant the minimum permissions required for logs, DynamoDB access, and reading SSM parameters

This is a better default than a generic JSON event because it matches a very common production Lambda shape while still being understandable for template users.

## Target Template Shape

The template should converge on a structure like this:

```text
cmd/
  lambda/
    main.go                 # lambda.Start(...)
internal/
  handler/
    api.go                 # API Gateway request/response mapping
  app/
    app.go                 # wiring config, logger, usecases
  adapter/
    dynamodb/
    ssm/
  config/
  domain/
  errors/
  logger/
  usecase/
pkg/
  version/
deploy/
  sam/
    template.yaml            # local build/invoke and optional quick deploy
  terraform/
test/
  fixtures/events/
  integration/lambda/
```

The key rule is that business logic stays outside the Lambda adapter. Developers should only need to touch the handler contract and usecase layer for most projects.

## Proposed Work Plan

### Phase 1: Reposition the Repository Around Lambda

1. Rename all template language from general `go-template` wording to Lambda-specific wording.
2. Update the repository identity in:
   - `README.md`
   - `TEMPLATE_GUIDE.md`
   - `WORKFLOW.md`
   - `go.mod`
   - `Makefile`
   - Terraform tags, resource names, and outputs
3. Rewrite the README quick start so it describes creating a Lambda function from this template, not adding Cobra commands.

Success outcome:

- a new user immediately understands this is an AWS Lambda template and how to start using it

### Phase 2: Remove Generic Template Baggage

Remove code and assets that increase cognitive load for Lambda users without helping them ship functions.

Planned removals:

- `internal/cli/hello.go`
- `internal/cli/root.go`
- Cobra dependency from `go.mod`
- CLI-focused tests under `test/unit/cli`
- Kubernetes manifests under `deploy/kubernetes/`
- container-oriented `docker-compose.yml`
- `pkg/health` unless there is a clear Lambda-specific use for it
- release steps that publish generic multi-platform binaries and Docker images

Files to reassess rather than blindly delete:

- `Dockerfile`: either remove it or repurpose it strictly for Lambda packaging
- `scripts/build.sh`: replace multi-platform binary build logic with Lambda packaging logic

Success outcome:

- the template surface area is smaller and clearly Lambda-specific

### Phase 3: Introduce a Lambda-First Runtime Skeleton

Add the runtime structure needed for real Lambda usage.

Planned changes:

1. Add `github.com/aws/aws-lambda-go` as a core dependency.
2. Replace `cmd/app/main.go` with a Lambda entrypoint, or move to `cmd/lambda/main.go`.
3. Add a dedicated Lambda adapter package, preferably `internal/handler` or `internal/adapter/lambda`.
4. Add an application wiring package to centralize dependency construction.
5. Standardize one default execution path:

```go
func main() {
    app := bootstrap()
    lambda.Start(app.Handler)
}
```

Recommended handler pattern:

- keep the API Gateway event mapping in the handler layer
- keep validation and orchestration in usecases
- keep pure business rules in domain packages

Success outcome:

- the template boots directly as an AWS Lambda function without any CLI leftovers

### Phase 4: Define the Developer Extension Point

The template needs one obvious place where developers put their own logic.

Recommended developer contract:

1. edit a single handler file for the event contract
2. implement business logic in a usecase
3. keep infrastructure and bootstrapping unchanged unless needed

Suggested initial scaffold:

- `internal/handler/api.go` using API Gateway proxy request and response types
- `internal/usecase/execute.go` with a minimal example business flow
- `internal/adapter/dynamodb/` with a default repository stub
- `internal/adapter/ssm/` with a parameter loader stub
- `test/fixtures/events/apigw-request.json`
- `test/integration/lambda/handler_test.go`

Important design choice:

- API Gateway should be the default trigger model
- do not try to support SQS, S3, EventBridge, and DynamoDB streams in the first scaffold all at once

Reason:

- a template that tries to model every Lambda trigger becomes harder to understand
- the better template is one clean default plus documentation on how to swap event types

Recommended default:

- start with API Gateway proxy integration handled through `lambda.Start`
- document how to replace API Gateway with another event source later if a team needs that shape

Success outcome:

- developers have one obvious customization point instead of multiple competing patterns

### Phase 5: Rebuild the Build and Packaging Workflow for Lambda

The current automation is optimized for generic binaries. That should be replaced with Lambda packaging.

Planned Makefile targets:

- `build` or `build-lambda`: compile Linux Lambda binary
- `package`: create `bootstrap` and zip artifact
- `test`: run unit and integration tests
- `sam-build`: run `sam build`
- `sam-invoke`: run `sam local invoke` against a sample event
- `sam-start-api`: run local API Gateway emulation against the default handler
- `clean`: remove zip and bootstrap artifacts

Recommended build defaults:

- `GOOS=linux`
- `GOARCH=arm64` by default, with `amd64` as optional
- output binary named `bootstrap`
- zip artifact named predictably, for example `dist/lambda.zip`

Likely file changes:

- `Makefile`
- `scripts/build.sh`
- `Dockerfile` if retained
- `deploy/sam/template.yaml`
- sample event files under `test/fixtures/events/`

Success outcome:

- any developer can build a deployable Lambda artifact with one command
- any developer can also run the function locally through SAM without inventing their own local workflow

### Phase 6: Fix Terraform So It Matches the New Template

Terraform already points at Lambda, but it should become the canonical infrastructure path.

Planned Terraform improvements:

1. align names with the new template identity
2. make artifact path and architecture explicit
3. set runtime to current custom runtime guidance such as `provided.al2023` if that is the chosen standard
4. provision API Gateway as the default trigger
5. provision a DynamoDB table as the default persistence resource
6. add log retention configuration through an explicit CloudWatch log group resource
7. add sane defaults for timeout, memory, and environment variables
8. make IAM role setup minimal but production-safe
9. add least-privilege permissions for:
  - CloudWatch Logs write access
  - DynamoDB access scoped to the template table
  - SSM Parameter Store read access scoped by path or prefix
10. document where users extend permissions for their own AWS integrations

Optional but useful additions:

- API Gateway stage configuration with access logging
- example EventBridge trigger module stub kept commented or isolated

Success outcome:

- the repo contains one clean deployment path that matches the generated artifact
- the repo includes the most common runtime dependencies by default instead of leaving them as first-user setup work

### Phase 6A: Add AWS SAM for Local Development

Add a focused SAM setup that improves local development without duplicating all Terraform concerns.

Planned SAM deliverables:

1. add `deploy/sam/template.yaml`
2. define one function resource wired to the default handler
3. define an API Gateway event in the SAM template by default
4. point SAM build metadata at the Go source and generated bootstrap artifact
5. include one or more sample API Gateway events under `test/fixtures/events/`
6. add Make targets for `sam build`, `sam local invoke`, and `sam local start-api`
7. document required local tools such as Docker and the SAM CLI

Important constraint:

- do not model every IAM policy, environment, and production integration twice in both SAM and Terraform

SAM should stay intentionally thin:

- enough for local invoke and local API testing
- enough for quick-start developer onboarding
- optionally enough for a basic demo deploy

Success outcome:

- developers can validate handler behavior locally with realistic event payloads in a standard AWS-native workflow

### Phase 7: Rebuild CI/CD Around Lambda Artifacts

The GitHub Actions workflows are useful, but the release outputs are wrong for the target product.

Planned CI/CD changes:

1. keep lint, test, security, and CodeQL workflows
2. update CI build job to produce a Lambda artifact instead of a generic CLI binary
3. update release workflow to publish:
   - Lambda zip artifact
   - checksums
   - SBOM
4. remove Docker image publishing unless the template explicitly supports Lambda container images
5. keep dependency automation and changelog workflow if they still fit the maintenance model
6. optionally add a CI validation step for `sam validate` if SAM is added

Success outcome:

- the release output is directly usable for Lambda deployment

### Phase 8: Rework the Tests Around Lambda Behavior

The current tests validate CLI metadata and greeting examples. They should validate Lambda behavior instead.

Planned test changes:

- remove CLI-specific tests
- add API Gateway handler unit tests
- add fixture-driven integration tests
- verify error mapping and logging behavior
- add repository tests around DynamoDB adapter behavior boundaries
- add configuration tests for SSM-backed parameter loading behavior
- keep existing domain and usecase tests where the logic still applies

Suggested test layout:

```text
test/
  fixtures/events/
    apigw-request.json
  integration/lambda/
    handler_test.go
  unit/
    adapter/
    config/
    domain/
    errors/
    handler/
    usecase/
```

Success outcome:

- tests prove that the template works as a Lambda project, not as a CLI demo

### Phase 9: Rewrite the Documentation for Template Consumers

Documentation needs a full rewrite, not just light edits.

Planned documentation deliverables:

1. `README.md`
   - what this template is for
   - quick start from GitHub template
  - how to replace the sample API route logic
   - how to build and deploy
  - what is provided by default: API Gateway, DynamoDB, SSM, CloudWatch logs
2. `TEMPLATE_GUIDE.md`
   - exact customization checklist after creating a repo
3. `WORKFLOW.md`
  - how local development, testing, packaging, SAM local invoke, and deployment work
4. optional `LAMBDA_GUIDE.md`
  - API Gateway event model guidance
   - architecture explanation
  - how DynamoDB and SSM are wired in by default
  - how to add additional AWS service integrations safely
5. document tool responsibilities clearly:
  - SAM for local workflow
  - Terraform for infrastructure and production deployment

Required tone for docs:

- opinionated
- short
- task-oriented
- focused on what the user edits versus what the template handles for them

Success outcome:

- a developer can fork the template and start coding their Lambda logic without reading the source tree first

## Proposed File-by-File Impact

Likely to keep and adapt:

- `internal/config/config.go`
- `internal/domain/greeter.go` or equivalent domain examples, if rewritten to be Lambda-relevant
- `internal/errors/errors.go`
- `internal/logger/logger.go`
- `internal/usecase/greeting.go` or equivalent usecase sample, if renamed and simplified
- `pkg/version/version.go`
- `deploy/terraform/*`
- `deploy/sam/*`
- `.github/workflows/ci.yml`
- `.github/workflows/codeql.yml`
- `.github/workflows/dependencies.yml`

Likely to remove:

- `internal/cli/*`
- `deploy/kubernetes/*`
- `docker-compose.yml`
- CLI-specific tests

Likely to add:

- `cmd/lambda/main.go`
- `internal/app/app.go`
- `internal/handler/api.go`
- `internal/adapter/dynamodb/*`
- `internal/adapter/ssm/*`
- `deploy/sam/template.yaml`
- `test/fixtures/events/apigw-request.json`
- `test/integration/lambda/handler_test.go`
- `LAMBDA_GUIDE.md` or equivalent

## Implementation Order

Recommended execution sequence:

1. rewrite README and template positioning
2. remove CLI/Cobra scaffolding
3. add Lambda runtime entrypoint and handler package
4. add SAM template and local invoke workflow
5. update Makefile and build script for Lambda packaging
6. align Terraform with the new artifact shape
7. update CI/release workflows
8. replace CLI tests with Lambda tests
9. finish documentation and cleanup

This order reduces churn because the code structure gets corrected before automation and docs are finalized.

## Acceptance Criteria

The conversion is complete when all of the following are true:

- `go test ./...` passes with Lambda-focused tests
- a single make target produces a deployable Lambda zip
- a SAM command can invoke the function locally using a sample event
- a SAM command can start a local API Gateway-compatible endpoint for development
- Terraform deploys the packaged function without manual file renaming
- Terraform provisions API Gateway, DynamoDB, SSM read access, and CloudWatch log retention by default
- no Cobra CLI code remains in the default template path
- README quick start matches the actual repository behavior
- developers can find one obvious file to edit for their Lambda logic

## Suggested Guardrails

To keep the template clean over time:

- do not add multiple trigger models into the default scaffold
- do not keep Kubernetes or generic server deployment examples in this repo
- do not introduce framework-heavy abstractions unless a real repeated need appears
- keep the handler thin and push logic into usecases
- keep Terraform example minimal and safe by default

## Final Direction

This repo already has a good Go engineering foundation. The problem is not quality, it is positioning and runtime shape.

The right move is to aggressively simplify it into a Lambda-first template:

- remove CLI-first scaffolding
- add a real Lambda runtime entrypoint
- standardize build and deploy around Lambda artifacts
- document one clear customization path for developers

That will turn this repository from a generic Go starter into a template that genuinely helps teams ship AWS Lambda functions quickly.
