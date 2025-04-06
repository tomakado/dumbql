package match

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_isZero(t *testing.T) { //nolint:funlen
	tests := []struct {
		name   string
		value  any
		expect bool
	}{
		{
			name:   "nil value",
			value:  nil,
			expect: true,
		},
		{
			name:   "zero int",
			value:  0,
			expect: true,
		},
		{
			name:   "non-zero int",
			value:  42,
			expect: false,
		},
		{
			name:   "zero int64",
			value:  int64(0),
			expect: true,
		},
		{
			name:   "non-zero int64",
			value:  int64(42),
			expect: false,
		},
		{
			name:   "zero float64",
			value:  0.0,
			expect: true,
		},
		{
			name:   "non-zero float64",
			value:  3.14,
			expect: false,
		},
		{
			name:   "empty string",
			value:  "",
			expect: true,
		},
		{
			name:   "non-empty string",
			value:  "hello",
			expect: false,
		},
		{
			name:   "false boolean",
			value:  false,
			expect: true,
		},
		{
			name:   "true boolean",
			value:  true,
			expect: false,
		},
		{
			name:   "empty slice",
			value:  []int{},
			expect: false, // Empty slice is not considered zero by reflect.DeepEqual
		},
		{
			name:   "non-empty slice",
			value:  []int{1, 2, 3},
			expect: false,
		},
		{
			name:   "empty map",
			value:  map[string]int{},
			expect: false, // Empty map is not considered zero by reflect.DeepEqual
		},
		{
			name:   "non-empty map",
			value:  map[string]int{"a": 1},
			expect: false,
		},
		{
			name:   "zero struct",
			value:  struct{}{},
			expect: true,
		},
		{
			name:   "zero time",
			value:  time.Time{},
			expect: true,
		},
		{
			name:   "non-zero time",
			value:  time.Now(),
			expect: false,
		},
		{
			name: "struct with zero fields",
			value: struct {
				Name string
				Age  int
			}{},
			expect: true,
		},
		{
			name: "struct with non-zero fields",
			value: struct {
				Name string
				Age  int
			}{
				Name: "John",
				Age:  30,
			},
			expect: false,
		},
		{
			name:   "nil pointer",
			value:  (*int)(nil),
			expect: true,
		},
		{
			name: "pointer to zero value",
			value: func() any {
				i := 0
				return &i
			}(),
			expect: false, // Not zero because it's a valid pointer
		},
		{
			name: "pointer to non-zero value",
			value: func() any {
				i := 42
				return &i
			}(),
			expect: false,
		},
		{
			name:   "invalid reflect value",
			value:  make(chan int), // Channels are not comparable with DeepEqual
			expect: false,          // Not zero because it's a valid channel
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isZero(tt.value)
			assert.Equal(t, tt.expect, result, "isZero(%v) = %v, want %v", tt.value, result, tt.expect)
		})
	}
}
