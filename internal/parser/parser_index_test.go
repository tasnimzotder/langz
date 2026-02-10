package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tasnimzotder/langz/internal/ast"
)

func TestParseArrayIndex(t *testing.T) {
	prog := parse(`val = arr[0]`)
	require.Len(t, prog.Statements, 1)
	assign := prog.Statements[0].(*ast.Assignment)

	idx, ok := assign.Value.(*ast.IndexExpr)
	require.True(t, ok, "expected IndexExpr")

	obj, ok := idx.Object.(*ast.Identifier)
	require.True(t, ok)
	assert.Equal(t, "arr", obj.Name)

	index, ok := idx.Index.(*ast.IntLiteral)
	require.True(t, ok)
	assert.Equal(t, "0", index.Value)
}

func TestParseMapIndex(t *testing.T) {
	prog := parse(`val = config["host"]`)
	require.Len(t, prog.Statements, 1)
	assign := prog.Statements[0].(*ast.Assignment)

	idx, ok := assign.Value.(*ast.IndexExpr)
	require.True(t, ok, "expected IndexExpr")

	index, ok := idx.Index.(*ast.StringLiteral)
	require.True(t, ok)
	assert.Equal(t, "host", index.Value)
}

func TestParseIndexAssign(t *testing.T) {
	prog := parse(`arr[0] = "new"`)
	require.Len(t, prog.Statements, 1)

	ia, ok := prog.Statements[0].(*ast.IndexAssignment)
	require.True(t, ok, "expected IndexAssignment")
	assert.Equal(t, "arr", ia.Object)
}
