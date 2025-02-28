package dumbql_test

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"go.tomakado.io/dumbql"
	"go.tomakado.io/dumbql/schema"
)

func ExampleParse() {
	const q = `profile.age >= 18 and profile.city = Barcelona and profile.verified = true`
	ast, err := dumbql.Parse(q)
	if err != nil {
		panic(err)
	}

	fmt.Println(ast)
	// Output: (and (and (>= profile.age 18) (= profile.city "Barcelona")) (= profile.verified true))
}

func ExampleQuery_Validate() {
	schm := schema.Schema{
		"status": schema.All(
			schema.Is[string](),
			schema.EqualsOneOf("pending", "approved", "rejected"),
		),
		"period_months": schema.Max(int64(3)),
		"title":         schema.LenInRange(1, 100),
		"active":        schema.Is[bool](),
	}

	// The following query is invalid against the schema:
	// 	- period_months == 4, but max allowed value is 3
	// 	- field `name` is not described in the schema
	//
	// Invalid parts of the query are dropped.
	const q = `status:pending and period_months:4 and active:true and (title:"hello world" or name:"John Doe")`
	expr, err := dumbql.Parse(q)
	if err != nil {
		panic(err)
	}

	validated, err := expr.Validate(schm)
	fmt.Println(validated)
	fmt.Println(err)
	// Output: (and (and (= status "pending") (= active true)) (= title "hello world"))
	// field "period_months": value must be equal or less than 3, got 4; field "name" not found in schema
}

func ExampleQuery_ToSql() {
	const q = `status:pending and period_months < 4 and is_active:true and (title:"hello world" or name:"John Doe")`
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
	// nolint:lll
	// Output: SELECT * FROM users WHERE (((status = ? AND period_months < ?) AND is_active = ?) AND (title = ? OR name = ?))
	// [pending 4 true hello world John Doe]
}

func ExampleParse_booleanFields() {
	const q = `verified and premium and not banned and (admin or moderator)`
	ast, err := dumbql.Parse(q)
	if err != nil {
		panic(err)
	}

	fmt.Println(ast)
	//nolint:lll
	// Output: (and (and (and (= verified true) (= premium true)) (not (= banned true))) (or (= admin true) (= moderator true)))
}
