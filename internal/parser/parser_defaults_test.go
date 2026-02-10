package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tasnimzotder/langz/internal/ast"
)

func TestParseDefaultParam(t *testing.T) {
	prog := parse(`fn deploy(target: str = "staging") { print(target) }`)
	require.Len(t, prog.Statements, 1)

	fn := prog.Statements[0].(*ast.FuncDecl)
	require.Len(t, fn.Params, 1)

	param := fn.Params[0]
	assert.Equal(t, "target", param.Name)
	assert.Equal(t, "str", param.Type)
	require.NotNil(t, param.Default)

	str, ok := param.Default.(*ast.StringLiteral)
	require.True(t, ok, "expected StringLiteral default")
	assert.Equal(t, "staging", str.Value)
}

func TestParseMixedParams(t *testing.T) {
	prog := parse(`fn greet(name: str, greeting: str = "Hello") { print(greeting) }`)
	require.Len(t, prog.Statements, 1)

	fn := prog.Statements[0].(*ast.FuncDecl)
	require.Len(t, fn.Params, 2)

	// First param: no default
	assert.Equal(t, "name", fn.Params[0].Name)
	assert.Nil(t, fn.Params[0].Default)

	// Second param: has default
	assert.Equal(t, "greeting", fn.Params[1].Name)
	require.NotNil(t, fn.Params[1].Default)
}

func TestParseDefaultIntParam(t *testing.T) {
	prog := parse(`fn retry(count: int = 3) { print(count) }`)
	require.Len(t, prog.Statements, 1)

	fn := prog.Statements[0].(*ast.FuncDecl)
	param := fn.Params[0]
	require.NotNil(t, param.Default)

	intLit, ok := param.Default.(*ast.IntLiteral)
	require.True(t, ok, "expected IntLiteral default")
	assert.Equal(t, "3", intLit.Value)
}
