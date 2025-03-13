package match

import (
	"reflect"
)

// For high-performance applications, consider using the code generator
// (cmd/dumbqlgen) to create a type-specific Router implementation instead
// of using ReflectRouter. The generated router avoids reflection at runtime,
// providing better performance, especially in hot paths.
//
// Example usage:
// //go:generate dumbqlgen -type User -package .

// ReflectRouter implements Router for struct targets.
// It supports struct tags using the `dumbql` tag name and nested field access using dot notation.
type ReflectRouter struct{}

// Route resolves a field path in the target struct and returns the value.
// It supports nested field access using dot notation (e.g., "address.city").
func (r *ReflectRouter) Route(target any, field string) (any, error) {
	var (
		cursor = target
		err    error
	)

	for field := range Path(field) {
		if field == "" {
			return nil, ErrFieldNotFound
		}

		cursor, err = r.resolveField(cursor, field)
		if err != nil {
			return nil, err
		}
	}

	return cursor, nil
}

func (r *ReflectRouter) resolveField(target any, field string) (any, error) {
	v := reflect.ValueOf(target)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, ErrNotAStruct // Nil pointer, can't resolve field
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, ErrNotAStruct // Not a struct, can't resolve field
	}

	t := v.Type()
	for i := range t.NumField() {
		f := t.Field(i)

		tag := f.Tag.Get("dumbql")
		if tag == "-" {
			// Field marked with dumbql:"-" is skipped
			return nil, ErrFieldNotFound
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
	return nil, ErrFieldNotFound
}
