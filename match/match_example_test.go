package match_test

import (
	"fmt"

	"go.tomakado.io/dumbql/match"
	"go.tomakado.io/dumbql/query"
)

type MatchUser struct {
	ID       int64   `dumbql:"id"`
	Name     string  `dumbql:"name"`
	Age      int64   `dumbql:"age"`
	Score    float64 `dumbql:"score"`
	Location string  `dumbql:"location"`
	Role     string  `dumbql:"role"`
}

func Example() {
	users := []MatchUser{
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
	ast, _ := query.Parse("test", []byte(q))
	expr := ast.(query.Expr)

	matcher := &match.StructMatcher{}

	filtered := make([]MatchUser, 0, len(users))

	for _, user := range users {
		if expr.Match(&user, matcher) {
			filtered = append(filtered, user)
		}
	}

	fmt.Println(filtered)
	// Output:
	// [{1 John Doe 30 4.5 New York admin} {2 Jane Smith 25 3.8 Los Angeles user} {3 Bob Johnson 35 4.2 Chicago user}]
}
