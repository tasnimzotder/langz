package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tasnimzotder/langz/internal/ast"
)

func TestParseMethodCall(t *testing.T) {
	prog := parse(`val = name.replace("old", "new")`)
	require.Len(t, prog.Statements, 1)
	assign := prog.Statements[0].(*ast.Assignment)

	mc, ok := assign.Value.(*ast.MethodCall)
	require.True(t, ok, "expected MethodCall")

	obj, ok := mc.Object.(*ast.Identifier)
	require.True(t, ok)
	assert.Equal(t, "name", obj.Name)
	assert.Equal(t, "replace", mc.Method)
	assert.Len(t, mc.Args, 2)
}

func TestParseMethodCallNoArgs(t *testing.T) {
	prog := parse(`val = name.upper()`)
	require.Len(t, prog.Statements, 1)
	assign := prog.Statements[0].(*ast.Assignment)

	mc, ok := assign.Value.(*ast.MethodCall)
	require.True(t, ok, "expected MethodCall")
	assert.Equal(t, "upper", mc.Method)
	assert.Len(t, mc.Args, 0)
}

func TestParseMethodChain(t *testing.T) {
	prog := parse(`val = name.replace("a", "b").upper()`)
	require.Len(t, prog.Statements, 1)
	assign := prog.Statements[0].(*ast.Assignment)

	mc, ok := assign.Value.(*ast.MethodCall)
	require.True(t, ok, "expected outer MethodCall")
	assert.Equal(t, "upper", mc.Method)

	inner, ok := mc.Object.(*ast.MethodCall)
	require.True(t, ok, "expected inner MethodCall")
	assert.Equal(t, "replace", inner.Method)
}
