package dumbql

import (
	"github.com/defer-panic/dumbql/query"
)

func Parse(q string) (query.Expr, error) {
	res, err := query.Parse("query", []byte(q))
	if err != nil {
		return nil, err
	}

	return res.(query.Expr), nil
}
