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
