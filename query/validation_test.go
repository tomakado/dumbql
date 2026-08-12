package query_test

import (
	"errors"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.tomakado.io/dumbql/query"
	"go.tomakado.io/dumbql/schema"
)

func TestBinaryExpr_Validate(t *testing.T) { //nolint:funlen
	t.Run("positive", func(t *testing.T) {
		b := schema.NewSchemaBuilder()
		b.AddField("left", schema.Any())
		b.AddField("right", schema.Any())

		schm := b.Build()

		expr := &query.BinaryExpr{
			Left: &query.FieldExpr{
				Field: "left",
				Op:    query.Equal,
				Value: &query.NumberLiteral{NumberValue: 42},
			},
			Op: query.And,
			Right: &query.FieldExpr{
				Field: "right",
				Op:    query.Equal,
				Value: &query.NumberLiteral{NumberValue: math.Pi},
			},
		}

		got, err := expr.Validate(schm)
		require.NoError(t, err)

		binaryExpr, isBinaryExpr := got.(*query.BinaryExpr)
		require.True(t, isBinaryExpr)

		leftFieldExpr, isLeftFieldExpr := binaryExpr.Left.(*query.FieldExpr)
		require.True(t, isLeftFieldExpr)

		rightFieldExpr, isRightFieldExpr := binaryExpr.Right.(*query.FieldExpr)
		require.True(t, isRightFieldExpr)

		leftNumberLiteral, isLeftNumberLiteral := leftFieldExpr.Value.(*query.NumberLiteral)
		require.True(t, isLeftNumberLiteral)

		rightNumberLiteral, isRightNumberLiteral := rightFieldExpr.Value.(*query.NumberLiteral)
		require.True(t, isRightNumberLiteral)

		require.InDelta(t, 42.0, leftNumberLiteral.NumberValue, 0.0001)
		require.InDelta(t, math.Pi, rightNumberLiteral.NumberValue, 0.01)
	})

	t.Run("negative", func(t *testing.T) {
		t.Run("left rule error", func(t *testing.T) {
			b := schema.NewSchemaBuilder()
			b.AddField("left", ruleError)
			b.AddField("right", schema.Any())
			schm := b.Build()

			expr := &query.BinaryExpr{
				Left: &query.FieldExpr{
					Field: "left",
					Op:    query.Equal,
					Value: &query.NumberLiteral{NumberValue: 42},
				},
				Op: query.And,
				Right: &query.FieldExpr{
					Field: "right",
					Op:    query.Equal,
					Value: &query.NumberLiteral{NumberValue: math.Pi},
				},
			}

			got, err := expr.Validate(schm)
			require.Error(t, err)

			fieldExpr, isFieldExpr := got.(*query.FieldExpr)
			require.True(t, isFieldExpr)

			numberLiteral, isNumberLiteral := fieldExpr.Value.(*query.NumberLiteral)
			require.True(t, isNumberLiteral)

			require.InDelta(t, math.Pi, numberLiteral.NumberValue, 0.01)
		})

		t.Run("right rule error", func(t *testing.T) {
			b := schema.NewSchemaBuilder()
			b.AddField("left", schema.Any())
			b.AddField("right", ruleError)
			schm := b.Build()

			expr := &query.BinaryExpr{
				Left: &query.FieldExpr{
					Field: "left",
					Op:    query.Equal,
					Value: &query.NumberLiteral{NumberValue: 42},
				},
				Op: query.And,
				Right: &query.FieldExpr{
					Field: "right",
					Op:    query.Equal,
					Value: &query.NumberLiteral{NumberValue: math.Pi},
				},
			}

			got, err := expr.Validate(schm)
			require.Error(t, err)

			fieldExpr, isFieldExpr := got.(*query.FieldExpr)
			require.True(t, isFieldExpr)

			numberLiteral, isNumberLiteral := fieldExpr.Value.(*query.NumberLiteral)
			require.True(t, isNumberLiteral)

			require.InDelta(t, 42.0, numberLiteral.NumberValue, 0.0001)
		})

		t.Run("left and right rule error", func(t *testing.T) {
			b := schema.NewSchemaBuilder()
			b.AddField("left", ruleError)
			b.AddField(
				"right", ruleError)
			schm := b.Build()

			expr := &query.BinaryExpr{
				Left: &query.FieldExpr{
					Field: "left",
					Op:    query.Equal,
					Value: &query.NumberLiteral{NumberValue: 42},
				},
				Right: &query.FieldExpr{
					Field: "right",
					Op:    query.Equal,
					Value: &query.NumberLiteral{NumberValue: math.Pi},
				},
			}

			got, err := expr.Validate(schm)
			require.Error(t, err)
			require.Nil(t, got)
		})

		t.Run("unknown field", func(t *testing.T) {
			schm := schema.Schema{}

			expr := &query.BinaryExpr{
				Left: &query.FieldExpr{
					Field: "left",
					Op:    query.Equal,
					Value: &query.NumberLiteral{NumberValue: 42},
				},
				Right: &query.FieldExpr{
					Field: "right",
					Op:    query.Equal,
					Value: &query.NumberLiteral{NumberValue: math.Pi},
				},
			}

			got, err := expr.Validate(schm)
			require.Error(t, err)
			require.Nil(t, got)
		})
	})
}

func TestNotExpr_Validate(t *testing.T) { //nolint:funlen
	t.Run("positive", func(t *testing.T) {
		b := schema.NewSchemaBuilder()
		b.AddField("field", schema.Any())
		schm := b.Build()

		expr := &query.NotExpr{
			Expr: &query.FieldExpr{
				Field: "field",
				Op:    query.Equal,
				Value: &query.NumberLiteral{NumberValue: 42},
			},
		}

		got, err := expr.Validate(schm)
		require.NoError(t, err)

		notExpr, isNotExpr := got.(*query.NotExpr)
		require.True(t, isNotExpr)

		fieldExpr, isFieldExpr := notExpr.Expr.(*query.FieldExpr)
		require.True(t, isFieldExpr)

		numberLiteral, isNumberLiteral := fieldExpr.Value.(*query.NumberLiteral)
		require.True(t, isNumberLiteral)

		require.InDelta(t, 42.0, numberLiteral.NumberValue, 0.0001)
	})

	t.Run("negative", func(t *testing.T) {
		t.Run("rule error", func(t *testing.T) {
			b := schema.NewSchemaBuilder()
			b.AddField("field", ruleError)
			schm := b.Build()

			expr := &query.NotExpr{
				Expr: &query.FieldExpr{
					Field: "field",
					Op:    query.Equal,
					Value: &query.NumberLiteral{NumberValue: 42},
				},
			}

			got, err := expr.Validate(schm)
			require.Error(t, err)
			require.Nil(t, got)
		})

		t.Run("unknown field", func(t *testing.T) {
			schm := schema.Schema{}

			expr := &query.NotExpr{
				Expr: &query.FieldExpr{
					Field: "field",
					Op:    query.Equal,
					Value: &query.NumberLiteral{NumberValue: 42},
				},
			}

			got, err := expr.Validate(schm)
			require.Error(t, err)
			require.Nil(t, got)
		})
	})
}

func TestFieldExpr_Validate(t *testing.T) { //nolint:funlen
	t.Run("positive", func(t *testing.T) {
		t.Run("primitive value", func(t *testing.T) {
			b := schema.NewSchemaBuilder()
			b.AddField("field", schema.Any())
			schm := b.Build()

			expr := &query.FieldExpr{
				Field: "field",
				Op:    query.Equal,
				Value: &query.NumberLiteral{NumberValue: 42},
			}

			got, err := expr.Validate(schm)
			require.NoError(t, err)

			fieldExpr, isFieldExpr := got.(*query.FieldExpr)
			require.True(t, isFieldExpr)

			numberLiteral, isNumberLiteral := fieldExpr.Value.(*query.NumberLiteral)
			require.True(t, isNumberLiteral)

			require.InDelta(t, 42.0, numberLiteral.NumberValue, 0.0001)
		})

		t.Run("one of", func(t *testing.T) {
			b := schema.NewSchemaBuilder()
			b.AddField("field", schema.Any())
			schm := b.Build()

			expr := &query.FieldExpr{
				Field: "field",
				Op:    query.Equal,
				Value: &query.OneOfExpr{
					Values: []query.Valuer{
						&query.NumberLiteral{NumberValue: 42},
						&query.NumberLiteral{NumberValue: math.Pi},
					},
				},
			}

			got, err := expr.Validate(schm)
			require.NoError(t, err)

			fieldExpr, isFieldExpr := got.(*query.FieldExpr)
			require.True(t, isFieldExpr)

			oneOfExpr, isOneOfExpr := fieldExpr.Value.(*query.OneOfExpr)
			require.True(t, isOneOfExpr)

			number1Literal, isNumber1Literal := oneOfExpr.Values[0].(*query.NumberLiteral)
			require.True(t, isNumber1Literal)

			number2Literal, isNumber2Literal := oneOfExpr.Values[1].(*query.NumberLiteral)
			require.True(t, isNumber2Literal)

			require.InDelta(t, 42.0, number1Literal.NumberValue, 0.0001)
			require.InDelta(t, math.Pi, number2Literal.NumberValue, 0.01)
		})
		t.Run("wildcard field", func(t *testing.T) {
			b := schema.NewSchemaBuilder()
			b.AddRegExpField(`^subject\..*$`, schema.Any())
			schm := b.Build()
			expr := &query.FieldExpr{
				Field: "subject.age",
				Op:    query.Equal,
				Value: &query.NumberLiteral{NumberValue: 42},
			}

			got, err := expr.Validate(schm)
			require.NoError(t, err)

			fieldExpr, isFieldExpr := got.(*query.FieldExpr)
			require.True(t, isFieldExpr)

			numberLiteral, isNumberLiteral := fieldExpr.Value.(*query.NumberLiteral)
			require.True(t, isNumberLiteral)

			require.InDelta(t, 42.0, numberLiteral.NumberValue, 0.0001)
		})
	})

	t.Run("negative", func(t *testing.T) {
		t.Run("primitive value rule error", func(t *testing.T) {
			b := schema.NewSchemaBuilder()
			b.AddField("field", ruleError)
			schm := b.Build()

			expr := &query.FieldExpr{
				Field: "field",
				Op:    query.Equal,
				Value: &query.NumberLiteral{NumberValue: 42},
			}

			got, err := expr.Validate(schm)
			require.Error(t, err)
			require.Nil(t, got)
		})

		t.Run("one of rule error", func(t *testing.T) {
			b := schema.NewSchemaBuilder()
			b.AddField("field", schema.All(
				schema.Is[float64](),
				schema.EqualsOneOf(42.0),
			))
			schm := b.Build()

			expr := &query.FieldExpr{
				Field: "field",
				Op:    query.Equal,
				Value: &query.OneOfExpr{
					Values: []query.Valuer{
						&query.NumberLiteral{NumberValue: 42}, // This will pass
						&query.NumberLiteral{NumberValue: 99}, // This will fail validation
					},
				},
			}

			got, err := expr.Validate(schm)
			require.Error(t, err)

			fieldExpr, isFieldExpr := got.(*query.FieldExpr)
			require.True(t, isFieldExpr)

			oneOfExpr, isOneOfExpr := fieldExpr.Value.(*query.OneOfExpr)
			require.True(t, isOneOfExpr)

			assert.Len(t, oneOfExpr.Values, 1)

			numberLiteral, isNumberLiteral := oneOfExpr.Values[0].(*query.NumberLiteral)
			require.True(t, isNumberLiteral)
			require.InDelta(t, 42.0, numberLiteral.NumberValue, 0.0001)
		})

		t.Run("unknown field", func(t *testing.T) {
			schm := schema.Schema{}

			expr := &query.FieldExpr{
				Field: "field",
				Op:    query.Equal,
				Value: &query.NumberLiteral{NumberValue: 42},
			}

			got, err := expr.Validate(schm)
			require.Error(t, err)
			require.Nil(t, got)
		})
		t.Run("regexp mismatch", func(t *testing.T) {
			b := schema.NewSchemaBuilder()
			b.AddRegExpField(`^subject\..*$`, schema.Any())
			schm := schema.Schema{}

			expr := &query.FieldExpr{
				Field: "subbjjj3ct.age",
				Op:    query.Equal,
				Value: &query.NumberLiteral{NumberValue: 42},
			}

			got, err := expr.Validate(schm)
			require.Error(t, err)
			require.Nil(t, got)
		})
	})
}

func ruleError(schema.Field, any) error {
	return errors.New("rule error")
}
