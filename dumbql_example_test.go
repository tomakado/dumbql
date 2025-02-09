package dumbql_test

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/defer-panic/dumbql"
	"github.com/defer-panic/dumbql/schema"
)

func ExampleParse() {
	const q = `profile.age >= 18 and profile.city = Barcelona`
	ast, err := dumbql.Parse(q)
	if err != nil {
		panic(err)
	}

	fmt.Println(ast)
	// Output: (and (>= profile.age 18) (= profile.city "Barcelona"))
}

func ExampleQuery_Validate() {
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
	fmt.Println(err)
	// Output: (and (= status "pending") (= title "hello world"))
	// field "period_months": value must be equal or less than 3, got 4; field "name" not found in schema
}

func ExampleQuery_ToSql() {
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
	// Output: SELECT * FROM users WHERE ((status = ? AND period_months < ?) AND (title = ? OR name = ?))
	// [pending 4 hello world John Doe]
}
