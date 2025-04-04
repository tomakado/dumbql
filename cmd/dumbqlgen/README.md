# DumbQL Code Generator

`dumbqlgen` is a code generation tool for DumbQL that generates type-specific matchers for structs.
This helps avoid reflection and provides better performance for matching queries.

## Installation

```bash
go install go.tomakado.io/dumbql/cmd/dumbqlgen@latest
```

## Usage

You can use the code generator in two ways:

### 1. Directly with go generate

Add a go:generate comment in your code:

```go
//go:generate go run go.tomakado.io/dumbql/cmd/dumbqlgen -type User -package .
```

Then run:

```bash
go generate ./...
```

### 2. Manually at the command line

```bash
dumbqlgen -type User -package ./path/to/package -output user_matcher.gen.go
```

## Flags

- `-type`: (Required) The name of the struct to generate a matcher for
- `-package`: (Optional) The path to the package containing the struct (defaults to ".")
- `-output`: (Optional) The output file path (defaults to "lowercase_type_matcher.gen.go")

## Example

Given a struct:

```go
type User struct {
    ID        int64  `dumbql:"id"`
    Name      string `dumbql:"name"`
    Email     string `dumbql:"email"`
    CreatedAt time.Time
    Private   bool   `dumbql:"-"` // Skip this field
}
```

Running:

```bash
dumbqlgen -type User -package .
```

Will generate:

```go
// Code generated by dumbqlgen; DO NOT EDIT.
package mypackage

import (
    "strings"
    
    "go.tomakado.io/dumbql/match"
)

// UserRouter is a generated Router implementation for User.
// It implements the match.Router interface.
type UserRouter struct{}

// NewUserMatcher creates a new StructMatcher with a generated router for User.
func NewUserMatcher() *match.StructMatcher {
    return match.NewStructMatcher(&UserRouter{})
}

// Route method and other implementation details...
```

## Benefits

- No reflection at runtime for better performance
- Type-safe field access
- Works with the existing DumbQL query engine
- Supports struct tags and nested fields