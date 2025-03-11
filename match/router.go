package match

import (
	"reflect"
	"strings"
)

// ReflectRouter implements Router for struct targets.
// It supports struct tags using the `dumbql` tag name and nested field access using dot notation.
type ReflectRouter struct{}

// Route resolves a field path in the target struct and returns the value.
// It supports nested field access using dot notation (e.g., "address.city").
func (r *ReflectRouter) Route(target any, field string) (any, error) {
	// Handle dot notation for nested fields
	parts := strings.Split(field, ".")

	// Process single field name - common case
	if len(parts) == 1 {
		return r.resolveDirectField(target, field)
	}

	// Handle nested fields traversal
	return r.resolveNestedField(target, parts)
}

// resolveDirectField handles resolving a direct field (no dots in the name)
func (r *ReflectRouter) resolveDirectField(target any, field string) (any, error) {
	v := reflect.ValueOf(target)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, errNotAStruct // Nil pointer, can't resolve field
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, errNotAStruct // Not a struct, can't resolve field
	}

	t := v.Type()
	for i := range t.NumField() {
		f := t.Field(i)

		tag := f.Tag.Get("dumbql")
		if tag == "-" {
			// Field marked with dumbql:"-" is skipped
			return nil, errFieldNotFound
		}

		fname := f.Name
		if tag != "" {
			fname = tag
		}

		if fname == field {
			return v.Field(i).Interface(), nil
		}
	}

	// Field not found
	return nil, errFieldNotFound
}

// resolveNestedField handles traversing nested fields using dot notation
func (r *ReflectRouter) resolveNestedField(target any, path []string) (any, error) {
	current := target

	// Navigate through all segments except the last one
	for i := range len(path) - 1 {
		v := reflect.ValueOf(current)
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				return nil, errNotAStruct // Nil pointer, can't traverse further
			}
			v = v.Elem()
		}

		if v.Kind() != reflect.Struct {
			return nil, errNotAStruct // Not a struct, cannot traverse
		}

		// Find field by name or tag
		t := v.Type()
		found := false

		for j := range t.NumField() {
			f := t.Field(j)

			tag := f.Tag.Get("dumbql")
			if tag == "-" {
				return nil, errFieldNotFound // Field marked with dumbql:"-" is skipped
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
			return nil, errFieldNotFound // Field not found
		}
	}

	// Resolve the final segment
	return r.resolveDirectField(current, path[len(path)-1])
}
