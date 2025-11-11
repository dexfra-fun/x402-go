# Linting Guide for x402-go

This project uses [golangci-lint](https://golangci-lint.run/) v2 for code quality and consistency.

## Installation

```bash
# macOS/Linux
brew install golangci-lint

# Or using go install
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

## Running Linters

```bash
# Run all linters
golangci-lint run

# Run with auto-fix
golangci-lint run --fix

# Run on specific packages
golangci-lint run ./pkg/...

# Run only specific linters
golangci-lint run --disable-all --enable=errcheck,govet
```

## Configuration

The project uses `.golangci.yml` configuration with **50+ linters** enabled, including:

### Core Linters
- **errcheck**: Checks for unchecked errors
- **govet**: Reports suspicious constructs
- **staticcheck**: Advanced static analysis (includes gosimple)
- **unused**: Finds unused code

### Security
- **gosec**: Security vulnerability scanner
- **errname**: Error naming conventions
- **errorlint**: Error wrapping checks

### Code Quality
- **gocognit**: Cognitive complexity
- **cyclop**: Cyclomatic complexity
- **funlen**: Function length
- **gocyclo**: Cyclomatic complexity
- **dupl**: Code duplication

### Style & Best Practices
- **revive**: Golint replacement
- **gocritic**: Meta-linter with many checks
- **godot**: Comment punctuation
- **whitespace**: Whitespace checks

### Performance
- **perfsprint**: Sprint optimization
- **makezero**: Slice initialization
- **prealloc**: Preallocated slices (optional)

## Excluding Issues

### Per-File Exclusions

Test files are automatically excluded from some strict checks:
- `funlen`, `gocognit`, `dupl` - Tests can be longer
- `errcheck`, `gosec` - Less strict in tests

### Inline Suppressions

Use `//nolint` directives when necessary:

```go
//nolint:errcheck // Reason: writing to stdout never fails
fmt.Println("hello")

//nolint:funlen,gocognit // Complex test case
func TestComplexScenario(t *testing.T) {
    // ...
}
```

**Note**: Our config requires explanations and specific linter names for nolint directives.

## Common Issues & Fixes

### Unchecked Errors

❌ **Bad**:
```go
defer resp.Body.Close()
```

✅ **Good**:
```go
defer func() {
    _ = resp.Body.Close()
}()
```

### Magic Numbers

❌ **Bad**:
```go
time.Sleep(5 * time.Second)
```

✅ **Good**:
```go
const defaultTimeout = 5 * time.Second
time.Sleep(defaultTimeout)
```

### Comment Formatting

❌ **Bad**:
```go
// NewClient creates a client
func NewClient() {}
```

✅ **Good**:
```go
// NewClient creates a new HTTP client.
func NewClient() {}
```

### Naked Returns

❌ **Bad**:
```go
func process(x int) (result int, err error) {
    result = x * 2
    return  // naked return
}
```

✅ **Good**:
```go
func process(x int) (int, error) {
    result := x * 2
    return result, nil
}
```

## CI/CD Integration

Linting runs automatically on:
- Pull requests
- Push to main/develop branches
- Pre-release checks

See [`.github/workflows/test.yml`](../.github/workflows/test.yml) for details.

## Customizing Configuration

To adjust linter settings, edit `.golangci.yml`:

```yaml
linters-settings:
  funlen:
    lines: 100  # Adjust max function length
    statements: 50
  
  gocognit:
    min-complexity: 20  # Adjust cognitive complexity threshold
```

## Recommended Workflow

1. **Before committing**:
   ```bash
   golangci-lint run --fix
   ```

2. **Review issues**:
   ```bash
   golangci-lint run
   ```

3. **Fix remaining issues** manually

4. **Commit** clean code

## Resources

- [golangci-lint Documentation](https://golangci-lint.run/)
- [Linters Reference](https://golangci-lint.run/usage/linters/)
- [Configuration Reference](https://github.com/golangci/golangci-lint/blob/master/.golangci.reference.yml)
- [Our Config](./.golangci.yml)
