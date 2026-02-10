package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tasnimzotder/langz/internal/ast"
)

func TestMapLiteralStringKeys(t *testing.T) {
	prog := parse(`headers = {"Content-Type": "application/json", "Accept": "text/html"}`)

	require.Len(t, prog.Statements, 1)
	assign := prog.Statements[0].(*ast.Assignment)
	m, ok := assign.Value.(*ast.MapLiteral)
	require.True(t, ok, "expected MapLiteral")
	require.Len(t, m.Keys, 2)
	assert.Equal(t, "Content-Type", m.Keys[0])
	assert.Equal(t, "Accept", m.Keys[1])
}

func TestMapLiteralMixedKeys(t *testing.T) {
	prog := parse(`m = {name: "Alice", "X-Custom": "val"}`)

	assign := prog.Statements[0].(*ast.Assignment)
	m, ok := assign.Value.(*ast.MapLiteral)
	require.True(t, ok, "expected MapLiteral")
	require.Len(t, m.Keys, 2)
	assert.Equal(t, "name", m.Keys[0])
	assert.Equal(t, "X-Custom", m.Keys[1])
}

func TestMapLiteralStringKeyInKwarg(t *testing.T) {
	prog := parse(`fetch("url", headers: {"Authorization": token})`)

	call := prog.Statements[0].(*ast.FuncCall)
	require.Len(t, call.KwArgs, 1)
	m, ok := call.KwArgs[0].Value.(*ast.MapLiteral)
	require.True(t, ok, "expected MapLiteral as kwarg value")
	assert.Equal(t, "Authorization", m.Keys[0])
}
