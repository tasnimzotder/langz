package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tasnimzotder/langz/internal/ast"
	"github.com/tasnimzotder/langz/internal/lexer"
)

func TestFuncCallWithOneKwarg(t *testing.T) {
	prog := parse(`fetch("https://api.com", method: "POST")`)

	require.Len(t, prog.Statements, 1)
	call, ok := prog.Statements[0].(*ast.FuncCall)
	require.True(t, ok, "expected FuncCall")
	assert.Equal(t, "fetch", call.Name)
	require.Len(t, call.Args, 1)
	require.Len(t, call.KwArgs, 1)
	assert.Equal(t, "method", call.KwArgs[0].Key)

	val, ok := call.KwArgs[0].Value.(*ast.StringLiteral)
	require.True(t, ok, "expected StringLiteral kwarg value")
	assert.Equal(t, "POST", val.Value)
}

func TestFuncCallWithMultipleKwargs(t *testing.T) {
	prog := parse(`fetch("url", method: "POST", timeout: 30)`)

	call := prog.Statements[0].(*ast.FuncCall)
	require.Len(t, call.Args, 1)
	require.Len(t, call.KwArgs, 2)
	assert.Equal(t, "method", call.KwArgs[0].Key)
	assert.Equal(t, "timeout", call.KwArgs[1].Key)

	timeout, ok := call.KwArgs[1].Value.(*ast.IntLiteral)
	require.True(t, ok, "expected IntLiteral")
	assert.Equal(t, "30", timeout.Value)
}

func TestFuncCallKwargsOnly(t *testing.T) {
	prog := parse(`fetch(url: "https://api.com", method: "GET")`)

	call := prog.Statements[0].(*ast.FuncCall)
	assert.Len(t, call.Args, 0)
	require.Len(t, call.KwArgs, 2)
	assert.Equal(t, "url", call.KwArgs[0].Key)
	assert.Equal(t, "method", call.KwArgs[1].Key)
}

func TestFuncCallNoKwargs(t *testing.T) {
	prog := parse(`print("hello", "world")`)

	call := prog.Statements[0].(*ast.FuncCall)
	assert.Len(t, call.Args, 2)
	assert.Len(t, call.KwArgs, 0)
}

func TestFuncCallKwargWithMapValue(t *testing.T) {
	prog := parse(`fetch("url", headers: {content_type: "json"})`)

	call := prog.Statements[0].(*ast.FuncCall)
	require.Len(t, call.Args, 1)
	require.Len(t, call.KwArgs, 1)
	assert.Equal(t, "headers", call.KwArgs[0].Key)

	_, ok := call.KwArgs[0].Value.(*ast.MapLiteral)
	assert.True(t, ok, "expected MapLiteral as kwarg value")
}

func TestFuncCallKwargWithIdentValue(t *testing.T) {
	prog := parse(`fetch("url", body: payload)`)

	call := prog.Statements[0].(*ast.FuncCall)
	require.Len(t, call.KwArgs, 1)
	assert.Equal(t, "body", call.KwArgs[0].Key)

	ident, ok := call.KwArgs[0].Value.(*ast.Identifier)
	require.True(t, ok, "expected Identifier as kwarg value")
	assert.Equal(t, "payload", ident.Name)
}

func TestPositionalAfterKwargError(t *testing.T) {
	tokens := lexer.New(`f(a: 1, "positional")`).Tokenize()
	p := New(tokens)
	_, err := p.ParseWithErrors()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "positional argument after keyword argument")
}
