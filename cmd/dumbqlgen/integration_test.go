package main_test

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.tomakado.io/dumbql"
	"go.tomakado.io/dumbql/cmd/dumbqlgen/testdata"
)

func TestMain(m *testing.M) {
	cmd := exec.Command("go", "generate", "./testdata")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	os.Exit(m.Run())
}

func TestGeneratedMatcherWorks(t *testing.T) {
	user := testdata.TestUser{
		ID:        123,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
		Address: testdata.Address{
			Street: "123 Main St",
			City:   "Anytown",
			State:  "CA",
			Zip:    "12345",
		},
	}

	tests := []struct {
		name  string
		query string
		want  bool
	}{
		{
			name:  "Simple field match",
			query: `name = "John Doe"`,
			want:  true,
		},
		{
			name:  "Field not equal",
			query: `name != "Jane Doe"`,
			want:  true,
		},
		{
			name:  "AND condition",
			query: `id = 123 AND email = "john@example.com"`,
			want:  true,
		},
		{
			name:  "OR condition with one true",
			query: `name = "John Doe" OR email = "wrong@example.com"`,
			want:  true,
		},
		{
			name:  "Nested condition",
			query: `(name = "John Doe" AND id = 123) OR email = "wrong@example.com"`,
			want:  true,
		},
		{
			name:  "NOT condition",
			query: `NOT name = "Jane Doe"`,
			want:  true,
		},
		{
			name:  "Negative match",
			query: `name = "Jane Doe" AND email = "john@example.com"`,
			want:  false,
		},
	}

	matcher := testdata.NewTestUserMatcher()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			q, err := dumbql.Parse(tc.query)
			require.NoError(t, err)

			got := q.Match(&user, matcher)
			require.Equal(t, tc.want, got)
		})
	}
}
