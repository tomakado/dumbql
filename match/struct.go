package match

import (
	"errors"
	"reflect"

	"go.tomakado.io/dumbql/query"
)

// Router is responsible for resolving a field path to a value in the target object.
// It abstracts the field resolution logic from the matcher implementation.
type Router interface {
	// Route resolves a field path in the target object and returns the value.
	// The boolean return value indicates whether the field was successfully resolved.
	Route(target any, field string) (any, error)
}

// StructMatcher is a basic implementation of the Matcher interface for evaluating query expressions against structs.
// It supports struct tags using the `dumbql` tag name, which allows you to specify a custom field name.
type StructMatcher struct {
	router Router
}

func NewStructMatcher(router Router) *StructMatcher {
	return &StructMatcher{
		router: router,
	}
}

func (m *StructMatcher) lazyInit() {
	if m.router == nil {
		m.router = &ReflectRouter{}
	}
}

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
	m.lazyInit()

	fieldValue, err := m.router.Route(target, field)

	switch {
	case op == query.Exists:
		return err == nil && !isZero(fieldValue)
	case err != nil:
		return errors.Is(err, ErrFieldNotFound)
	default:
		return m.MatchValue(fieldValue, value, op)
	}
}

func (m *StructMatcher) MatchValue(target any, value query.Valuer, op query.FieldOperator) bool {
	return value.Match(target, op)
}

func isZero(tv any) bool {
	// FIXME: Need to find a way to do it faster
	v := reflect.ValueOf(tv)
	return !v.IsValid() || reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}
