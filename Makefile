.PHONY: help install build package docker-build terraform-bootstrap ecr-push test test-unit test-integration coverage clean lint fmt vet security dev-setup changelog changelog-check sam-build sam-invoke sam-start-api terraform-plan terraform-apply terraform-plan-zip terraform-apply-zip

APP_NAME?=go-lambda-template
VERSION?=dev
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GO_VERSION?=$(shell go version | awk '{print $$3}')
ARCH?=arm64
DOCKER_PLATFORM?=linux/arm64
IMAGE_TAG?=latest
ARTIFACT_DIR?=dist
BOOTSTRAP=$(ARTIFACT_DIR)/bootstrap
ZIP_ARTIFACT=$(ARTIFACT_DIR)/lambda.zip
TERRAFORM_DIR=deploy/terraform

LD_FLAGS=-ldflags "-X github.com/PlatformStackPulse/go-lambda-template/pkg/version.Version=$(VERSION) \
	-X github.com/PlatformStackPulse/go-lambda-template/pkg/version.Commit=$(COMMIT) \
	-X github.com/PlatformStackPulse/go-lambda-template/pkg/version.BuildTime=$(BUILD_TIME) \
	-X github.com/PlatformStackPulse/go-lambda-template/pkg/version.GoVersion=$(GO_VERSION)"

help: ## Show available targets
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

install: ## Download and tidy dependencies
	@go mod download
	@go mod tidy

build: ## Build the Lambda bootstrap binary
	@mkdir -p $(ARTIFACT_DIR)
	@CGO_ENABLED=0 GOOS=linux GOARCH=$(ARCH) go build $(LD_FLAGS) -o $(BOOTSTRAP) ./cmd/lambda
	@echo "Built $(BOOTSTRAP)"

package: build ## Package the Lambda bootstrap binary as a zip artifact
	@rm -f $(ZIP_ARTIFACT)
	@cd $(ARTIFACT_DIR) && zip -q lambda.zip bootstrap
	@echo "Packaged $(ZIP_ARTIFACT)"

docker-build: ## Build the Lambda container image for ECR deployment
	@docker build --platform $(DOCKER_PLATFORM) \
		--build-arg TARGETARCH=$(ARCH) \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		--build-arg GO_VERSION=$(GO_VERSION) \
		-t $(APP_NAME):$(IMAGE_TAG) .

terraform-bootstrap: ## Create the ECR repository before the first image push
	@terraform -chdir=$(TERRAFORM_DIR) init
	@terraform -chdir=$(TERRAFORM_DIR) apply -target=aws_ecr_repository.lambda -target=aws_ecr_lifecycle_policy.lambda -var-file=terraform.dev.tfvars -var="deployment_package_type=image" -var="image_tag=$(IMAGE_TAG)"

ecr-push: docker-build terraform-bootstrap ## Push the Lambda container image to ECR
	@chmod +x scripts/push-ecr.sh
	@scripts/push-ecr.sh $(IMAGE_TAG) $(TERRAFORM_DIR) $(APP_NAME):$(IMAGE_TAG)

test: ## Run all tests with coverage
	@go test -v -race -covermode=atomic -coverpkg=./internal/...,./pkg/... -coverprofile=coverage.txt ./...
	@go tool cover -func=coverage.txt | tail -1

test-unit: ## Run unit tests only
	@go test -v -race ./test/unit/...

test-integration: ## Run integration tests only
	@go test -v ./test/integration/...

coverage: test ## Generate an HTML coverage report
	@go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report written to coverage.html"

clean: ## Remove build and coverage artifacts
	@rm -rf $(ARTIFACT_DIR) .aws-sam coverage.txt coverage.html
	@go clean

lint: ## Run golangci-lint
	@golangci-lint run ./...

fmt: ## Format Go code
	@go fmt ./...
	@gofmt -w .

vet: ## Run go vet
	@go vet ./...

security: ## Run gosec and govulncheck
	@gosec ./...
	@govulncheck ./...

dev-setup: ## Install local developer tooling
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@go install github.com/git-chglog/git-chglog/cmd/git-chglog@latest
	@echo "Install the AWS SAM CLI separately: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html"

changelog: ## Regenerate CHANGELOG.md from Conventional Commits
	@chmod +x scripts/update-changelog.sh
	@scripts/update-changelog.sh

changelog-check: ## Verify CHANGELOG.md is up to date
	@cp CHANGELOG.md CHANGELOG.md.bak
	@chmod +x scripts/update-changelog.sh
	@scripts/update-changelog.sh
	@cmp -s CHANGELOG.md CHANGELOG.md.bak || (echo "CHANGELOG.md is outdated. Run 'make changelog'." && rm -f CHANGELOG.md.bak && exit 1)
	@rm -f CHANGELOG.md.bak

sam-build: build ## Build the Lambda artifact using AWS SAM
	@sam build --template-file deploy/sam/template.yaml --build-dir .aws-sam/build

sam-invoke: sam-build ## Invoke the function locally with a sample API Gateway event
	@sam local invoke ApiFunction --template-file deploy/sam/template.yaml --event test/fixtures/events/apigw-request.json

sam-start-api: sam-build ## Start the local API Gateway emulator
	@sam local start-api --template-file deploy/sam/template.yaml

terraform-plan: ## Run terraform plan for the default image-based deployment
	@terraform -chdir=$(TERRAFORM_DIR) init
	@terraform -chdir=$(TERRAFORM_DIR) plan -var-file=terraform.dev.tfvars -var="deployment_package_type=image" -var="image_tag=$(IMAGE_TAG)"

terraform-apply: ecr-push ## Build, push, and deploy the default image-based Lambda stack
	@terraform -chdir=$(TERRAFORM_DIR) init
	@terraform -chdir=$(TERRAFORM_DIR) apply -var-file=terraform.dev.tfvars -var="deployment_package_type=image" -var="image_tag=$(IMAGE_TAG)"

terraform-plan-zip: package ## Run terraform plan for the optional zip-based deployment
	@terraform -chdir=$(TERRAFORM_DIR) init
	@terraform -chdir=$(TERRAFORM_DIR) plan -var-file=terraform.dev.tfvars -var="deployment_package_type=zip"

terraform-apply-zip: package ## Apply terraform for the optional zip-based Lambda deployment
	@terraform -chdir=$(TERRAFORM_DIR) init
	@terraform -chdir=$(TERRAFORM_DIR) apply -var-file=terraform.dev.tfvars -var="deployment_package_type=zip"
