package query_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.tomakado.io/dumbql/query"
)

func TestLiteralValues(t *testing.T) {
	t.Run("StringLiteral.Value", func(t *testing.T) {
		sl := &query.StringLiteral{StringValue: "hello"}
		assert.Equal(t, "hello", sl.Value())
	})

	t.Run("NumberLiteral.Value", func(t *testing.T) {
		nl := &query.NumberLiteral{NumberValue: 42.5}
		assert.InDelta(t, 42.5, nl.Value(), 0.0001)
	})

	t.Run("BoolLiteral.Value", func(t *testing.T) {
		bl1 := &query.BoolLiteral{BoolValue: true}
		assert.Equal(t, true, bl1.Value())

		bl2 := &query.BoolLiteral{BoolValue: false}
		assert.Equal(t, false, bl2.Value())
	})

	t.Run("Identifier.Value", func(t *testing.T) {
		id := query.Identifier("field")
		assert.Equal(t, "field", id.Value())
	})
}

// Test String methods for operators
func TestOperatorString(t *testing.T) {
	t.Run("BooleanOperator.String", func(t *testing.T) {
		// Test valid operators
		assert.Equal(t, "and", query.And.String())
		assert.Equal(t, "or", query.Or.String())
		
		// Test invalid operator (default case)
		type CustomBoolOp query.BooleanOperator
		invalidOp := query.BooleanOperator(CustomBoolOp(255)) // Invalid operator (max uint8)
		assert.Equal(t, "unknown!", invalidOp.String())
	})

	t.Run("FieldOperator.String", func(t *testing.T) {
		// Test all valid operators
		assert.Equal(t, "=", query.Equal.String())
		assert.Equal(t, "!=", query.NotEqual.String())
		assert.Equal(t, ">", query.GreaterThan.String())
		assert.Equal(t, ">=", query.GreaterThanOrEqual.String())
		assert.Equal(t, "<", query.LessThan.String())
		assert.Equal(t, "<=", query.LessThanOrEqual.String())
		assert.Equal(t, "~", query.Like.String())
		
		// Test invalid operator (default case)
		type CustomFieldOp query.FieldOperator
		invalidOp := query.FieldOperator(CustomFieldOp(255)) // Invalid operator (max uint8)
		assert.Equal(t, "unknown!", invalidOp.String())
	})
}
