# Contributing to Go Template

Thank you for your interest in contributing! This project follows a set of guidelines to ensure code quality and consistency.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork:**
   ```bash
   git clone https://github.com/PlatformStackPulse/go-template.git
   cd go-template
   ```

3. **Setup development environment:**
   ```bash
   make dev-setup
   ```

4. **Create a feature branch:**
   ```bash
   git checkout -b feature/my-awesome-feature
   ```

## Development Workflow

### Before You Start

- Review existing issues and PRs to avoid duplicates
- Open an issue first for significant changes
- Discuss your approach with maintainers

### Making Changes

1. **Ensure tests pass:**
   ```bash
   make test
   ```

2. **Follow code style:**
   ```bash
   make fmt lint
   ```

3. **Run security checks:**
   ```bash
   make security
   ```

4. **Commit with conventional format:**
   ```bash
   git commit -m "feat: add new feature"
   git commit -m "fix: resolve issue"
   ```

### Conventional Commits

All commits must follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

**Format:**
```
<type>(<scope>): <description>

<optional body>

<optional footer>
```

**Types:**
- `feat` — A new feature
- `fix` — A bug fix
- `docs` — Documentation only changes
- `style` — Changes that don't affect code meaning (formatting, etc.)
- `refactor` — Code change that neither fixes bugs nor adds features
- `perf` — Code change that improves performance
- `test` — Adding missing tests or correcting existing tests
- `chore` — Changes to build process, dependencies, etc.
- `ci` — Changes to CI configuration
- `build` — Changes to build system

**Examples:**
```
feat: add support for feature flags
feat(cli): add hello command
fix: resolve nil pointer exception
fix(logger): fix timestamp formatting
docs: update README with examples
chore: upgrade Go to 1.22
```

## Testing

### Write Tests

- Add tests for new features
- Update tests for bug fixes
- Follow table-driven test pattern
- Use `testify` for assertions

```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "test case 1",
            input:    "input",
            expected: "output",
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            result := MySomething(tc.input)
            assert.Equal(t, tc.expected, result)
        })
    }
}
```

### Run Tests Before Submitting

```bash
make test           # Run all tests
make test-unit      # Unit tests only
make test-integration # Integration tests only
make coverage       # Generate coverage report
```

## Code Style

### Go Code

- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use `gofmt` for formatting
- Run `golangci-lint` regularly

```bash
make fmt
make lint
```

### File Organization

The template follows Clean Architecture principles:

```
internal/
├── domain/         # Pure business entities (no external deps)
│   └── greeter.go  # Example: Domain entity with business logic
├── usecase/        # Business logic orchestration
│   └── greeting.go # Example: Orchestrates domain + adapters
├── adapter/        # External integrations (database, HTTP, etc)
│   └── (empty - add as needed)
├── cli/            # CLI interface layer
│   ├── root.go     # Root command
│   └── hello.go    # Example command
├── config/         # Configuration loading
├── logger/         # Structured logging
└── errors/         # Error types (if needed)

test/
├── unit/           # Unit tests (mirrors internal/)
│   ├── domain/
│   ├── usecase/
│   ├── cli/
│   └── config/
└── integration/    # Integration/end-to-end tests
```

### Testing in Each Layer

**Domain Layer (Pure Business Logic)**
```go
// test/unit/domain/greeter_test.go
func TestGreeterGreet(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {name: "greet with name", input: "Alice", expected: "Hello, Alice!"},
        {name: "greet without name", input: "", expected: "Hello, World!"},
    }
    
    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            greeter := domain.NewGreeter(tc.input)
            assert.Equal(t, tc.expected, greeter.Greet())
        })
    }
}
```

**Usecase Layer (Orchestration)**
- Mock adapters/repositories
- Verify business logic flow
- Test error handling

**CLI Layer**
- Mock logger and usecase
- Test command parsing and flags
- Verify output format

## Documentation

- Update README if adding features
- Add comments to exported functions
- Update CHANGELOG for significant changes
- Use clear, concise language

## Submitting a Pull Request

1. **Push your changes:**
   ```bash
   git push origin feature/my-awesome-feature
   ```

2. **Create a Pull Request on GitHub**

3. **PR Title:** Must follow Conventional Commits
   ```
   feat: add new feature
   fix: resolve issue
   ```

4. **PR Description:**
   - Clear description of changes
   - Reference related issues (#123)
   - Explain motivation and impact

5. **Wait for review:**
   - Ensure all checks pass
   - Respond to feedback
   - Make requested changes

### PR Checklist

- [ ] PR title follows Conventional Commits
- [ ] All tests pass (`make test`)
- [ ] Coverage maintained/improved
- [ ] Linting passes (`make lint`)
- [ ] Security checks pass (`make security`)
- [ ] Documentation updated
- [ ] No breaking changes (or documented if breaking)
- [ ] Commits follow Conventional Commits

## Adding Features

### Adding a New CLI Command

1. Create command file in `internal/cli/`:
```go
// internal/cli/mycommand.go
package cli

import (
    "github.com/spf13/cobra"
    "github.com/PlatformStackPulse/go-template/internal/logger"
)

func NewMyCommand(log *logger.Logger) *cobra.Command {
    return &cobra.Command{
        Use:   "mycommand",
        Short: "Description",
        RunE: func(cmd *cobra.Command, args []string) error {
            log.Info("Command started")
            // Your logic here
            return nil
        },
    }
}
```

2. Register in `cmd/app/main.go`:
```go
cmd.AddCommand(cli.NewMyCommand(log))
```

3. Add tests in `test/unit/cli/`:
```go
func TestMyCommand(t *testing.T) {
    // Test command execution
}
```

### Adding Domain Logic

1. Define entity in `internal/domain/`
2. Create usecase in `internal/usecase/`
3. Write unit tests for both
4. Update CLI to use the usecase

### Coverage Requirements

- **Minimum**: 70% coverage required
- **Target**: 80%+ coverage for new code
- Check coverage locally: `make coverage`
- View report: `open coverage.html`

## Code Review Process

1. At least one approval required
2. All checks must pass
3. CI/CD pipeline must succeed
4. Maintainer will merge when ready

## Questions or Need Help?

- 📖 Check [documentation](./README.md)
- 🐛 [Open an issue](https://github.com/PlatformStackPulse/go-template/issues)
- 💬 [Start a discussion](https://github.com/PlatformStackPulse/go-template/discussions)

## License

By contributing to this project, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing! 🙏
