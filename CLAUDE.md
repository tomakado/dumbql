# DumbQL Project Guide

## Build/Test Commands
- `task lint` - Run linting with golangci-lint
- `task test` - Run all tests with coverage
- `go test ./...` - Run all tests
- `go test ./path/to/package` - Run tests in specific package
- `go test -run TestName ./path/to/package` - Run specific test
- `task generate` - Run code generators
- `go run ./...` - Run the project

## Code Style Guidelines
- Naming: Use camelCase for variables, PascalCase for exported items
- Errors: Return meaningful error messages, use multierr for combining errors
- Formatting: Use standard Go formatting (`gofmt`)
- Types: Prefer strong typing, use int64/float64 consistently for numeric values
- Testing: Write thorough tests with testify, use table-driven tests
- Documentation: Document public APIs with complete sentences
- Structure: Keep packages focused on specific functionality
- Imports: Group standard library, external, and internal packages
- Complexity: Keep functions small and focused, avoid deep nesting
- In general, follow Go's best practices and idioms and the project's existing style

## Project Structure
- /query - Query parsing, validation, and SQL generation
- /match - Struct matching functionality
- /schema - Schema definition and validation rules

## Workflow Instructions
- When it's said that we're fixing the issue, request this issue with gh-cli and its comments to get the context
- Write tests before writing code
- After changes and before committing, run `task lint` and `task test`. If any of these fail, fix the issues before proceeding
- After creating PR, wait for CI/CD to pass
- Then check if test coverage decreased. If it did, add or modify tests to keep at least the same coverage
