package query

import (
	"fmt"

	"github.com/defer-panic/dumbql/schema"
	"go.uber.org/multierr"
)

// Validate checks if the binary expression is valid against the schema.
// If either the left or right expression is invalid, the only valid expression is returned.
func (b *BinaryExpr) Validate(schema schema.Schema) (Expr, error) {
	left, err := b.Left.Validate(schema)

	right, rightErr := b.Right.Validate(schema)
	if rightErr != nil {
		err = multierr.Append(err, rightErr)
		if right == nil {
			return left, err
		}
	}

	if left == nil {
		return right, err
	}

	return &BinaryExpr{
		Left:  left,
		Op:    b.Op,
		Right: right,
	}, err
}

// Validate checks if the not expression is valid against the schema.
func (n *NotExpr) Validate(schema schema.Schema) (Expr, error) {
	expr, err := n.Expr.Validate(schema)
	if err != nil {
		if expr == nil {
			return nil, err
		}

		return expr, err
	}

	return n, nil
}

// Validate checks if the field expression is valid against the corresponding schema rule.
func (f *FieldExpr) Validate(schm schema.Schema) (Expr, error) {
	field := schema.Field(f.Field)

	rule, ok := schm[field]
	if !ok {
		return nil, fmt.Errorf("field %q not found in schema", f.Field)
	}

	oneOf, isOneOf := f.Value.(*OneOfExpr)
	if !isOneOf {
		if err := rule(field, f.Value.Value()); err != nil {
			return nil, err
		}
		return f, nil
	}

	var (
		values = make([]Valuer, 0, len(oneOf.Values))
		err    error
	)

	for _, v := range oneOf.Values {
		if ruleErr := rule(field, v.Value()); ruleErr != nil {
			err = multierr.Append(err, ruleErr)
			continue
		}
		values = append(values, v)
	}

	return &FieldExpr{
		Field: f.Field,
		Op:    f.Op,
		Value: &OneOfExpr{Values: values},
	}, err
}
