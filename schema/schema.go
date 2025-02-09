package schema

type Field string

type RuleFunc func(field Field, value any) error

type Schema map[Field]RuleFunc

type ValueType interface {
	string | Numeric
}

type Numeric interface {
	float64 | int64
}
