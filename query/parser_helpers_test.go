package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseNumber(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    float64
		wantErr bool
	}{
		{
			name:    "positive integer",
			input:   "42",
			want:    42.0,
			wantErr: false,
		},
		{
			name:    "negative integer",
			input:   "-42",
			want:    -42.0,
			wantErr: false,
		},
		{
			name:    "zero",
			input:   "0",
			want:    0.0,
			wantErr: false,
		},
		{
			name:    "positive float",
			input:   "3.14159",
			want:    3.14159,
			wantErr: false,
		},
		{
			name:    "negative float",
			input:   "-3.14159",
			want:    -3.14159,
			wantErr: false,
		},
		{
			name:    "scientific notation positive",
			input:   "1.23e5",
			want:    123000.0,
			wantErr: false,
		},
		{
			name:    "scientific notation negative",
			input:   "-1.23e-5",
			want:    -0.0000123,
			wantErr: false,
		},
		{
			name:    "invalid number - contains letters",
			input:   "42abc",
			want:    0.0,
			wantErr: true,
		},
		{
			name:    "invalid number - empty string",
			input:   "",
			want:    0.0,
			wantErr: true,
		},
		{
			name:    "invalid number - only decimal point",
			input:   ".",
			want:    0.0,
			wantErr: true,
		},
		{
			name:    "invalid number - multiple decimal points",
			input:   "3.14.159",
			want:    0.0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			curr := &current{
				text: []byte(tt.input),
			}
			
			result, err := parseNumber(curr)
			
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			
			require.NoError(t, err)
			numLiteral, ok := result.(*NumberLiteral)
			require.True(t, ok)
			assert.InDelta(t, tt.want, numLiteral.NumberValue, 0.0001)
		})
	}
}