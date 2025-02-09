package schema

type Field string

// RuleFunc defines a function type for validating a field value and returning an error if validation fails.
type RuleFunc func(field Field, value any) error

// Schema is a set of Field to RuleFunc pairs which defines constraints for the query validation.
type Schema map[Field]RuleFunc

type ValueType interface {
	string | Numeric
}

type Numeric interface {
	float64 | int64
}
