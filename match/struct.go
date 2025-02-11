package match

import (
	"reflect"

	"github.com/defer-panic/dumbql/query"
)

// StructMatcher is a basic implementation of the Matcher interface for evaluating query expressions against structs.
// It supports struct tags using the `dumbql` tag name, which allows you to specify a custom field name.
type StructMatcher struct{}

func (m *StructMatcher) MatchAnd(target any, left, right query.Expr) bool {
	return left.Match(target, m) && right.Match(target, m)
}

func (m *StructMatcher) MatchOr(target any, left, right query.Expr) bool {
	return left.Match(target, m) || right.Match(target, m)
}

func (m *StructMatcher) MatchNot(target any, expr query.Expr) bool {
	return !expr.Match(target, m)
}

// MatchField matches a field in the target struct using the provided value and operator. It supports struct tags using
// the `dumbql` tag name, which allows you to specify a custom field name. If struct tag is not provided, it will use
// the field name as is.
func (m *StructMatcher) MatchField(target any, field string, value query.Valuer, op query.FieldOperator) bool {
	t := reflect.TypeOf(target)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		tag := f.Tag.Get("dumbql")
		if tag == "-" {
			return true // Field marked with dumbql:"-" always match (in other words does not affect the result)
		}

		fname := f.Name
		if tag != "" {
			fname = tag
		}

		if fname == field {
			v := reflect.ValueOf(target)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}

			return m.MatchValue(v.Field(i).Interface(), value, op)
		}
	}

	return true
}

func (m *StructMatcher) MatchValue(target any, value query.Valuer, op query.FieldOperator) bool {
	return value.Match(target, op)
}
