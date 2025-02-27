package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToFloat64(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    float64
		wantOk  bool
	}{
		{
			name:    "float64",
			input:   float64(42.5),
			want:    42.5,
			wantOk:  true,
		},
		{
			name:    "float32",
			input:   float32(42.5),
			want:    42.5,
			wantOk:  true,
		},
		{
			name:    "int",
			input:   int(42),
			want:    42.0,
			wantOk:  true,
		},
		{
			name:    "int8",
			input:   int8(42),
			want:    42.0,
			wantOk:  true,
		},
		{
			name:    "int16",
			input:   int16(42),
			want:    42.0,
			wantOk:  true,
		},
		{
			name:    "int32",
			input:   int32(42),
			want:    42.0,
			wantOk:  true,
		},
		{
			name:    "int64",
			input:   int64(42),
			want:    42.0,
			wantOk:  true,
		},
		{
			name:    "uint",
			input:   uint(42),
			want:    42.0,
			wantOk:  true,
		},
		{
			name:    "uint8",
			input:   uint8(42),
			want:    42.0,
			wantOk:  true,
		},
		{
			name:    "uint16",
			input:   uint16(42),
			want:    42.0,
			wantOk:  true,
		},
		{
			name:    "uint32",
			input:   uint32(42),
			want:    42.0,
			wantOk:  true,
		},
		{
			name:    "uint64",
			input:   uint64(42),
			want:    42.0,
			wantOk:  true,
		},
		{
			name:    "string - not convertible",
			input:   "42",
			want:    0.0,
			wantOk:  false,
		},
		{
			name:    "bool - not convertible",
			input:   true,
			want:    0.0,
			wantOk:  false,
		},
		{
			name:    "nil - not convertible",
			input:   nil,
			want:    0.0,
			wantOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := convertToFloat64(tt.input)
			assert.Equal(t, tt.wantOk, ok)
			if tt.wantOk {
				assert.InDelta(t, tt.want, got, 0.0001)
			}
		})
	}
}

func TestNumberLiteral_Match_TypeConversion(t *testing.T) {
	tests := []struct {
		name     string
		literal  *NumberLiteral
		target   any
		operator FieldOperator
		want     bool
	}{
		{
			name:     "float64 equal int64",
			literal:  &NumberLiteral{NumberValue: 42.0},
			target:   int64(42),
			operator: Equal,
			want:     true,
		},
		{
			name:     "float64 not equal int64",
			literal:  &NumberLiteral{NumberValue: 42.0},
			target:   int64(43),
			operator: Equal,
			want:     false,
		},
		{
			name:     "float64 greater than int64",
			literal:  &NumberLiteral{NumberValue: 43.0},
			target:   int64(42),
			operator: GreaterThan,
			want:     false, // target is not greater than literal
		},
		{
			name:     "float64 less than int64",
			literal:  &NumberLiteral{NumberValue: 41.0},
			target:   int64(42),
			operator: LessThan,
			want:     false, // target is not less than literal
		},
		{
			name:     "float64 greater than or equal int64 (equal)",
			literal:  &NumberLiteral{NumberValue: 42.0},
			target:   int64(42),
			operator: GreaterThanOrEqual,
			want:     true,
		},
		{
			name:     "float64 less than or equal int64 (equal)",
			literal:  &NumberLiteral{NumberValue: 42.0},
			target:   int64(42),
			operator: LessThanOrEqual,
			want:     true,
		},
		{
			name:     "float64 with int",
			literal:  &NumberLiteral{NumberValue: 42.0},
			target:   42,
			operator: Equal,
			want:     true,
		},
		{
			name:     "float64 with float32",
			literal:  &NumberLiteral{NumberValue: 42.0},
			target:   float32(42.0),
			operator: Equal,
			want:     true,
		},
		{
			name:     "float64 with uint64",
			literal:  &NumberLiteral{NumberValue: 42.0},
			target:   uint64(42),
			operator: Equal,
			want:     true,
		},
		{
			name:     "float64 with non-numeric type",
			literal:  &NumberLiteral{NumberValue: 42.0},
			target:   "42",
			operator: Equal,
			want:     false,
		},
		{
			name:     "float64 with invalid operator",
			literal:  &NumberLiteral{NumberValue: 42.0},
			target:   int64(42),
			operator: Like,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.literal.Match(tt.target, tt.operator)
			assert.Equal(t, tt.want, got)
		})
	}
}