package dumbql_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.tomakado.io/dumbql"
)

func TestQuery_ToSql(t *testing.T) {
	q, err := dumbql.Parse("status:200")
	require.NoError(t, err)

	sql, args, err := q.ToSql()
	require.NoError(t, err)
	assert.Contains(t, sql, "status = ?")
	assert.Equal(t, []any{float64(200)}, args)
}