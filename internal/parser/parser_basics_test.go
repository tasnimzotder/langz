package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tasnimzotder/langz/internal/ast"
	"github.com/tasnimzotder/langz/internal/lexer"
)

func TestAssignString(t *testing.T) {
	prog := parse(`name = "hello"`)

	require.Len(t, prog.Statements, 1)
	assign, ok := prog.Statements[0].(*ast.Assignment)
	require.True(t, ok, "expected Assignment")
	assert.Equal(t, "name", assign.Name)

	str, ok := assign.Value.(*ast.StringLiteral)
	require.True(t, ok, "expected StringLiteral")
	assert.Equal(t, "hello", str.Value)
}

func TestAssignInt(t *testing.T) {
	prog := parse(`count = 42`)

	require.Len(t, prog.Statements, 1)
	assign := prog.Statements[0].(*ast.Assignment)
	assert.Equal(t, "count", assign.Name)

	num, ok := assign.Value.(*ast.IntLiteral)
	require.True(t, ok, "expected IntLiteral")
	assert.Equal(t, "42", num.Value)
}

func TestAssignBool(t *testing.T) {
	prog := parse(`ok = true`)

	require.Len(t, prog.Statements, 1)
	assign := prog.Statements[0].(*ast.Assignment)
	assert.Equal(t, "ok", assign.Name)

	b, ok := assign.Value.(*ast.BoolLiteral)
	require.True(t, ok, "expected BoolLiteral")
	assert.True(t, b.Value)
}

func TestMultipleStatements(t *testing.T) {
	prog := parse("x = 1\ny = 2\nprint(x)")

	require.Len(t, prog.Statements, 3)
}

func TestFuncCall(t *testing.T) {
	prog := parse(`print("hello")`)

	require.Len(t, prog.Statements, 1)
	call, ok := prog.Statements[0].(*ast.FuncCall)
	require.True(t, ok, "expected FuncCall")
	assert.Equal(t, "print", call.Name)
	require.Len(t, call.Args, 1)

	arg, ok := call.Args[0].(*ast.StringLiteral)
	require.True(t, ok, "expected StringLiteral arg")
	assert.Equal(t, "hello", arg.Value)
}

func TestFuncDecl(t *testing.T) {
	prog := parse(`fn greet(name: str) { print(name) }`)

	require.Len(t, prog.Statements, 1)
	fn, ok := prog.Statements[0].(*ast.FuncDecl)
	require.True(t, ok, "expected FuncDecl")
	assert.Equal(t, "greet", fn.Name)

	require.Len(t, fn.Params, 1)
	assert.Equal(t, "name", fn.Params[0].Name)
	assert.Equal(t, "str", fn.Params[0].Type)
	assert.Equal(t, "", fn.ReturnType)

	require.Len(t, fn.Body, 1)
}

func TestFuncDeclWithReturn(t *testing.T) {
	prog := parse(`fn add(a: int, b: int) -> int { return a }`)

	require.Len(t, prog.Statements, 1)
	fn := prog.Statements[0].(*ast.FuncDecl)
	assert.Equal(t, "add", fn.Name)
	assert.Equal(t, "int", fn.ReturnType)

	require.Len(t, fn.Params, 2)
	assert.Equal(t, "a", fn.Params[0].Name)
	assert.Equal(t, "int", fn.Params[0].Type)
	assert.Equal(t, "b", fn.Params[1].Name)
	assert.Equal(t, "int", fn.Params[1].Type)
}

func TestReturnStatement(t *testing.T) {
	prog := parse(`return 42`)

	require.Len(t, prog.Statements, 1)
	ret, ok := prog.Statements[0].(*ast.ReturnStmt)
	require.True(t, ok, "expected ReturnStmt")

	val, ok := ret.Value.(*ast.IntLiteral)
	require.True(t, ok, "expected IntLiteral")
	assert.Equal(t, "42", val.Value)
}

func TestContinueStatement(t *testing.T) {
	prog := parse(`continue`)

	require.Len(t, prog.Statements, 1)
	_, ok := prog.Statements[0].(*ast.ContinueStmt)
	require.True(t, ok, "expected ContinueStmt")
}

func TestBreakStatement(t *testing.T) {
	prog := parse(`break`)
	require.Len(t, prog.Statements, 1)
	_, ok := prog.Statements[0].(*ast.BreakStmt)
	require.True(t, ok, "expected BreakStmt")
}

func TestListLiteral(t *testing.T) {
	prog := parse(`items = ["a", "b", "c"]`)

	require.Len(t, prog.Statements, 1)
	assign := prog.Statements[0].(*ast.Assignment)

	list, ok := assign.Value.(*ast.ListLiteral)
	require.True(t, ok, "expected ListLiteral")
	require.Len(t, list.Elements, 3)

	assert.Equal(t, "a", list.Elements[0].(*ast.StringLiteral).Value)
	assert.Equal(t, "b", list.Elements[1].(*ast.StringLiteral).Value)
	assert.Equal(t, "c", list.Elements[2].(*ast.StringLiteral).Value)
}

func TestEmptyList(t *testing.T) {
	prog := parse(`items = []`)

	require.Len(t, prog.Statements, 1)
	assign := prog.Statements[0].(*ast.Assignment)

	list, ok := assign.Value.(*ast.ListLiteral)
	require.True(t, ok, "expected ListLiteral")
	assert.Len(t, list.Elements, 0)
}

func TestMapLiteral(t *testing.T) {
	prog := parse(`config = {port: 8080, host: "localhost"}`)

	require.Len(t, prog.Statements, 1)
	assign := prog.Statements[0].(*ast.Assignment)

	m, ok := assign.Value.(*ast.MapLiteral)
	require.True(t, ok, "expected MapLiteral")
	require.Len(t, m.Keys, 2)

	assert.Equal(t, "port", m.Keys[0])
	assert.Equal(t, "8080", m.Values[0].(*ast.IntLiteral).Value)
	assert.Equal(t, "host", m.Keys[1])
	assert.Equal(t, "localhost", m.Values[1].(*ast.StringLiteral).Value)
}

func TestDotAccess(t *testing.T) {
	prog := parse(`x = f.name`)

	require.Len(t, prog.Statements, 1)
	assign := prog.Statements[0].(*ast.Assignment)

	dot, ok := assign.Value.(*ast.DotExpr)
	require.True(t, ok, "expected DotExpr")
	assert.Equal(t, "name", dot.Field)

	obj, ok := dot.Object.(*ast.Identifier)
	require.True(t, ok, "expected Identifier")
	assert.Equal(t, "f", obj.Name)
}

func TestParseErrorMessage(t *testing.T) {
	tokens := lexer.New("x = \n+ 1").Tokenize()
	p := New(tokens)
	prog, err := p.ParseWithErrors()

	// Should have errors because + is not a valid expression start
	assert.NotNil(t, prog, "should return partial program even with errors")
	if err != nil {
		assert.Contains(t, err.Error(), "line")
	}
}

func TestParseErrorUnexpectedToken(t *testing.T) {
	tokens := lexer.New("fn (").Tokenize()
	p := New(tokens)
	_, err := p.ParseWithErrors()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "expected IDENT")
}
