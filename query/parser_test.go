package query_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.tomakado.io/dumbql/query"
)

func TestParser(t *testing.T) { //nolint:funlen
	tests := []struct {
		input string
		want  string
	}{
		// Simple field expression.
		{
			input: "status:200",
			want:  "(= status 200)",
		},
		// Floating-point number.
		{
			input: "eps<0.003",
			want:  "(< eps 0.003000)",
		},
		// Using <= operator.
		{
			input: "eps<=0.003",
			want:  "(<= eps 0.003000)",
		},
		// Using >= operator.
		{
			input: "eps>=0.003",
			want:  "(>= eps 0.003000)",
		},
		// Using > operator.
		{
			input: "eps>0.003",
			want:  "(> eps 0.003000)",
		},
		// Using not-equals with !: operator.
		{
			input: "eps!:0.003",
			want:  "(!= eps 0.003000)",
		},
		// Combined with AND.
		{
			input: "status:200 and eps < 0.003",
			want:  "(and (= status 200) (< eps 0.003000))",
		},
		// Combined with OR.
		{
			input: "status:200 or eps<0.003",
			want:  "(or (= status 200) (< eps 0.003000))",
		},
		// Mixed operators: AND with not-equals.
		{
			input: "status:200 and eps!=0.003",
			want:  "(and (= status 200) (!= eps 0.003000))",
		},
		// Nested parentheses.
		{
			input: "((status:200))",
			want:  "(= status 200)",
		},
		// Extra whitespace.
		{
			input: "   status  :   200    and   eps  <  0.003   ",
			want:  "(and (= status 200) (< eps 0.003000))",
		},
		// Uppercase boolean operator.
		{
			input: "status:200 AND eps<0.003",
			want:  "(and (= status 200) (< eps 0.003000))",
		},
		// Array literal in a field expression.
		{
			input: "req.fields.ext:[\"jpg\", \"png\"]",
			want:  "(= req.fields.ext [\"jpg\" \"png\"])",
		},
		// Array with a single element.
		{
			input: "tags:[\"urgent\"]",
			want:  "(= tags [\"urgent\"])",
		},
		// Empty array literal.
		{
			input: "tags:[]",
			want:  "(= tags [])",
		},
		// A complex expression combining several constructs.
		{
			input: "status : 200 and eps < 0.003 and (req.fields.ext:[\"jpg\", \"png\"])",
			want:  "(and (and (= status 200) (< eps 0.003000)) (= req.fields.ext [\"jpg\" \"png\"]))",
		},
		// NOT with parentheses.
		{
			input: "not (status:200)",
			want:  "(not (= status 200))",
		},
		// Boolean true.
		{
			input: "enabled:true",
			want:  "(= enabled true)",
		},
		// Boolean false.
		{
			input: "enabled:false",
			want:  "(= enabled false)",
		},
		// Boolean with not equals operator.
		{
			input: "enabled!=false",
			want:  "(!= enabled false)",
		},
		// Complex query with boolean.
		{
			input: "status:200 and enabled:true",
			want:  "(and (= status 200) (= enabled true))",
		},
		// Boolean field shorthand syntax.
		{
			input: "enabled",
			want:  "(= enabled true)",
		},
		// Multiple boolean field shorthand syntax.
		{
			input: "enabled and verified",
			want:  "(and (= enabled true) (= verified true))",
		},
		// Boolean field shorthand syntax with other expressions.
		{
			input: "enabled and status:200",
			want:  "(and (= enabled true) (= status 200))",
		},
		// Boolean field shorthand with NOT.
		{
			input: "not enabled",
			want:  "(not (= enabled true))",
		},
		// Complex query with boolean shorthand.
		{
			input: "verified and (status:200 or not enabled)",
			want:  "(and (= verified true) (or (= status 200) (not (= enabled true))))",
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			ast, err := query.Parse("input", []byte(test.input))
			require.NoError(t, err, "parsing error for input: %s", test.input)

			require.Equal(t, test.want, ast.(query.Expr).String())
		})
	}
}
