package query_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.tomakado.io/dumbql/query"
)

func TestNumberParsing(t *testing.T) { //nolint:funlen
	tests := []struct {
		name    string
		input   string
		want    float64
		wantErr bool
	}{
		{
			name:    "positive integer",
			input:   "field:42",
			want:    42.0,
			wantErr: false,
		},
		{
			name:    "negative integer",
			input:   "field:-42",
			want:    -42.0,
			wantErr: false,
		},
		{
			name:    "zero",
			input:   "field:0",
			want:    0.0,
			wantErr: false,
		},
		{
			name:    "positive float",
			input:   "field:3.14159",
			want:    3.14159,
			wantErr: false,
		},
		{
			name:    "negative float",
			input:   "field:-3.14159",
			want:    -3.14159,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Since we can't directly test parseNumber in the query_test package,
			// we'll use the Parse function to parse a simple query with a number
			if tt.wantErr {
				_, err := query.Parse("test", []byte(tt.input))
				assert.Error(t, err)
				return
			}

			result, err := query.Parse("test", []byte(tt.input))
			require.NoError(t, err)

			fieldExpr, ok := result.(*query.FieldExpr)
			require.True(t, ok, "Expected *query.FieldExpr, got %T", result)

			numLiteral, ok := fieldExpr.Value.(*query.NumberLiteral)
			require.True(t, ok, "Expected *query.NumberLiteral, got %T", fieldExpr.Value)
			assert.InDelta(t, tt.want, numLiteral.NumberValue, 0.0001)
		})
	}
}
