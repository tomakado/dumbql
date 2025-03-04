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
// It supports struct tags using the `dumbql` tag name, which allows you to specify a custom field name.
// It also supports nested field access using dot notation (e.g., "top_level.second_level.third_level").
func (m *StructMatcher) MatchField(target any, field string, value query.Valuer, op query.FieldOperator) bool {
	// Check if this is a nested field access (contains dots)
	segments := strings.Split(field, ".")
	if len(segments) == 1 {
		// Use existing non-nested implementation
		return m.matchSingleField(target, field, value, op)
	}

	// Handle nested field traversal
	return m.matchNestedField(target, segments, value, op)
}

// matchSingleField is the implementation for matching a single field
func (m *StructMatcher) matchSingleField(target any, field string, value query.Valuer, op query.FieldOperator) bool {
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

// matchNestedField traverses nested structs to find and match the target field
func (m *StructMatcher) matchNestedField(
	target any,
	segments []string,
	value query.Valuer,
	op query.FieldOperator,
) bool {
	currentTarget := target

	// Traverse the struct hierarchy for all segments except the last one
	for i := 0; i < len(segments)-1; i++ {
		currentSegment := segments[i]

		// Find the field in the current struct
		v := reflect.ValueOf(currentTarget)
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				return true // Nil pointer, can't traverse further
			}
			v = v.Elem()
		}

		if v.Kind() != reflect.Struct {
			return false // Not a struct, cannot traverse
		}

		// Find the field by name or tag
		field, found := m.findField(v, currentSegment)
		if !found {
			return true // Field not found, behave same as current implementation
		}

		// Move to the next level in the hierarchy
		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				return true // Nil pointer, can't traverse further
			}
			currentTarget = field.Elem().Interface()
		} else {
			currentTarget = field.Interface()
		}
	}

	// Match the final field
	lastSegment := segments[len(segments)-1]
	return m.matchSingleField(currentTarget, lastSegment, value, op)
}

// findField looks for a field by name or tag in a struct value
func (m *StructMatcher) findField(v reflect.Value, fieldName string) (reflect.Value, bool) {
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		tag := f.Tag.Get("dumbql")
		if tag == "-" {
			continue // Skip fields marked with dumbql:"-"
		}

		fname := f.Name
		if tag != "" {
			fname = tag
		}

		if fname == fieldName {
			return v.Field(i), true
		}
	}

	return reflect.Value{}, false
}

func (m *StructMatcher) MatchValue(target any, value query.Valuer, op query.FieldOperator) bool {
	return value.Match(target, op)
}
