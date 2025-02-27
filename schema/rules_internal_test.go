package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInRange_CrossTypeComparisons(t *testing.T) {
	tests := []struct {
		name    string
		min     any
		max     any
		value   any
		wantErr bool
	}{
		{
			name:    "float64 value with int64 bounds - within range",
			min:     int64(10),
			max:     int64(20),
			value:   float64(15.5),
			wantErr: false,
		},
		{
			name:    "float64 value with int64 bounds - below range",
			min:     int64(10),
			max:     int64(20),
			value:   float64(5.5),
			wantErr: true,
		},
		{
			name:    "float64 value with int64 bounds - above range",
			min:     int64(10),
			max:     int64(20),
			value:   float64(25.5),
			wantErr: true,
		},
		{
			name:    "float64 value with int64 bounds - equal to min",
			min:     int64(10),
			max:     int64(20),
			value:   float64(10.0),
			wantErr: false,
		},
		{
			name:    "float64 value with int64 bounds - equal to max",
			min:     int64(10),
			max:     int64(20),
			value:   float64(20.0),
			wantErr: false,
		},
		{
			name:    "int64 value with float64 bounds - within range",
			min:     float64(10.5),
			max:     float64(20.5),
			value:   int64(15),
			wantErr: false,
		},
		{
			name:    "int64 value with float64 bounds - below range",
			min:     float64(10.5),
			max:     float64(20.5),
			value:   int64(5),
			wantErr: true,
		},
		{
			name:    "int64 value with float64 bounds - above range",
			min:     float64(10.5),
			max:     float64(20.5),
			value:   int64(25),
			wantErr: true,
		},
		{
			name:    "int64 value with float64 bounds - between min and integer",
			min:     float64(10.5),
			max:     float64(20.5),
			value:   int64(10),
			wantErr: true,
		},
		{
			name:    "int64 value with float64 bounds - between integer and max",
			min:     float64(10.5),
			max:     float64(20.5),
			value:   int64(21),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rule RuleFunc
			
			switch min := tt.min.(type) {
			case int64:
				rule = InRange[int64](min, tt.max.(int64))
			case float64:
				rule = InRange[float64](min, tt.max.(float64))
			}
			
			err := rule("test_field", tt.value)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMin_CrossTypeComparisons(t *testing.T) {
	tests := []struct {
		name    string
		min     any
		value   any
		wantErr bool
	}{
		{
			name:    "float64 value with int64 min - above min",
			min:     int64(10),
			value:   float64(15.5),
			wantErr: false,
		},
		{
			name:    "float64 value with int64 min - below min",
			min:     int64(10),
			value:   float64(5.5),
			wantErr: true,
		},
		{
			name:    "float64 value with int64 min - equal to min",
			min:     int64(10),
			value:   float64(10.0),
			wantErr: false,
		},
		{
			name:    "int64 value with float64 min - above min",
			min:     float64(10.5),
			value:   int64(15),
			wantErr: false,
		},
		{
			name:    "int64 value with float64 min - below min",
			min:     float64(10.5),
			value:   int64(5),
			wantErr: true,
		},
		{
			name:    "int64 value with float64 min - between min and integer",
			min:     float64(10.5),
			value:   int64(10),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rule RuleFunc
			
			switch min := tt.min.(type) {
			case int64:
				rule = Min[int64](min)
			case float64:
				rule = Min[float64](min)
			}
			
			err := rule("test_field", tt.value)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMax_CrossTypeComparisons(t *testing.T) {
	tests := []struct {
		name    string
		max     any
		value   any
		wantErr bool
	}{
		{
			name:    "float64 value with int64 max - below max",
			max:     int64(20),
			value:   float64(15.5),
			wantErr: false,
		},
		{
			name:    "float64 value with int64 max - above max",
			max:     int64(20),
			value:   float64(25.5),
			wantErr: true,
		},
		{
			name:    "float64 value with int64 max - equal to max",
			max:     int64(20),
			value:   float64(20.0),
			wantErr: false,
		},
		{
			name:    "int64 value with float64 max - below max",
			max:     float64(20.5),
			value:   int64(15),
			wantErr: false,
		},
		{
			name:    "int64 value with float64 max - above max",
			max:     float64(20.5),
			value:   int64(25),
			wantErr: true,
		},
		{
			name:    "int64 value with float64 max - between integer and max",
			max:     float64(20.5),
			value:   int64(21),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rule RuleFunc
			
			switch max := tt.max.(type) {
			case int64:
				rule = Max[int64](max)
			case float64:
				rule = Max[float64](max)
			}
			
			err := rule("test_field", tt.value)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}