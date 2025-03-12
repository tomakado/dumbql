package examples

import (
	"time"

	"go.tomakado.io/dumbql"
)

//go:generate go run go.tomakado.io/dumbql/cmd/dumbqlgen -type User -package . -output user_matcher.gen.go

// User is a sample struct to demonstrate the code generator
type User struct {
	ID        int64  `dumbql:"id"`
	Name      string `dumbql:"name"`
	Email     string `dumbql:"email"`
	CreatedAt time.Time
	Address   Address `dumbql:"address"`
	Private   bool    `dumbql:"-"` // Skip this field
}

// Address is a nested struct
type Address struct {
	Street string `dumbql:"street"`
	City   string `dumbql:"city"`
	State  string `dumbql:"state"`
	Zip    string `dumbql:"zip"`
}

// Below is an example of how to use the generated matcher
func Example() {
	// Create a user
	user := User{
		ID:        123,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
		Address: Address{
			Street: "123 Main St",
			City:   "Anytown",
			State:  "CA",
			Zip:    "12345",
		},
	}

	// Create a query
	q, _ := dumbql.Parse(`name = "John Doe" AND email = "john@example.com"`)

	// Use the generated matcher (this would be generated after running go generate)
	matcher := NewUserMatcher()

	// Match the query against the user
	result := q.Match(&user, matcher)
	_ = result // result will be true in this case
}
