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

## Project Structure
- /query - Query parsing, validation, and SQL generation
- /match - Struct matching functionality
- /schema - Schema definition and validation rules