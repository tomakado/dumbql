package query_test

import (
	"testing"

	"github.com/defer-panic/dumbql/query"
	"github.com/stretchr/testify/require"

	sq "github.com/Masterminds/squirrel"
)

func TestToSql(t *testing.T) {
	tests := []struct {
		input    string
		want     string
		wantArgs []any
	}{
		{
			// Simple equality using colon (converted to "=")
			input:    "status:200",
			want:     "SELECT * FROM dummy_table WHERE status = ?",
			wantArgs: []any{int64(200)},
		},
		{
			// Floating-point comparison with ">"
			input:    "eps>0.003",
			want:     "SELECT * FROM dummy_table WHERE eps > ?",
			wantArgs: []any{0.003},
		},
		{
			// Boolean AND between two conditions.
			input:    "status:200 and eps < 0.003",
			want:     "SELECT * FROM dummy_table WHERE (status = ? AND eps < ?)",
			wantArgs: []any{int64(200), 0.003},
		},
		{
			// Boolean OR between two conditions.
			input:    "status:200 or eps < 0.003",
			want:     "SELECT * FROM dummy_table WHERE (status = ? OR eps < ?)",
			wantArgs: []any{int64(200), 0.003},
		},
		{
			// NOT operator applied to a field expression.
			input:    "not status:200",
			want:     "SELECT * FROM dummy_table WHERE NOT status = ?",
			wantArgs: []any{int64(200)},
		},
		{
			// Parenthesized expression.
			input:    "(status:200 and eps<0.003)",
			want:     "SELECT * FROM dummy_table WHERE (status = ? AND eps < ?)",
			wantArgs: []any{int64(200), 0.003},
		},
		{
			// Array literal conversion (using IN).
			input:    "req.fields.ext:[\"jpg\", \"png\"]",
			want:     "SELECT * FROM dummy_table WHERE req.fields.ext IN (?,?)",
			wantArgs: []any{"jpg", "png"},
		},
		{
			// Complex expression combining AND and a parenthesized array literal.
			input:    "status:200 and eps<0.003 and (req.fields.ext:[\"jpg\", \"png\"])",
			want:     "SELECT * FROM dummy_table WHERE ((status = ? AND eps < ?) AND req.fields.ext IN (?,?))",
			wantArgs: []any{int64(200), 0.003, "jpg", "png"},
		},
		{
			// Greater than or equal operator.
			input:    "cmp>=100",
			want:     "SELECT * FROM dummy_table WHERE cmp >= ?",
			wantArgs: []any{int64(100)},
		},
		{
			// Less than or equal operator.
			input:    "price<=50",
			want:     "SELECT * FROM dummy_table WHERE price <= ?",
			wantArgs: []any{int64(50)},
		},
		{
			// Nested NOT with a parenthesized expression.
			input:    "not (status:200 and eps < 0.003)",
			want:     "SELECT * FROM dummy_table WHERE NOT (status = ? AND eps < ?)",
			wantArgs: []any{int64(200), 0.003},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			ast, err := query.Parse("test", []byte(test.input))
			require.NoError(t, err)
			require.NotNil(t, ast)

			expr, ok := ast.(query.Expr)
			require.True(t, ok)

			got, gotArgs, err := sq.Select("*").From("dummy_table").Where(expr).ToSql()
			require.NoError(t, err, "Unexpected error for input: %s", test.input)
			require.Equal(t, test.want, got, "Mismatch for input: %s", test.input)
			require.ElementsMatch(t, test.wantArgs, gotArgs)
		})
	}
}
