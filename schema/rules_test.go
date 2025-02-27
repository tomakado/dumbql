package schema_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.tomakado.io/dumbql/schema"
)

func TestAny(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		rule := schema.Any(schema.Is[float64](), schema.Is[string]())
		require.NoError(t, rule("positive_int", float64(42)))
		require.NoError(t, rule("positive_string", "Hello, world!"))
	})

	t.Run("negative", func(t *testing.T) {
		rule := schema.Any(schema.Is[float64](), schema.Is[string]())
		require.Error(t, rule("negative", true))
	})
}

func TestAll(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		rule := schema.All(schema.Is[float64](), schema.Min[float64](42))
		require.NoError(t, rule("positive", float64(42)))
	})

	t.Run("negative", func(t *testing.T) {
		rule := schema.All(schema.Is[float64](), schema.Min[float64](42))
		require.Error(t, rule("negative", float64(41)))
	})
}

func TestInRange(t *testing.T) {
	t.Run("int64", func(t *testing.T) {
		t.Run("positive", func(t *testing.T) {
			rule := schema.InRange[int64](5, 10)
			require.NoError(t, rule("positive", int64(7)))
		})

		t.Run("negative", func(t *testing.T) {
			rule := schema.InRange[int64](5, 10)
			require.Error(t, rule("negative", int64(42)))
		})
	})

	t.Run("float64", func(t *testing.T) {
		t.Run("positive", func(t *testing.T) {
			rule := schema.InRange[float64](5.0, 10.0)
			require.NoError(t, rule("positive", 7.5))
		})

		t.Run("negative", func(t *testing.T) {
			rule := schema.InRange[float64](5.0, 10.0)
			require.Error(t, rule("negative", 42.0))
		})
	})
}

func TestMin(t *testing.T) {
	t.Run("int64", func(t *testing.T) {
		t.Run("positive", func(t *testing.T) {
			rule := schema.Min[int64](42)
			require.NoError(t, rule("positive", int64(42)))
			// Test our new feature - float64 values should work with int64 min
			require.NoError(t, rule("positive_float", 42.5))
		})

		t.Run("negative", func(t *testing.T) {
			t.Run("wrong value", func(t *testing.T) {
				rule := schema.Min[int64](42)
				require.Error(t, rule("negative", int64(41)))
				// Even with our new feature, lower values should fail
				require.Error(t, rule("negative_float", 41.5))
			})

			t.Run("wrong type", func(t *testing.T) {
				rule := schema.Min[int64](42)
				// With our improved implementation, float64 values work with int64 minimums
				require.NoError(t, rule("float_works", 42.42))
				// But other non-numeric types should still fail
				require.Error(t, rule("string_fails", "42"))
			})
		})
	})

	t.Run("float64", func(t *testing.T) {
		t.Run("positive", func(t *testing.T) {
			rule := schema.Min[float64](42.42)
			require.NoError(t, rule("positive", 42.42))
		})

		t.Run("negative", func(t *testing.T) {
			t.Run("wrong value", func(t *testing.T) {
				rule := schema.Min[float64](42.42)
				require.Error(t, rule("negative", 42.41))
			})

			t.Run("wrong type", func(t *testing.T) {
				rule := schema.Min[float64](42.42)
				// With our improved implementation, int64 values work with float64 minimums
				require.NoError(t, rule("integer_works", int64(43)))
				// But other non-numeric types should still fail
				require.Error(t, rule("string_fails", "42.42"))
			})
		})
	})
}

func TestMax(t *testing.T) {
	t.Run("int64", func(t *testing.T) {
		t.Run("positive", func(t *testing.T) {
			rule := schema.Max[int64](42)
			require.NoError(t, rule("positive", int64(42)))
			// Test our new feature - float64 values should work with int64 max
			require.NoError(t, rule("positive_float", 41.5))
		})

		t.Run("negative", func(t *testing.T) {
			rule := schema.Max[int64](42)
			require.Error(t, rule("negative", int64(43)))
			// Even with our new feature, higher values should fail
			require.Error(t, rule("negative_float", 42.5))
		})
	})

	t.Run("float64", func(t *testing.T) {
		t.Run("positive", func(t *testing.T) {
			rule := schema.Max[float64](42.42)
			require.NoError(t, rule("positive", 42.42))
		})

		t.Run("negative", func(t *testing.T) {
			rule := schema.Max[float64](42.42)
			require.Error(t, rule("negative", 42.43))
		})
	})
}

func TestLenInRange(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		rule := schema.LenInRange(5, 10)
		require.NoError(t, rule("positive", "hello"))
	})

	t.Run("negative", func(t *testing.T) {
		t.Run("wrong len", func(t *testing.T) {
			rule := schema.LenInRange(5, 10)
			require.Error(t, rule("negative", "hi"))
		})

		t.Run("wrong type", func(t *testing.T) {
			rule := schema.LenInRange(5, 10)
			require.Error(t, rule("negative", 42))
		})
	})
}

func TestMinLen(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		rule := schema.MinLen(5)
		require.NoError(t, rule("positive", "hello, world!"))
	})

	t.Run("negative", func(t *testing.T) {
		t.Run("wrong len", func(t *testing.T) {
			rule := schema.MinLen(5)
			require.Error(t, rule("negative", "hi"))
		})

		t.Run("wrong type", func(t *testing.T) {
			rule := schema.MinLen(5)
			require.Error(t, rule("negative", 42))
		})
	})
}

func TestMaxLen(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		rule := schema.MaxLen(5)
		require.NoError(t, rule("positive", "hello"))
	})

	t.Run("negative", func(t *testing.T) {
		t.Run("wrong len", func(t *testing.T) {
			rule := schema.MaxLen(5)
			require.Error(t, rule("negative_len", "hello, world!"))
		})

		t.Run("wrong type", func(t *testing.T) {
			rule := schema.MaxLen(5)
			require.Error(t, rule("negative_type", 42))
		})
	})
}

func TestIs(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		t.Run("positive", func(t *testing.T) {
			rule := schema.Is[string]()
			require.NoError(t, rule("string_positive", "Hello, world!"))
		})

		t.Run("negative", func(t *testing.T) {
			rule := schema.Is[string]()
			require.Error(t, rule("string_negative", 42))
		})
	})

	t.Run("int64", func(t *testing.T) {
		t.Run("positive", func(t *testing.T) {
			rule := schema.Is[int64]()
			require.NoError(t, rule("int64_positive", int64(42)))
		})

		t.Run("negative", func(t *testing.T) {
			rule := schema.Is[int64]()
			require.Error(t, rule("int64_negative", "Hello, world!"))
		})
	})

	t.Run("float64", func(t *testing.T) {
		t.Run("positive", func(t *testing.T) {
			rule := schema.Is[float64]()
			require.NoError(t, rule("float64_positive", 42.42))
		})

		t.Run("negative", func(t *testing.T) {
			rule := schema.Is[float64]()
			require.Error(t, rule("float64_negative", "Hello, world!"))
		})
	})
}

func TestEqualsOneOf(t *testing.T) {
	values := []any{"positive", "hello", "world", 42.0, 0.75}

	t.Run("positive", func(t *testing.T) {
		rule := schema.EqualsOneOf(values...)

		for i, value := range values {
			field := schema.Field(fmt.Sprintf("positive_%d", i))
			assert.NoError(t, rule(field, value))
		}
	})

	t.Run("negative", func(t *testing.T) {
		rule := schema.EqualsOneOf(values...)
		require.Error(t, rule("negative", math.Pi))
	})
}