package query

import (
	"strings"
)

type Matcher interface {
	MatchAnd(target any, left, right Expr) bool
	MatchOr(target any, left, right Expr) bool
	MatchNot(target any, expr Expr) bool
	MatchField(target any, field string, value Valuer, op FieldOperator) bool
	MatchValue(target any, value Valuer, op FieldOperator) bool
}

func (b *BinaryExpr) Match(target any, matcher Matcher) bool {
	switch b.Op {
	case And:
		return matcher.MatchAnd(target, b.Left, b.Right)
	case Or:
		return matcher.MatchOr(target, b.Left, b.Right)
	default:
		return false
	}
}

func (n *NotExpr) Match(target any, matcher Matcher) bool {
	return matcher.MatchNot(target, n.Expr)
}

func (f *FieldExpr) Match(target any, matcher Matcher) bool {
	return matcher.MatchField(target, f.Field.String(), f.Value, f.Op)
}

func (s *StringLiteral) Match(target any, op FieldOperator) bool {
	str, ok := target.(string)
	if !ok {
		return false
	}

	return matchString(str, s.StringValue, op)
}

func (n *NumberLiteral) Match(target any, op FieldOperator) bool {
	// Convert target to float64 regardless of its type
	var targetFloat float64
	var ok bool

	switch v := target.(type) {
	case float64:
		targetFloat = v
		ok = true
	case float32:
		targetFloat = float64(v)
		ok = true
	case int:
		targetFloat = float64(v)
		ok = true
	case int8:
		targetFloat = float64(v)
		ok = true
	case int16:
		targetFloat = float64(v)
		ok = true
	case int32:
		targetFloat = float64(v)
		ok = true
	case int64:
		targetFloat = float64(v)
		ok = true
	case uint:
		targetFloat = float64(v)
		ok = true
	case uint8:
		targetFloat = float64(v)
		ok = true
	case uint16:
		targetFloat = float64(v)
		ok = true
	case uint32:
		targetFloat = float64(v)
		ok = true
	case uint64:
		targetFloat = float64(v)
		ok = true
	}

	if !ok {
		return false
	}

	return matchNum(targetFloat, n.NumberValue, op)
}

func (i Identifier) Match(target any, op FieldOperator) bool {
	str, ok := target.(string)
	if !ok {
		return false
	}

	return matchString(str, i.String(), op)
}

func (o *OneOfExpr) Match(target any, op FieldOperator) bool {
	switch op { //nolint:exhaustive
	case Equal, Like:
		for _, v := range o.Values {
			if v.Match(target, op) {
				return true
			}
		}

		return false

	default:
		return false
	}
}

func matchString(a, b string, op FieldOperator) bool {
	switch op { //nolint:exhaustive
	case Equal:
		return a == b
	case NotEqual:
		return a != b
	case Like:
		return strings.Contains(a, b)
	default:
		return false
	}
}

func matchNum(a, b float64, op FieldOperator) bool {
	switch op { //nolint:exhaustive
	case Equal:
		return a == b
	case NotEqual:
		return a != b
	case GreaterThan:
		return a > b
	case GreaterThanOrEqual:
		return a >= b
	case LessThan:
		return a < b
	case LessThanOrEqual:
		return a <= b
	default:
		return false
	}
}
