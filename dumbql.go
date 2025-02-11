// Package dumbql provides simple (dumb) query language and it's parser.
package dumbql

import (
	"github.com/defer-panic/dumbql/query"
	"github.com/defer-panic/dumbql/schema"
)

type Query struct {
	query.Expr
}

// Parse parses the input query string q, returning a Query reference or an error in case of invalid input.
func Parse(q string, opts ...query.Option) (*Query, error) {
	res, err := query.Parse("query", []byte(q), opts...)
	if err != nil {
		return nil, err
	}

	return &Query{res.(query.Expr)}, nil
}

// Validate checks the query against the provided schema, returning a validated expression or an error
// if any rule is violated. Even when error returned Validate can return query AST with invalided nodes dropped.
func (q *Query) Validate(s schema.Schema) (query.Expr, error) {
	return q.Expr.Validate(s)
}

// ToSql converts the Query into an SQL string, returning the SQL string, arguments slice,
// and any potential error encountered.
func (q *Query) ToSql() (string, []any, error) { //nolint:revive
	return q.Expr.ToSql()
}
