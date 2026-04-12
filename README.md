# Go Template

![Go Version](https://img.shields.io/badge/Go-1.23+-blue?style=flat-square&logo=go)
![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)
![CI Status](https://github.com/PlatformStackPulse/go-template/actions/workflows/ci.yml/badge.svg)
[![codecov](https://codecov.io/gh/PlatformStackPulse/go-template/branch/main/graph/badge.svg)](https://codecov.io/gh/PlatformStackPulse/go-template)
![DevContainer](https://img.shields.io/static/v1?label=DevContainer&message=Ready&color=blue&style=flat-square&logo=visual-studio-code)

<p align="center">
  <strong>Slim, Production-Ready Go Template</strong><br>
  Enterprise CI/CD, DevSecOps, clean architecture. Optimized for CLI tools and API servers.
</p>

---

## 🎯 Overview

A **minimal, reusable GitHub template** for building production-ready Go applications. Supports both **CLI single-binary projects** and **API servers** with zero bloat.

**What you get:**
- ✅ Clean architecture (Domain/Usecase/Adapter layers)
- ✅ CLI foundation (Cobra framework, ready for commands)
- ✅ Structured logging (slog)
- ✅ Configuration management
- ✅ Comprehensive testing (unit, integration)
- ✅ DevSecOps (gosec, govulncheck, CodeQL)
- ✅ GitHub Actions CI/CD (linting, testing, releases)
- ✅ Docker & multi-platform builds
- ✅ DevContainer with pre-configured tools
- ✅ Kubernetes & Terraform examples

**What you don't get (keep it slim!):**
- No bloated feature flag systems
- No unused HTTP/health endpoints (add only if needed)
- No example commands cluttering production code (we include one example to remove)
- No over-engineered abstractions

---

## 🚀 Quick Start

### Using as GitHub Template

```bash
# Create a new repo from this template
gh repo create my-app --template go-template

# Setup
cd my-app
make dev-setup

# Build & run
make build
./bin/go-template hello --name "World"
```

### Example: Add Your First Command

The template includes an example command in `internal/cli/hello.go` (via `NewExampleCommand`). Replace it with your own:

**1. Create your command:**
```go
// internal/cli/mycommand.go
package cli

import (
	"github.com/spf13/cobra"
	"github.com/PlatformStackPulse/go-template/internal/logger"
)

func NewMyCommand(log *logger.Logger) *cobra.Command {
	return &cobra.Command{
		Use:        "mycommand",
		Short:      "What my command does",
		Long:       "Detailed description of what my command does",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info("Running mycommand")
			// Your logic here
			return nil
		},
	}
}
```

**2. Register in main.go:**
```go
// cmd/app/main.go
cmd.AddCommand(cli.NewMyCommand(log))
```

**3. Remove the example command:**
```bash
rm internal/cli/hello.go
# Update cmd/app/main.go - remove the cli.NewExampleCommand line
```

---

## 📁 Lean Project Structure

```
go-template/
├── cmd/app/
│   └── main.go                 # Entry point (minimal, ~30 lines)
├── internal/                   # Private packages
│   ├── cli/                    # Cobra commands
│   │   ├── root.go            # Root command
│   │   └── hello.go           # Example command (remove/rename)
│   ├── config/                # Configuration loading
│   ├── domain/                # Domain entities
│   ├── logger/                # Structured logging (slog)
│   ├── usecase/               # Business logic
│   └── adapter/               # External integrations (add as needed)
├── pkg/
│   └── version/               # Version info (injected at build)
├── test/
│   ├── unit/                  # Unit tests
│   └── integration/           # Integration tests (if needed)
├── deploy/
│   ├── kubernetes/            # K8s manifests (optional)
│   └── terraform/             # IaC (optional)
├── Makefile                   # Core build targets
├── Dockerfile                 # Multi-stage build
├── go.mod / go.sum            # Dependencies
└── .github/workflows/         # GitHub Actions (6 workflows)
```

---

## 🎯 Two Modes: CLI vs API

### Mode 1: CLI (Single Binary)

Your main.go stays **slim**:
```go
func main() {
	cfg := config.Load()
	log := logger.NewLogger(cfg.Debug)
	cmd := cli.NewRootCommand(log)
	cmd.Version = version.Version
	
	if err := cmd.ExecuteContext(context.Background()); err != nil {
		os.Exit(1)
	}
}
```

**Add your commands** to `internal/cli/` and register them in `cmd/app/main.go`.

### Mode 2: API Server

Extend main.go with HTTP:
```go
func main() {
	cfg := config.Load()
	log := logger.NewLogger(cfg.Debug)
	
	// Create HTTP server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: setupRoutes(log),
	}
	
	// Run with graceful shutdown
	// (add your HTTP handler logic)
}
```

**Optional**: Add files as needed:
- `internal/http/server.go` — HTTP server setup
- `internal/http/handlers.go` — Route handlers  
- `pkg/health/` — Health check endpoints

---

## 📐 Architecture

**Clean layers** with clear separation:

```
CLI/HTTP Layer (User interaction)
         ↓
Usecases (Business logic)
         ↓
Domain (Entities, rules)
         ↓
Adapters (External services)
```

Each layer is **independent and testable**.

---

## 🪵 Design Philosophy

- **Slim** — Only essential structure, no bloat
- **Extensible** — Easy to add features without refactoring
- **Example-first** — Rename/remove example command, add yours
- **Test-friendly** — Proper layering makes testing straightforward
- **Single responsibility** — Each package does ONE thing well

---

## ⚙️ Common Tasks

```bash
make dev-setup       # Install tools
make build           # Build binary
make test            # Run tests
make coverage        # Coverage report
make fmt             # Format code
make lint            # Lint check
make security        # Security scan
make docker-build    # Build Docker image
```

See [Makefile](Makefile) for all targets.

---

## 🔄 Customization Checklist

When using this template:

- [ ] Update `go.mod` module name
- [ ] Rename `go-template` binary in `Makefile` and `README.md`
- [ ] Remove/rename `internal/cli/hello.go` example command
- [ ] Update repository references in documentation
- [ ] Add your commands to `internal/cli/`
- [ ] Update tests in `test/unit/`

---

## 📦 Deployment

### Build for Production

```bash
make build                # Build locally
make release             # Build all platforms (macOS/Linux/Windows)
make docker-build        # Build Docker image
```

### Docker

```dockerfile
# Already configured in Dockerfile (multi-stage build)
docker build -t my-app .
docker run my-app hello --name "Docker"
```

### Kubernetes (Optional Example)

```bash
kubectl apply -f deploy/kubernetes/
```

See [deploy/kubernetes/](deploy/kubernetes/) for manifests.

### Terraform (Optional Example)

```bash
cd deploy/terraform
terraform init
terraform plan -var-file=terraform.dev.tfvars
```

---

## 🧪 Testing

```bash
# All tests
make test

# With coverage
make coverage

# Specific package
go test ./internal/usecase/...

# Race detection
go test -race ./...
```

---

## 🔐 Security

- Integrated security scanning (gosec, govulncheck, CodeQL)
- Vulnerability notifications via Dependabot
- Branch protection on main
- Commit signing support

See [SECURITY.md](SECURITY.md) for details.

---

## 📚 Documentation

- [CONTRIBUTING.md](CONTRIBUTING.md) — How to contribute
- [SECURITY.md](SECURITY.md) — Security policy
- [WORKFLOW.md](WORKFLOW.md) — Git workflow & branch protection
- [Makefile](Makefile) — Build targets (run `make help`)

---

## 📄 License

MIT License — See [LICENSE](LICENSE)

---

## 🎓 Learn More

- [Cobra Framework](https://cobra.dev)
- [Go Best Practices](https://golang.org/doc/effective_go)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

---

**Ready to build?** Fork this template and start coding! 🚀

### Design Patterns Included

- **Command Pattern** — Cobra commands
- **Mediator Pattern** — Feature manager
- **Orchestrator Pattern** — Usecase layer
- **Dependency Injection** — Idiomatic Go

---

## 🛠 Available Commands

### Makefile Targets

```bash
make help              # Show all available targets
make build             # Build the application
make run               # Build and run
make test              # Run all tests with coverage
make test-unit         # Run unit tests only
make test-integration  # Run integration tests only
make coverage          # Generate coverage report
make clean             # Remove build artifacts
make install           # Install dependencies
make lint              # Run linters
make fmt               # Format code
make vet               # Run go vet
make security          # Run security checks
make sec-update        # Check for security updates
make dev-setup         # Setup development environment
make changelog         # Regenerate CHANGELOG.md from commits
make changelog-check   # Verify CHANGELOG.md is current
make watch             # Watch for changes (air)
make version           # Show version info
make all               # Run all tasks
```

### CLI Commands

```bash
./bin/go-template --help
./bin/go-template hello                    # Print "Hello, World!"
./bin/go-template hello --name Alice       # Print "Hello, Alice!"
./bin/go-template version                  # Show version information
```

---

## 🧪 Testing

### Running Tests

```bash
# All tests with coverage
make test

# Specific test suite
make test-unit
make test-integration

# Generate HTML coverage report
make coverage

# Run with race detector
go test -race ./...
```

### Test Coverage

- **Minimum threshold:** 70%
- **Enforced in CI:** Yes
- **Table-driven tests:** Domain and usecase layers
- **Mock patterns:** Via interfaces

---

## 🔐 Security & Quality

### Integrated Security Tools

| Tool | Purpose | Status |
|------|---------|--------|
| `golangci-lint` | Code quality & style | ✅ CI/CD |
| `gosec` | Security issues | ✅ CI/CD |
| `govulncheck` | Dependency vulnerabilities | ✅ CI/CD |
| `CodeQL` | Static analysis | ✅ CI/CD |
| `Dependabot` | Dependency updates | ✅ Automated |

### Security Checks

```bash
# Run all security scans
make security

# Check for vulnerabilities
make sec-update

# Manual gosec scan
gosec ./...

# Manual govulncheck
govulncheck ./...
```

---

## 📦 Versioning & Releases

### Semantic Versioning

This project follows [Semantic Versioning](https://semver.org/):

```
MAJOR.MINOR.PATCH
  ↑      ↑      ↑
  │      │      └─ Bug fixes
  │      └────────── New features (backward compatible)
  └───────────────── Breaking changes
```

### Conventional Commits

Commit messages must follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <description>

<body>

<footer>
```

**Types:**
- `feat` — New feature
- `fix` — Bug fix
- `docs` — Documentation
- `style` — Code style
- `refactor` — Code refactoring
- `perf` — Performance improvement
- `test` — Test additions/updates
- `chore` — Build, dependencies, etc.

**Examples:**
```
feat: add hello command with name parameter
fix: resolve panic in logger initialization
docs: update README with quick start guide
chore: upgrade Go from 1.21 to 1.22
```

### Automated Release Process

1. **Create and push tag:**
   ```bash
   git tag v1.2.3
   git push origin v1.2.3
   ```

2. **Release workflow triggers:**
   - Multi-platform builds
   - GitHub Release creation
   - Docker image push
   - SBOM generation

3. **Artifacts available at:** `https://github.com/PlatformStackPulse/my-app/releases/tag/v1.2.3`

### Changelog Strategy

- `CHANGELOG.md` is generated from Conventional Commits using `git-chglog`.
- On merges/pushes to `main`, `.github/workflows/changelog.yml` updates and commits `CHANGELOG.md` automatically.
- Maintainers can regenerate locally with:

```bash
make changelog
```

---

## 🌟 Feature Flags

### Quick Example

```go
fm := feature.NewManager()

if fm.IsEnabled(feature.FeatureHello) {
    fmt.Println("Feature is enabled!")
}
```

### Define Flags

Edit `internal/feature/flags.go`:

```go
const (
    FeatureNewAPI    Flag = "feature_new_api"
    FeatureMetrics   Flag = "feature_metrics"
)
```

### Enable Flags

```bash
# Environment variable format
export FEATURE_FLAGS="feature_new_api=true,feature_metrics=false"
./bin/go-template hello
```

### Future Extensions

- Config file support (YAML/JSON)
- Remote flag service integration
- Gradual rollout (percentage)
- User-based targeting
- A/B testing support

---

## 🐳 Docker

### Build Image

```bash
# Build locally
docker build -t go-template:latest .

# Run container
docker run --rm go-template:latest hello

# Run with feature flags
docker run --rm \
  -e FEATURE_FLAGS="feature_hello=true" \
  go-template:latest hello
```

### Docker Compose

```bash
# Development
docker-compose up dev

# Run tests
docker-compose run test

# Run linter
docker-compose run lint
```

---

## 🚀 CI/CD Pipelines

### GitHub Actions Workflows

#### 1. **CI Pipeline** (`.github/workflows/ci.yml`)
Runs on every push and PR:
- Code linting (golangci-lint)
- Format verification
- Unit & integration tests
- Coverage enforcement (≥70%)
- Security scans (gosec, CodeQL)
- PR commit message validation
- Multi-version Go testing (1.21, 1.22)

#### 2. **Release Pipeline** (`.github/workflows/release.yml`)
Triggers on version tags (`v*.*.*`):
- Multi-platform builds (Linux, macOS, Windows)
- Checksum generation
- GitHub Release creation with artifacts
- Docker image build and push
- SBOM generation

#### 3. **Dependency Management** (`.github/workflows/dependencies.yml`)
Weekly automated tasks:
- Dependency updates
- Security vulnerability scans
- Automated PR creation

#### 4. **CodeQL Analysis** (`.github/workflows/codeql.yml`)
Runs code scanning on pushes/PRs to `main` and on a weekly schedule.

#### 5. **Changelog Update** (`.github/workflows/changelog.yml`)
Updates `CHANGELOG.md` from Conventional Commits on pushes to `main`.

---

## 📚 Infrastructure Templates

### Kubernetes

Deploy to Kubernetes:

```bash
kubectl apply -f deploy/kubernetes/

# Check deployment
kubectl get pods -l app=go-template
```

### Terraform

Deploy to AWS Lambda:

```bash
cd deploy/terraform

terraform init
terraform plan -var-file=terraform.dev.tfvars
terraform apply -var-file=terraform.dev.tfvars
```

---

## 📖 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DEBUG` | Enable debug logging | `false` |
| `APP_NAME` | Application name | `go-template` |
| `APP_VERSION` | Application version | `dev` |
| `FEATURE_FLAGS` | Enabled feature flags | "" |

### Example `.env`

```env
DEBUG=true
APP_NAME=my-app
APP_VERSION=v1.0.0
FEATURE_FLAGS=feature_hello=true,feature_metrics=true
```

---

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/awesome-feature`
3. Commit with conventional format: `git commit -m "feat: add awesome feature"`
4. Push to branch: `git push origin feature/awesome-feature`
5. Open a Pull Request

### PR Requirements

- ✅ Must follow Conventional Commits
- ✅ Tests passing (coverage ≥ 70%)
- ✅ Code linting passing
- ✅ No security warnings
- ✅ Documentation updated

---

## 📋 Checklist for New Projects

When using this template:

- [ ] Update `go.mod` module path
- [ ] Update `module path in code imports`
- [ ] Create `.github/CODEOWNERS` file
- [ ] Update README with project-specific info
- [ ] Configure branch protection rules
- [ ] Setup Dependabot alerts
- [ ] Update repository description
- [ ] Add repository topics
- [ ] Customize GitHub Actions (if needed)
- [ ] Update Terraform variables for your environment

---

## 📜 License

This project is licensed under the MIT License — see [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

Built with best practices from:
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [Semantic Versioning](https://semver.org/)
- [Go Best Practices](https://golang.org/doc/effective_go)
- [GoReleaser](https://goreleaser.com/)
- [Cobra](https://cobra.dev/)

---

## 📞 Support

- 📖 [Documentation](https://github.com/PlatformStackPulse/go-template/wiki)
- 🐛 [Report Issues](https://github.com/PlatformStackPulse/go-template/issues)
- 💬 [Discussions](https://github.com/PlatformStackPulse/go-template/discussions)

---

<p align="center">
  <sub>Built with ❤️ for platform engineers and Go developers</sub>
</p>
