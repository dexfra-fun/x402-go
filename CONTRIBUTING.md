# Contributing to x402-go

Thank you for your interest in contributing to x402-go! We welcome contributions from the community.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR-USERNAME/x402-go.git
   cd x402-go
   ```
3. **Create a feature branch**:
   ```bash
   git checkout -b feature/my-new-feature
   ```

## Development Setup

### Prerequisites

- Go 1.22 or higher
- Git

### Install Dependencies

```bash
go mod download
```

### Run Tests

```bash
go test -v ./...
```

### Run Linter

```bash
golangci-lint run
```

## Making Changes

1. **Write tests** for your changes
2. **Ensure all tests pass**: `go test -v ./...`
3. **Follow Go conventions**: `go fmt ./...` and `go vet ./...`
4. **Update documentation** if needed
5. **Commit your changes** with clear commit messages:
   ```bash
   git commit -m "Add feature: description of feature"
   ```

## Commit Message Guidelines

- Use present tense ("Add feature" not "Added feature")
- Use imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit first line to 72 characters
- Reference issues and pull requests after first line

Examples:
```
Add support for EVM networks

- Implement EVM wallet adapter
- Add network configuration for Ethereum
- Update documentation

Fixes #123
```

## Pull Request Process

1. **Update documentation** (README, code comments, etc.)
2. **Ensure CI passes** (tests, linting)
3. **Request review** from maintainers
4. **Address feedback** from reviewers
5. **Squash commits** if requested
6. **Maintainer will merge** when approved

## Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Add comments for exported functions, types, and packages
- Keep functions small and focused
- Write meaningful variable names

## Testing

- Write unit tests for all new code
- Aim for >80% code coverage
- Use table-driven tests where appropriate
- Mock external dependencies

Example test:
```go
func TestPricingFixed(t *testing.T) {
    tests := []struct {
        name     string
        amount   float64
        expected float64
    }{
        {"zero", 0, 0},
        {"small", 0.001, 0.001},
        {"large", 10.5, 10.5},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            p := pricing.NewFixed(tt.amount)
            got, err := p.GetPrice(context.Background(), x402.Resource{})
            if err != nil {
                t.Fatalf("unexpected error: %v", err)
            }
            if got != tt.expected {
                t.Errorf("expected %f, got %f", tt.expected, got)
            }
        })
    }
}
```

## Reporting Bugs

1. **Check existing issues** to avoid duplicates
2. **Use the bug report template**
3. **Provide detailed information**:
   - Go version
   - Operating system
   - Steps to reproduce
   - Expected vs actual behavior
   - Relevant logs/screenshots

## Requesting Features

1. **Check existing issues** to avoid duplicates
2. **Use the feature request template**
3. **Describe the use case** and why it's valuable
4. **Provide examples** if possible

## Questions?

- Open a [GitHub Discussion](https://github.com/dexfra-fun/x402-go/discussions)
- Join our [Discord community](https://discord.gg/dexfra)

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
