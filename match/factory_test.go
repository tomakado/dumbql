package match_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.tomakado.io/dumbql/match"
	"go.tomakado.io/dumbql/query"
)

type testStruct struct {
	Name string `dumbql:"name"`
	Age  int    `dumbql:"age"`
}

// MockMatcher is a mock implementation of query.Matcher for testing
type MockMatcher struct{}

func (m *MockMatcher) MatchAnd(_ any, _, _ query.Expr) bool {
	return true
}

func (m *MockMatcher) MatchOr(_ any, _, _ query.Expr) bool {
	return true
}

func (m *MockMatcher) MatchNot(_ any, _ query.Expr) bool {
	return true
}

func (m *MockMatcher) MatchField(_ any, _ string, _ query.Valuer, _ query.FieldOperator) bool {
	return true
}

func (m *MockMatcher) MatchValue(_ any, _ query.Valuer, _ query.FieldOperator) bool {
	return true
}

func TestCreateMatcher(t *testing.T) {
	t.Run("reflection matcher", func(t *testing.T) {
		target := testStruct{}
		matcher, err := match.CreateMatcher(target, match.MatcherTypeReflection)
		require.NoError(t, err)
		assert.IsType(t, &match.StructMatcher{}, matcher)
	})

	t.Run("generated matcher", func(t *testing.T) {
		target := testStruct{}
		mockMatcher := &MockMatcher{}

		// Register mock matcher
		match.RegisterGeneratedMatcher("testStruct", mockMatcher)

		matcher, err := match.CreateMatcher(target, match.MatcherTypeGenerated)
		require.NoError(t, err)
		assert.Equal(t, mockMatcher, matcher)
	})

	t.Run("generated matcher fallback", func(t *testing.T) {
		type unregisteredStruct struct{}
		target := unregisteredStruct{}

		matcher, err := match.CreateMatcher(target, match.MatcherTypeGenerated)
		require.NoError(t, err)
		assert.IsType(t, &match.StructMatcher{}, matcher)
	})

	t.Run("unknown matcher type", func(t *testing.T) {
		target := testStruct{}
		_, err := match.CreateMatcher(target, 999)
		assert.Error(t, err)
	})
}
