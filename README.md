<div align="center">
<h1>DumbQL</h1>
    
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/tomakado/dumbql) ![GitHub License](https://img.shields.io/github/license/tomakado/dumbql) ![GitHub Tag](https://img.shields.io/github/v/tag/tomakado/dumbql) [![Go Report Card](https://goreportcard.com/badge/go.tomakado.io/dumbql)](https://goreportcard.com/report/go.tomakado.io/dumbql) [![CI](https://github.com/tomakado/dumbql/actions/workflows/main.yml/badge.svg)](https://github.com/tomakado/dumbql/actions/workflows/main.yml) [![codecov](https://codecov.io/gh/tomakado/dumbql/graph/badge.svg?token=15IWJO0R0K)](https://codecov.io/gh/tomakado/dumbql) [![Go Reference](https://pkg.go.dev/badge/go.tomakado.io/dumbql.svg)](https://pkg.go.dev/go.tomakado.io/dumbql)

Simple (dumb?) query language and parser for Go.

</div>

## Features

- Field expressions (`age >= 18`, `field.name:"field value"`, etc.)
- Boolean expressions (`age >= 18 and city = Barcelona`, `occupation = designer or occupation = "ux analyst"`)
- One-of/In expressions (`occupation = [designer, "ux analyst"]`)
- Boolean fields support with shorthand syntax (`is_active`, `verified and premium`)
- Schema validation
- Drop-in usage with [squirrel](https://github.com/Masterminds/squirrel) or SQL drivers directly
- Struct matching with `dumbql` struct tag

## Examples

### Simple parse

```go
package main

import (
    "fmt"

    "go.tomakado.io/dumbql"
)

func main() {
    const q = `profile.age >= 18 and profile.city = Barcelona`
    ast, err := dumbql.Parse(q)
    if err != nil {
        panic(err)
    }

    fmt.Println(ast)
    // Output: (and (>= profile.age 18) (= profile.city "Barcelona"))
}
```

### Validation against schema

```go
package main

import (
    "fmt"

    "go.tomakado.io/dumbql"
    "go.tomakado.io/dumbql/schema"
)

func main() {
    schm := schema.Schema{
        "status": schema.All(
            schema.Is[string](),
            schema.EqualsOneOf("pending", "approved", "rejected"),
        ),
        "period_months": schema.Max(int64(3)),
        "title":         schema.LenInRange(1, 100),
    }

    // The following query is invalid against the schema:
    // 	- period_months == 4, but max allowed value is 3
    // 	- field `name` is not described in the schema
    //
    // Invalid parts of the query are dropped.
    const q = `status:pending and period_months:4 and (title:"hello world" or name:"John Doe")`
    expr, err := dumbql.Parse(q)
    if err != nil {
        panic(err)
    }

    validated, err := expr.Validate(schm)
    fmt.Println(validated)
    fmt.Printf("validation error: %v\n", err)
    // Output: 
    // (and (= status "pending") (= title "hello world"))
    // validation error: field "period_months": value must be equal or less than 3, got 4; field "name" not found in schema
}
```

### Convert to SQL

```go
package main

import (
  "fmt"

  sq "github.com/Masterminds/squirrel"
  "go.tomakado.io/dumbql"
)

func main() {
  const q = `status:pending and period_months < 4 and (title:"hello world" or name:"John Doe")`
  expr, err := dumbql.Parse(q)
  if err != nil {
    panic(err)
  }

  sql, args, err := sq.Select("*").
    From("users").
    Where(expr).
    ToSql()
  if err != nil {
    panic(err)
  }

  fmt.Println(sql)
  fmt.Println(args)
  // Output: 
  // SELECT * FROM users WHERE ((status = ? AND period_months < ?) AND (title = ? OR name = ?))
  // [pending 4 hello world John Doe]
}
```

See [dumbql_example_test.go](dumbql_example_test.go)

### Match against structs

```go
package main

import (
  "fmt"

  "go.tomakado.io/dumbql"
  "go.tomakado.io/dumbql/match"
  "go.tomakado.io/dumbql/query"
)

type User struct {
  ID       int64   `dumbql:"id"`
  Name     string  `dumbql:"name"`
  Age      int64   `dumbql:"age"`
  Score    float64 `dumbql:"score"`
  Location string  `dumbql:"location"`
  Role     string  `dumbql:"role"`
}

func main() {
  users := []User{
    {
      ID:       1,
      Name:     "John Doe",
      Age:      30,
      Score:    4.5,
      Location: "New York",
      Role:     "admin",
    },
    {
      ID:       2,
      Name:     "Jane Smith",
      Age:      25,
      Score:    3.8,
      Location: "Los Angeles",
      Role:     "user",
    },
    {
      ID:       3,
      Name:     "Bob Johnson",
      Age:      35,
      Score:    4.2,
      Location: "Chicago",
      Role:     "user",
    },
    // This one will be dropped:
    {
      ID:       4,
      Name:     "Alice Smith",
      Age:      25,
      Score:    3.8,
      Location: "Los Angeles",
      Role:     "admin",
    },
  }

  q := `(age >= 30 and score > 4.0) or (location:"Los Angeles" and role:"user")`
  ast, err := dumbql.Parse(q)
  if err != nil {
    panic(err)
  }
  matcher := &match.StructMatcher{}
  filtered := make([]User, 0, len(users))

  for _, user := range users {
    if expr.Match(&user, matcher) {
      filtered = append(filtered, user)
    }
  }

  fmt.Println(filtered)
  // [{1 John Doe 30 4.5 New York admin} {2 Jane Smith 25 3.8 Los Angeles user} {3 Bob Johnson 35 4.2 Chicago user}]
}
```

See [match_example_test.go](match_example_test.go) for more examples.

## Query syntax

This section is a non-formal description of DumbQL syntax. For strict description see [grammar file](query/grammar.peg).

### Field expression

Field name & value pair divided by operator. Field name is any alphanumeric identifier (with underscore), value can be string, int64, float64, or bool.
One-of expression is also supported (see below).

```
<field_name> <operator> <value>
```

for example

```
period_months < 4
is_active:true
```

### Field expression operators

| Operator             | Meaning                       | Supported types                      |
|----------------------|-------------------------------|--------------------------------------|
| `:` or `=`           | Equal, one of                 | `int64`, `float64`, `string`, `bool` |
| `!=` or `!:`         | Not equal                     | `int64`, `float64`, `string`, `bool` |
| `~`                  | "Like" or "contains" operator | `string`                             |
| `>`, `>=`, `<`, `<=` | Comparison                    | `int64`, `float64`                   |


### Boolean operators

Multiple field expression can be combined into boolean expressions with `and` (`AND`) or `or` (`OR`) operators:

```
status:pending and period_months < 4 and (title:"hello world" or name:"John Doe")
```

### Boolean Field Shorthand

Boolean fields can be expressed in a simpler shorthand syntax:

```
verified                  # equivalent to verified:true
verified and premium      # equivalent to verified:true and premium:true  
not verified              # equivalent to not (verified:true)
verified or admin         # equivalent to verified:true or admin:true
```

### "One of" expression

Sometimes instead of multiple `and`/`or` clauses against the same field:

```
occupation = designer or occupation = "ux analyst"
```

it's more convenient to use equivalent “one of” expressions:

```
occupation: [designer, "ux analyst"]
```

### Numbers

If number does not have digits after `.` it's treated as integer and stored as `int64`. And it's `float64` otherwise.

### Strings

String is a sequence of Unicode characters surrounded by double quotes (`"`). In some cases like single word it's possible to write string value without double quotes.

### Booleans

Boolean values are represented by `true` or `false` literals and can be used with the equality operators (`=`, `:`, `!=`, `!:`).

```
is_active:true
verified = true
is_banned != false
```
