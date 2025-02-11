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
		rule := schema.Any(schema.Is[int64](), schema.Is[string]())
		require.NoError(t, rule("positive_int", int64(42)))
		require.NoError(t, rule("positive_string", "Hello, world!"))
	})

	t.Run("negative", func(t *testing.T) {
		rule := schema.Any(schema.Is[int64](), schema.Is[string]())
		require.Error(t, rule("negative", 0.75))
	})
}

func TestAll(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		rule := schema.All(schema.Is[int64](), schema.Min[int64](42))
		require.NoError(t, rule("positive", int64(42)))
	})

	t.Run("negative", func(t *testing.T) {
		rule := schema.All(schema.Is[int64](), schema.Min[int64](42))
		require.Error(t, rule("negative", int64(41)))
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
		})

		t.Run("negative", func(t *testing.T) {
			t.Run("wrong value", func(t *testing.T) {
				rule := schema.Min[int64](42)
				require.Error(t, rule("negative", int64(41)))
			})

			t.Run("wrong type", func(t *testing.T) {
				rule := schema.Min[int64](42)
				require.Error(t, rule("negative", 42.42))
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
				require.Error(t, rule("negative", int64(42)))
			})
		})
	})
}

func TestMax(t *testing.T) {
	t.Run("int64", func(t *testing.T) {
		t.Run("positive", func(t *testing.T) {
			rule := schema.Max[int64](42)
			require.NoError(t, rule("positive", int64(42)))
		})

		t.Run("negative", func(t *testing.T) {
			rule := schema.Max[int64](42)
			require.Error(t, rule("negative", int64(43)))
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
	values := []any{"positive", "hello", "world", 42, 0.75}

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
