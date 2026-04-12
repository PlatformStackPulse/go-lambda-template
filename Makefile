.PHONY: help build run test clean lint fmt vet security coverage install dev-setup changelog changelog-check

# Variables
BINARY_NAME=go-template
VERSION?=dev
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION?=$(shell go version | awk '{print $$3}')

LD_FLAGS=-ldflags "-X github.com/PlatformStackPulse/go-template/pkg/version.Version=$(VERSION) \
                   -X github.com/PlatformStackPulse/go-template/pkg/version.Commit=$(COMMIT) \
                   -X github.com/PlatformStackPulse/go-template/pkg/version.BuildTime=$(BUILD_TIME) \
                   -X github.com/PlatformStackPulse/go-template/pkg/version.GoVersion=$(GO_VERSION)"

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@go build $(LD_FLAGS) -o bin/$(BINARY_NAME) cmd/app/main.go
	@echo "Build complete: bin/$(BINARY_NAME)"

run: build ## Build and run the application
	@bin/$(BINARY_NAME) hello

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race -covermode=atomic -coverpkg=./internal/...,./pkg/... -coverprofile=coverage.txt ./test/...
	@go tool cover -func=coverage.txt | tail -1

test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	@go test -v -race ./test/unit/...

test-integration: ## Run integration tests only
	@echo "Running integration tests..."
	@go test -v -race ./test/integration/...

coverage: test ## Run tests with coverage report
	@go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/ dist/ coverage.txt coverage.html
	@go clean
	@echo "Clean complete"

install: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed"

lint: ## Run linters
	@echo "Running linters..."
	@golangci-lint run ./... || true
	@echo "Lint complete"

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@gofmt -w .
	@echo "Format complete"

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...
	@echo "Vet complete"

security: ## Run security checks
	@echo "Running security checks..."
	@gosec ./... || true
	@echo "Security checks complete"

sec-update: ## Check for security updates
	@echo "Checking for security updates..."
	@go list -u -m all
	@govulncheck ./...

dev-setup: ## Setup development environment
	@echo "Setting up development environment..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@go install github.com/git-chglog/git-chglog/cmd/git-chglog@latest
	@go install github.com/cosmtrek/air@latest
	@echo "Development environment ready"

changelog: ## Regenerate CHANGELOG.md from Conventional Commits
	@chmod +x scripts/update-changelog.sh
	@scripts/update-changelog.sh

changelog-check: ## Verify CHANGELOG.md is up to date
	@cp CHANGELOG.md CHANGELOG.md.bak
	@chmod +x scripts/update-changelog.sh
	@scripts/update-changelog.sh
	@cmp -s CHANGELOG.md CHANGELOG.md.bak || (echo "CHANGELOG.md is outdated. Run 'make changelog'." && rm -f CHANGELOG.md.bak && exit 1)
	@rm -f CHANGELOG.md.bak
	@echo "CHANGELOG.md is up to date"

watch: ## Watch for changes and rebuild (requires air)
	@air

version: ## Show version information
	@bin/$(BINARY_NAME) version || echo "Binary not built. Run 'make build' first."

all: clean install lint test build ## Run all targets
	@echo "All tasks completed"
