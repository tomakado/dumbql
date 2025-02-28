package query_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.tomakado.io/dumbql/query"
)

func TestDirectSqlGeneration(t *testing.T) {
	t.Run("StringLiteral", func(t *testing.T) {
		sl := &query.StringLiteral{StringValue: "hello"}
		sql, args, err := sl.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "?", sql)
		assert.Equal(t, []any{"hello"}, args)
	})

	t.Run("NumberLiteral", func(t *testing.T) {
		nl := &query.NumberLiteral{NumberValue: 42.5}
		sql, args, err := nl.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "?", sql)
		assert.InDelta(t, 42.5, args[0], 0.0001)
	})

	t.Run("BoolLiteral", func(t *testing.T) {
		bl := &query.BoolLiteral{BoolValue: true}
		sql, args, err := bl.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "?", sql)
		assert.Equal(t, []any{true}, args)
	})

	t.Run("Identifier", func(t *testing.T) {
		id := query.Identifier("fieldname")
		sql, args, err := id.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "fieldname", sql)
		assert.Empty(t, args)
	})

	t.Run("OneOfExpr", func(t *testing.T) {
		values := []query.Valuer{
			&query.StringLiteral{StringValue: "one"},
			&query.StringLiteral{StringValue: "two"},
			&query.NumberLiteral{NumberValue: 3},
		}
		oe := &query.OneOfExpr{Values: values}
		sql, args, err := oe.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "?", sql)
		// The Values should be a slice of string/number values
		assert.Contains(t, args[0], "one")
		assert.Contains(t, args[0], "two")
		assert.Contains(t, args[0], float64(3))
	})

	t.Run("BinaryExpr_OR_operator", func(t *testing.T) {
		// Testing the OR operator branch in BinaryExpr.ToSql
		left := &query.FieldExpr{
			Field: query.Identifier("status"),
			Op:    query.Equal,
			Value: &query.NumberLiteral{NumberValue: 200},
		}
		right := &query.FieldExpr{
			Field: query.Identifier("code"),
			Op:    query.Equal,
			Value: &query.NumberLiteral{NumberValue: 400},
		}
		be := &query.BinaryExpr{
			Left:  left,
			Op:    query.Or,
			Right: right,
		}
		sql, args, err := be.ToSql()
		assert.NoError(t, err)
		assert.Contains(t, sql, "OR")
		assert.Len(t, args, 2)
	})

	t.Run("BinaryExpr_unknown_operator", func(t *testing.T) {
		// Test the unknown operator branch
		// Create a custom type that embeds BooleanOperator but with a value not defined in the enum
		type CustomBoolOp query.BooleanOperator
		be := &query.BinaryExpr{
			Left:  &query.FieldExpr{},
			Op:    query.BooleanOperator(CustomBoolOp(255)), // Invalid operator (max uint8)
			Right: &query.FieldExpr{},
		}
		_, _, err := be.ToSql()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown operator")
	})

	t.Run("NotExpr_error_handling", func(t *testing.T) {
		// Test error handling in NotExpr.ToSql
		// Use a BinaryExpr with an invalid operator to generate an error
		type CustomBoolOp query.BooleanOperator
		invalidExpr := &query.BinaryExpr{
			Left:  &query.FieldExpr{},
			Op:    query.BooleanOperator(CustomBoolOp(255)), // Invalid operator (max uint8)
			Right: &query.FieldExpr{},
		}
		ne := &query.NotExpr{
			Expr: invalidExpr,
		}
		_, _, err := ne.ToSql()
		assert.Error(t, err)
	})

	t.Run("FieldExpr_unknown_operator", func(t *testing.T) {
		// Test the unknown operator branch in FieldExpr.ToSql
		type CustomFieldOp query.FieldOperator
		fe := &query.FieldExpr{
			Field: query.Identifier("status"),
			Op:    query.FieldOperator(CustomFieldOp(255)), // Invalid operator (max uint8)
			Value: &query.NumberLiteral{NumberValue: 200},
		}
		_, _, err := fe.ToSql()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown operator")
	})
}