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

func (i *IntegerLiteral) Match(target any, op FieldOperator) bool {
	intVal, ok := target.(int64)
	if !ok {
		return false
	}

	return matchNum(intVal, i.IntegerValue, op)
}

func (n *NumberLiteral) Match(target any, op FieldOperator) bool {
	floatVal, ok := target.(float64)
	if !ok {
		return false
	}

	return matchNum(floatVal, n.NumberValue, op)
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

func matchNum[T int64 | float64](a, b T, op FieldOperator) bool {
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
