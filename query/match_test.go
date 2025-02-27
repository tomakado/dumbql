package query_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.tomakado.io/dumbql/match"
	"go.tomakado.io/dumbql/query"
)

type person struct {
	Name     string  `dumbql:"name"`
	Age      int64   `dumbql:"age"`
	Height   float64 `dumbql:"height"`
	IsMember bool
}

func TestBinaryExpr_Match(t *testing.T) { //nolint:funlen
	target := person{Name: "John", Age: 30}
	matcher := &match.StructMatcher{}

	tests := []struct {
		name string
		expr *query.BinaryExpr
		want bool
	}{
		{
			name: "AND - both true",
			expr: &query.BinaryExpr{
				Left: &query.FieldExpr{
					Field: "name",
					Op:    query.Equal,
					Value: &query.StringLiteral{StringValue: "John"},
				},
				Op: query.And,
				Right: &query.FieldExpr{
					Field: "age",
					Op:    query.Equal,
					Value: &query.NumberLiteral{NumberValue: 30},
				},
			},
			want: true,
		},
		{
			name: "AND - left false",
			expr: &query.BinaryExpr{
				Left: &query.FieldExpr{
					Field: "name",
					Op:    query.Equal,
					Value: &query.StringLiteral{StringValue: "Jane"},
				},
				Op: query.And,
				Right: &query.FieldExpr{
					Field: "age",
					Op:    query.Equal,
					Value: &query.NumberLiteral{NumberValue: 30},
				},
			},
			want: false,
		},
		{
			name: "OR - both true",
			expr: &query.BinaryExpr{
				Left: &query.FieldExpr{
					Field: "name",
					Op:    query.Equal,
					Value: &query.StringLiteral{StringValue: "John"},
				},
				Op: query.Or,
				Right: &query.FieldExpr{
					Field: "age",
					Op:    query.Equal,
					Value: &query.NumberLiteral{NumberValue: 30},
				},
			},
			want: true,
		},
		{
			name: "OR - one true",
			expr: &query.BinaryExpr{
				Left: &query.FieldExpr{
					Field: "name",
					Op:    query.Equal,
					Value: &query.StringLiteral{StringValue: "John"},
				},
				Op: query.Or,
				Right: &query.FieldExpr{
					Field: "age",
					Op:    query.Equal,
					Value: &query.NumberLiteral{NumberValue: 25},
				},
			},
			want: true,
		},
		{
			name: "OR - both false",
			expr: &query.BinaryExpr{
				Left: &query.FieldExpr{
					Field: "name",
					Op:    query.Equal,
					Value: &query.StringLiteral{StringValue: "Jane"},
				},
				Op: query.Or,
				Right: &query.FieldExpr{
					Field: "age",
					Op:    query.Equal,
					Value: &query.NumberLiteral{NumberValue: 25},
				},
			},
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.expr.Match(target, matcher)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestNotExpr_Match(t *testing.T) {
	target := person{Name: "John", Age: 30}
	matcher := &match.StructMatcher{}

	tests := []struct {
		name string
		expr *query.NotExpr
		want bool
	}{
		{
			name: "negate true condition",
			expr: &query.NotExpr{
				Expr: &query.FieldExpr{
					Field: "name",
					Op:    query.Equal,
					Value: &query.StringLiteral{StringValue: "John"},
				},
			},
			want: false,
		},
		{
			name: "negate false condition",
			expr: &query.NotExpr{
				Expr: &query.FieldExpr{
					Field: "name",
					Op:    query.Equal,
					Value: &query.StringLiteral{StringValue: "Jane"},
				},
			},
			want: true,
		},
		{
			name: "negate AND expression",
			expr: &query.NotExpr{
				Expr: &query.BinaryExpr{
					Left: &query.FieldExpr{
						Field: "name",
						Op:    query.Equal,
						Value: &query.StringLiteral{StringValue: "John"},
					},
					Op: query.And,
					Right: &query.FieldExpr{
						Field: "age",
						Op:    query.Equal,
						Value: &query.NumberLiteral{NumberValue: 30},
					},
				},
			},
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.expr.Match(target, matcher)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestFieldExpr_Match(t *testing.T) { //nolint:funlen
	target := person{
		Name:     "John",
		Age:      30,
		Height:   1.75,
		IsMember: true,
	}
	matcher := &match.StructMatcher{}

	tests := []struct {
		name string
		expr *query.FieldExpr
		want bool
	}{
		{
			name: "string equal - match",
			expr: &query.FieldExpr{
				Field: "name",
				Op:    query.Equal,
				Value: &query.StringLiteral{StringValue: "John"},
			},
			want: true,
		},
		{
			name: "string not equal - match",
			expr: &query.FieldExpr{
				Field: "name",
				Op:    query.NotEqual,
				Value: &query.StringLiteral{StringValue: "Jane"},
			},
			want: true,
		},
		{
			name: "integer greater than - match",
			expr: &query.FieldExpr{
				Field: "age",
				Op:    query.GreaterThan,
				Value: &query.NumberLiteral{NumberValue: 25},
			},
			want: true,
		},
		{
			name: "float less than - match",
			expr: &query.FieldExpr{
				Field: "height",
				Op:    query.LessThan,
				Value: &query.NumberLiteral{NumberValue: 1.80},
			},
			want: true,
		},
		{
			name: "non-existent field",
			expr: &query.FieldExpr{
				Field: "invalid",
				Op:    query.Equal,
				Value: &query.StringLiteral{StringValue: "test"},
			},
			want: true,
		},
		{
			name: "field without dumbql tag",
			expr: &query.FieldExpr{
				Field: "IsMember",
				Op:    query.Equal,
				Value: &query.StringLiteral{StringValue: "true"},
			},
			want: false,
		},
		{
			name: "type mismatch",
			expr: &query.FieldExpr{
				Field: "age",
				Op:    query.Equal,
				Value: &query.StringLiteral{StringValue: "30"},
			},
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.expr.Match(target, matcher)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestIdentifier_Match(t *testing.T) { //nolint:funlen
	tests := []struct {
		name   string
		id     query.Identifier
		target any
		op     query.FieldOperator
		want   bool
	}{
		{
			name:   "equal - match",
			id:     query.Identifier("test"),
			target: "test",
			op:     query.Equal,
			want:   true,
		},
		{
			name:   "equal - no match",
			id:     query.Identifier("test"),
			target: "other",
			op:     query.Equal,
			want:   false,
		},
		{
			name:   "not equal - match",
			id:     query.Identifier("test"),
			target: "other",
			op:     query.NotEqual,
			want:   true,
		},
		{
			name:   "not equal - no match",
			id:     query.Identifier("test"),
			target: "test",
			op:     query.NotEqual,
			want:   false,
		},
		{
			name:   "like - match",
			id:     query.Identifier("world"),
			target: "hello world",
			op:     query.Like,
			want:   true,
		},
		{
			name:   "like - no match",
			id:     query.Identifier("universe"),
			target: "hello world",
			op:     query.Like,
			want:   false,
		},
		{
			name:   "with non-string target",
			id:     query.Identifier("42"),
			target: 42,
			op:     query.Equal,
			want:   false,
		},
		{
			name:   "with invalid operator",
			id:     query.Identifier("test"),
			target: "test",
			op:     query.GreaterThan,
			want:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.id.Match(test.target, test.op)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestOneOfExpr_Match(t *testing.T) { //nolint:funlen
	tests := []struct {
		name   string
		expr   *query.OneOfExpr
		target any
		op     query.FieldOperator
		want   bool
	}{
		{
			name: "string equal - match",
			expr: &query.OneOfExpr{
				Values: []query.Valuer{
					&query.StringLiteral{StringValue: "apple"},
					&query.StringLiteral{StringValue: "banana"},
					&query.StringLiteral{StringValue: "orange"},
				},
			},
			target: "banana",
			op:     query.Equal,
			want:   true,
		},
		{
			name: "string equal - no match",
			expr: &query.OneOfExpr{
				Values: []query.Valuer{
					&query.StringLiteral{StringValue: "apple"},
					&query.StringLiteral{StringValue: "banana"},
					&query.StringLiteral{StringValue: "orange"},
				},
			},
			target: "grape",
			op:     query.Equal,
			want:   false,
		},
		{
			name: "integer equal - match",
			expr: &query.OneOfExpr{
				Values: []query.Valuer{
					&query.NumberLiteral{NumberValue: 1},
					&query.NumberLiteral{NumberValue: 2},
					&query.NumberLiteral{NumberValue: 3},
				},
			},
			target: int64(2),
			op:     query.Equal,
			want:   true,
		},
		{
			name: "integer equal - no match",
			expr: &query.OneOfExpr{
				Values: []query.Valuer{
					&query.NumberLiteral{NumberValue: 1},
					&query.NumberLiteral{NumberValue: 2},
					&query.NumberLiteral{NumberValue: 3},
				},
			},
			target: int64(4),
			op:     query.Equal,
			want:   false,
		},
		{
			name: "float equal - match",
			expr: &query.OneOfExpr{
				Values: []query.Valuer{
					&query.NumberLiteral{NumberValue: 1.1},
					&query.NumberLiteral{NumberValue: 2.2},
					&query.NumberLiteral{NumberValue: 3.3},
				},
			},
			target: 2.2,
			op:     query.Equal,
			want:   true,
		},
		{
			name: "float equal - no match",
			expr: &query.OneOfExpr{
				Values: []query.Valuer{
					&query.NumberLiteral{NumberValue: 1.1},
					&query.NumberLiteral{NumberValue: 2.2},
					&query.NumberLiteral{NumberValue: 3.3},
				},
			},
			target: 4.4,
			op:     query.Equal,
			want:   false,
		},
		{
			name: "mixed types",
			expr: &query.OneOfExpr{
				Values: []query.Valuer{
					&query.StringLiteral{StringValue: "one"},
					&query.NumberLiteral{NumberValue: 2},
					&query.NumberLiteral{NumberValue: 3.3},
				},
			},
			target: "one",
			op:     query.Equal,
			want:   true,
		},
		{
			name: "empty values",
			expr: &query.OneOfExpr{
				Values: []query.Valuer{},
			},
			target: "test",
			op:     query.Equal,
			want:   false,
		},
		{
			name: "nil values",
			expr: &query.OneOfExpr{
				Values: nil,
			},
			target: "test",
			op:     query.Equal,
			want:   false,
		},
		{
			name: "string like - match",
			expr: &query.OneOfExpr{
				Values: []query.Valuer{
					&query.StringLiteral{StringValue: "world"},
					&query.StringLiteral{StringValue: "universe"},
				},
			},
			target: "hello world",
			op:     query.Like,
			want:   true,
		},
		{
			name: "invalid operator",
			expr: &query.OneOfExpr{
				Values: []query.Valuer{
					&query.StringLiteral{StringValue: "test"},
				},
			},
			target: "test",
			op:     query.GreaterThan,
			want:   false,
		},
		{
			name: "type mismatch",
			expr: &query.OneOfExpr{
				Values: []query.Valuer{
					&query.StringLiteral{StringValue: "42"},
				},
			},
			target: 42,
			op:     query.Equal,
			want:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.expr.Match(test.target, test.op)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestStructFieldOmission(t *testing.T) { //nolint:funlen
	type User struct {
		ID       int64   `dumbql:"id"`
		Name     string  `dumbql:"name"`
		Password string  `dumbql:"-"` // Should always match
		Internal bool    `dumbql:"-"` // Should always match
		Score    float64 `dumbql:"score"`
	}

	matcher := &match.StructMatcher{}
	user := &User{
		ID:       1,
		Name:     "John",
		Password: "secret123",
		Internal: true,
		Score:    4.5,
	}

	tests := []struct {
		name  string
		query string
		want  bool
	}{
		{
			name:  "visible field",
			query: `id:1`,
			want:  true,
		},
		{
			name:  "multiple visible fields",
			query: `id:1 and name:"John" and score:4.5`,
			want:  true,
		},
		{
			name:  "omitted field - always true",
			query: `password:"wrong_password"`,
			want:  true,
		},
		{
			name:  "another omitted field - always true",
			query: `internal:false`,
			want:  true,
		},
		{
			name:  "visible and omitted fields",
			query: `id:1 and password:"wrong_password"`,
			want:  true,
		},
		{
			name:  "non-existent field",
			query: `unknown:"value"`,
			want:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ast, err := query.Parse("test", []byte(test.query))
			require.NoError(t, err)
			expr := ast.(query.Expr)

			got := expr.Match(user, matcher)
			assert.Equal(t, test.want, got)
		})
	}
}
