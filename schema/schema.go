package schema

import (
	"fmt"
	"regexp"
)

type Field string

// RuleFunc defines a function type for validating a field value and returning an error if validation fails.
type RuleFunc func(field Field, value any) error

// Schema is a set of Field to RuleFunc pairs which defines constraints for the query validation.
type Schema struct {
	fields      map[Field]RuleFunc
	regexFields map[*regexp.Regexp]RuleFunc
}

// GetRule looks for a matching rule accross fields and regular expressions.
func (s Schema) GetRule(field Field) (RuleFunc, error) {

	rule, ok := s.fields[field]
	if ok {
		return rule, nil
	}
	for re, rule := range s.regexFields {
		if re.MatchString(string(field)) {
			return rule, nil
		}
	}

	return nil, fmt.Errorf("field %q not found in schema", field)
}

type ValueType interface {
	string | Numeric | bool
}

type Numeric interface {
	float64 | int64
}

type SchemaBuilder struct {
	fields      map[Field]RuleFunc
	regexFields map[*regexp.Regexp]RuleFunc
}

func NewSchemaBuilder() SchemaBuilder {
	return SchemaBuilder{fields: make(map[Field]RuleFunc),
		regexFields: make(map[*regexp.Regexp]RuleFunc)}
}

func (b *SchemaBuilder) AddRegExpField(rex string, rule RuleFunc) *SchemaBuilder {
	re := regexp.MustCompile(rex)
	b.regexFields[re] = rule
	return b
}

func (b *SchemaBuilder) AddField(field Field, rule RuleFunc) *SchemaBuilder {
	b.fields[field] = rule
	return b
}

func (b *SchemaBuilder) Build() Schema {
	return Schema{fields: b.fields, regexFields: b.regexFields}
}
