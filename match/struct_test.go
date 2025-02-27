package match_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.tomakado.io/dumbql/match"
	"go.tomakado.io/dumbql/query"
)

type person struct {
	Name     string  `dumbql:"name"`
	Age      int64   `dumbql:"age"`
	Height   float64 `dumbql:"height"`
	IsMember bool
}

func TestStructMatcher_MatchAnd(t *testing.T) { //nolint:funlen
	matcher := &match.StructMatcher{}
	target := person{Name: "John", Age: 30}

	tests := []struct {
		name  string
		left  query.Expr
		right query.Expr
		want  bool
	}{
		{
			name: "both conditions true",
			left: &query.FieldExpr{
				Field: "name",
				Op:    query.Equal,
				Value: &query.StringLiteral{StringValue: "John"},
			},
			right: &query.FieldExpr{
				Field: "age",
				Op:    query.Equal,
				Value: &query.NumberLiteral{NumberValue: 30},
			},
			want: true,
		},
		{
			name: "left condition false",
			left: &query.FieldExpr{
				Field: "name",
				Op:    query.Equal,
				Value: &query.StringLiteral{StringValue: "Jane"},
			},
			right: &query.FieldExpr{
				Field: "age",
				Op:    query.Equal,
				Value: &query.NumberLiteral{NumberValue: 30},
			},
			want: false,
		},
		{
			name: "right condition false",
			left: &query.FieldExpr{
				Field: "name",
				Op:    query.Equal,
				Value: &query.StringLiteral{StringValue: "John"},
			},
			right: &query.FieldExpr{
				Field: "age",
				Op:    query.Equal,
				Value: &query.NumberLiteral{NumberValue: 25},
			},
			want: false,
		},
		{
			name: "both conditions false",
			left: &query.FieldExpr{
				Field: "name",
				Op:    query.Equal,
				Value: &query.StringLiteral{StringValue: "Jane"},
			},
			right: &query.FieldExpr{
				Field: "age",
				Op:    query.Equal,
				Value: &query.NumberLiteral{NumberValue: 25},
			},
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := matcher.MatchAnd(target, test.left, test.right)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestStructMatcher_MatchOr(t *testing.T) { //nolint:funlen
	matcher := &match.StructMatcher{}
	target := person{Name: "John", Age: 30}

	tests := []struct {
		name  string
		left  query.Expr
		right query.Expr
		want  bool
	}{
		{
			name: "both conditions true",
			left: &query.FieldExpr{
				Field: "name",
				Op:    query.Equal,
				Value: &query.StringLiteral{StringValue: "John"},
			},
			right: &query.FieldExpr{
				Field: "age",
				Op:    query.Equal,
				Value: &query.NumberLiteral{NumberValue: 30},
			},
			want: true,
		},
		{
			name: "left condition true only",
			left: &query.FieldExpr{
				Field: "name",
				Op:    query.Equal,
				Value: &query.StringLiteral{StringValue: "John"},
			},
			right: &query.FieldExpr{
				Field: "age",
				Op:    query.Equal,
				Value: &query.NumberLiteral{NumberValue: 25},
			},
			want: true,
		},
		{
			name: "right condition true only",
			left: &query.FieldExpr{
				Field: "name",
				Op:    query.Equal,
				Value: &query.StringLiteral{StringValue: "Jane"},
			},
			right: &query.FieldExpr{
				Field: "age",
				Op:    query.Equal,
				Value: &query.NumberLiteral{NumberValue: 30},
			},
			want: true,
		},
		{
			name: "both conditions false",
			left: &query.FieldExpr{
				Field: "name",
				Op:    query.Equal,
				Value: &query.StringLiteral{StringValue: "Jane"},
			},
			right: &query.FieldExpr{
				Field: "age",
				Op:    query.Equal,
				Value: &query.NumberLiteral{NumberValue: 25},
			},
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := matcher.MatchOr(target, test.left, test.right)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestStructMatcher_MatchNot(t *testing.T) {
	matcher := &match.StructMatcher{}
	target := person{Name: "John", Age: 30}

	tests := []struct {
		name string
		expr query.Expr
		want bool
	}{
		{
			name: "negate true condition",
			expr: &query.FieldExpr{
				Field: "name",
				Op:    query.Equal,
				Value: &query.StringLiteral{StringValue: "John"},
			},
			want: false,
		},
		{
			name: "negate false condition",
			expr: &query.FieldExpr{
				Field: "name",
				Op:    query.Equal,
				Value: &query.StringLiteral{StringValue: "Jane"},
			},
			want: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := matcher.MatchNot(target, test.expr)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestStructMatcher_MatchField(t *testing.T) {
	matcher := &match.StructMatcher{}
	target := person{
		Name:     "John",
		Age:      30,
		Height:   1.75,
		IsMember: true,
	}

	tests := []struct {
		name  string
		field string
		value query.Valuer
		op    query.FieldOperator
		want  bool
	}{
		{
			name:  "string equal match",
			field: "name",
			value: &query.StringLiteral{StringValue: "John"},
			op:    query.Equal,
			want:  true,
		},
		{
			name:  "string not equal match",
			field: "name",
			value: &query.StringLiteral{StringValue: "Jane"},
			op:    query.NotEqual,
			want:  true,
		},
		{
			name:  "integer equal match",
			field: "age",
			value: &query.NumberLiteral{NumberValue: 30},
			op:    query.Equal,
			want:  true,
		},
		{
			name:  "float greater than match",
			field: "height",
			value: &query.NumberLiteral{NumberValue: 1.70},
			op:    query.GreaterThan,
			want:  true,
		},
		{
			name:  "non-existent field",
			field: "invalid",
			value: &query.StringLiteral{StringValue: "test"},
			op:    query.Equal,
			want:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := matcher.MatchField(target, test.field, test.value, test.op)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestStructMatcher_MatchValue(t *testing.T) {
	t.Run("string", testMatchValueString)
	t.Run("integer", testMatchValueInteger)
	t.Run("float", testMatchValueFloat)
	t.Run("type mismatch", testMatchValueTypeMismatch)
}

func testMatchValueString(t *testing.T) { //nolint:funlen
	matcher := &match.StructMatcher{}
	tests := []struct {
		name   string
		target any
		value  query.Valuer
		op     query.FieldOperator
		want   bool
	}{
		{
			name:   "equal - match",
			target: "hello",
			value:  &query.StringLiteral{StringValue: "hello"},
			op:     query.Equal,
			want:   true,
		},
		{
			name:   "equal - no match",
			target: "hello",
			value:  &query.StringLiteral{StringValue: "world"},
			op:     query.Equal,
			want:   false,
		},
		{
			name:   "not equal - match",
			target: "hello",
			value:  &query.StringLiteral{StringValue: "world"},
			op:     query.NotEqual,
			want:   true,
		},
		{
			name:   "not equal - no match",
			target: "hello",
			value:  &query.StringLiteral{StringValue: "hello"},
			op:     query.NotEqual,
			want:   false,
		},
		{
			name:   "like - match",
			target: "hello world",
			value:  &query.StringLiteral{StringValue: "world"},
			op:     query.Like,
			want:   true,
		},
		{
			name:   "like - no match",
			target: "hello world",
			value:  &query.StringLiteral{StringValue: "universe"},
			op:     query.Like,
			want:   false,
		},
		{
			name:   "greater than - invalid",
			target: "hello",
			value:  &query.StringLiteral{StringValue: "world"},
			op:     query.GreaterThan,
			want:   false,
		},
		{
			name:   "greater than or equal - invalid",
			target: "hello",
			value:  &query.StringLiteral{StringValue: "world"},
			op:     query.GreaterThanOrEqual,
			want:   false,
		},
		{
			name:   "less than - invalid",
			target: "hello",
			value:  &query.StringLiteral{StringValue: "world"},
			op:     query.LessThan,
			want:   false,
		},
		{
			name:   "less than or equal - invalid",
			target: "hello",
			value:  &query.StringLiteral{StringValue: "world"},
			op:     query.LessThanOrEqual,
			want:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := matcher.MatchValue(test.target, test.value, test.op)
			assert.Equal(t, test.want, result)
		})
	}
}

func testMatchValueInteger(t *testing.T) { //nolint:funlen
	matcher := &match.StructMatcher{}
	tests := []struct {
		name   string
		target any
		value  query.Valuer
		op     query.FieldOperator
		want   bool
	}{
		{
			name:   "equal - match",
			target: int64(42),
			value:  &query.NumberLiteral{NumberValue: 42},
			op:     query.Equal,
			want:   true,
		},
		{
			name:   "equal - no match",
			target: int64(42),
			value:  &query.NumberLiteral{NumberValue: 24},
			op:     query.Equal,
			want:   false,
		},
		{
			name:   "not equal - match",
			target: int64(42),
			value:  &query.NumberLiteral{NumberValue: 24},
			op:     query.NotEqual,
			want:   true,
		},
		{
			name:   "not equal - no match",
			target: int64(42),
			value:  &query.NumberLiteral{NumberValue: 42},
			op:     query.NotEqual,
			want:   false,
		},
		{
			name:   "greater than - match",
			target: int64(42),
			value:  &query.NumberLiteral{NumberValue: 24},
			op:     query.GreaterThan,
			want:   true,
		},
		{
			name:   "greater than - no match",
			target: int64(24),
			value:  &query.NumberLiteral{NumberValue: 42},
			op:     query.GreaterThan,
			want:   false,
		},
		{
			name:   "greater than or equal - match (greater)",
			target: int64(42),
			value:  &query.NumberLiteral{NumberValue: 24},
			op:     query.GreaterThanOrEqual,
			want:   true,
		},
		{
			name:   "greater than or equal - match (equal)",
			target: int64(42),
			value:  &query.NumberLiteral{NumberValue: 42},
			op:     query.GreaterThanOrEqual,
			want:   true,
		},
		{
			name:   "greater than or equal - no match",
			target: int64(24),
			value:  &query.NumberLiteral{NumberValue: 42},
			op:     query.GreaterThanOrEqual,
			want:   false,
		},
		{
			name:   "less than - match",
			target: int64(24),
			value:  &query.NumberLiteral{NumberValue: 42},
			op:     query.LessThan,
			want:   true,
		},
		{
			name:   "less than - no match",
			target: int64(42),
			value:  &query.NumberLiteral{NumberValue: 24},
			op:     query.LessThan,
			want:   false,
		},
		{
			name:   "less than or equal - match (less)",
			target: int64(24),
			value:  &query.NumberLiteral{NumberValue: 42},
			op:     query.LessThanOrEqual,
			want:   true,
		},
		{
			name:   "less than or equal - match (equal)",
			target: int64(42),
			value:  &query.NumberLiteral{NumberValue: 42},
			op:     query.LessThanOrEqual,
			want:   true,
		},
		{
			name:   "less than or equal - no match",
			target: int64(42),
			value:  &query.NumberLiteral{NumberValue: 24},
			op:     query.LessThanOrEqual,
			want:   false,
		},
		{
			name:   "like - invalid",
			target: int64(42),
			value:  &query.NumberLiteral{NumberValue: 24},
			op:     query.Like,
			want:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := matcher.MatchValue(test.target, test.value, test.op)
			assert.Equal(t, test.want, result)
		})
	}
}

func testMatchValueFloat(t *testing.T) { //nolint:funlen
	matcher := &match.StructMatcher{}
	tests := []struct {
		name   string
		target any
		value  query.Valuer
		op     query.FieldOperator
		want   bool
	}{
		{
			name:   "equal - match",
			target: 3.14,
			value:  &query.NumberLiteral{NumberValue: 3.14},
			op:     query.Equal,
			want:   true,
		},
		{
			name:   "equal - no match",
			target: 3.14,
			value:  &query.NumberLiteral{NumberValue: 2.718},
			op:     query.Equal,
			want:   false,
		},
		{
			name:   "not equal - match",
			target: 3.14,
			value:  &query.NumberLiteral{NumberValue: 2.718},
			op:     query.NotEqual,
			want:   true,
		},
		{
			name:   "not equal - no match",
			target: 3.14,
			value:  &query.NumberLiteral{NumberValue: 3.14},
			op:     query.NotEqual,
			want:   false,
		},
		{
			name:   "greater than - match",
			target: 3.14,
			value:  &query.NumberLiteral{NumberValue: 2.718},
			op:     query.GreaterThan,
			want:   true,
		},
		{
			name:   "greater than - no match",
			target: 2.718,
			value:  &query.NumberLiteral{NumberValue: 3.14},
			op:     query.GreaterThan,
			want:   false,
		},
		{
			name:   "greater than or equal - match (greater)",
			target: 3.14,
			value:  &query.NumberLiteral{NumberValue: 2.718},
			op:     query.GreaterThanOrEqual,
			want:   true,
		},
		{
			name:   "greater than or equal - match (equal)",
			target: 3.14,
			value:  &query.NumberLiteral{NumberValue: 3.14},
			op:     query.GreaterThanOrEqual,
			want:   true,
		},
		{
			name:   "greater than or equal - no match",
			target: 2.718,
			value:  &query.NumberLiteral{NumberValue: 3.14},
			op:     query.GreaterThanOrEqual,
			want:   false,
		},
		{
			name:   "less than - match",
			target: 2.718,
			value:  &query.NumberLiteral{NumberValue: 3.14},
			op:     query.LessThan,
			want:   true,
		},
		{
			name:   "less than - no match",
			target: 3.14,
			value:  &query.NumberLiteral{NumberValue: 2.718},
			op:     query.LessThan,
			want:   false,
		},
		{
			name:   "less than or equal - match (less)",
			target: 2.718,
			value:  &query.NumberLiteral{NumberValue: 3.14},
			op:     query.LessThanOrEqual,
			want:   true,
		},
		{
			name:   "less than or equal - match (equal)",
			target: 3.14,
			value:  &query.NumberLiteral{NumberValue: 3.14},
			op:     query.LessThanOrEqual,
			want:   true,
		},
		{
			name:   "less than or equal - no match",
			target: 3.14,
			value:  &query.NumberLiteral{NumberValue: 2.718},
			op:     query.LessThanOrEqual,
			want:   false,
		},
		{
			name:   "like - invalid",
			target: 3.14,
			value:  &query.NumberLiteral{NumberValue: 2.718},
			op:     query.Like,
			want:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := matcher.MatchValue(test.target, test.value, test.op)
			assert.Equal(t, test.want, result)
		})
	}
}

func testMatchValueTypeMismatch(t *testing.T) {
	matcher := &match.StructMatcher{}
	tests := []struct {
		name   string
		target any
		value  query.Valuer
		op     query.FieldOperator
		want   bool
	}{
		{
			name:   "string target with number value",
			target: "42",
			value:  &query.NumberLiteral{NumberValue: 42},
			op:     query.Equal,
			want:   false,
		},
		{
			name:   "integer target with string value",
			target: int64(42),
			value:  &query.StringLiteral{StringValue: "42"},
			op:     query.Equal,
			want:   false,
		},
		{
			name:   "float target with string value",
			target: 3.14,
			value:  &query.StringLiteral{StringValue: "3.14"},
			op:     query.Equal,
			want:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := matcher.MatchValue(test.target, test.value, test.op)
			assert.Equal(t, test.want, result)
		})
	}
}