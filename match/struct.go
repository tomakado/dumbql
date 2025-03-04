package match

import (
	"reflect"
	"strings"

	"go.tomakado.io/dumbql/query"
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

// MatchField matches a field in the target struct using the provided value and operator.
// It supports struct tags using the `dumbql` tag name and nested field access using dot notation.
// For example: "address.city" to access the city field in the address struct.
func (m *StructMatcher) MatchField(target any, field string, value query.Valuer, op query.FieldOperator) bool {
	// Handle dot notation for nested fields
	parts := strings.Split(field, ".")

	// Process single field name - common case
	if len(parts) == 1 {
		return m.matchDirectField(target, field, value, op)
	}

	// Handle nested fields traversal
	return m.matchNestedField(target, parts, value, op)
}

// matchDirectField handles matching a direct field (no dots in the name)
func (m *StructMatcher) matchDirectField(target any, field string, value query.Valuer, op query.FieldOperator) bool {
	v := reflect.ValueOf(target)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return false
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		tag := f.Tag.Get("dumbql")
		if tag == "-" {
			// Field marked with dumbql:"-" always match (in other words does not affect the result)
			return true
		}

		fname := f.Name
		if tag != "" {
			fname = tag
		}

		if fname == field {
			return m.MatchValue(v.Field(i).Interface(), value, op)
		}
	}

	// If field not found, return true (same behavior as original)
	return true
}

// matchNestedField handles traversing nested fields using dot notation
func (m *StructMatcher) matchNestedField(target any, path []string, value query.Valuer, op query.FieldOperator) bool {
	current := target

	// Navigate through all segments except the last one
	for i := 0; i < len(path)-1; i++ {
		v := reflect.ValueOf(current)
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				return true // Nil pointer, can't traverse further
			}
			v = v.Elem()
		}

		if v.Kind() != reflect.Struct {
			return false // Not a struct, cannot traverse
		}

		// Find field by name or tag
		t := v.Type()
		found := false

		for j := 0; j < t.NumField(); j++ {
			f := t.Field(j)

			tag := f.Tag.Get("dumbql")
			if tag == "-" {
				return true // Field marked with dumbql:"-" always matches
			}

			fname := f.Name
			if tag != "" {
				fname = tag
			}

			if fname == path[i] {
				current = v.Field(j).Interface()
				found = true
				break
			}
		}

		if !found {
			return true // Field not found, behave same as current implementation
		}
	}

	// Match the final segment
	return m.matchDirectField(current, path[len(path)-1], value, op)
}

func (m *StructMatcher) MatchValue(target any, value query.Valuer, op query.FieldOperator) bool {
	return value.Match(target, op)
}
