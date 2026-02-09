package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tasnimzotder/langz/ast"
	"github.com/tasnimzotder/langz/lexer"
)

func parse(input string) *ast.Program {
	tokens := lexer.New(input).Tokenize()
	return New(tokens).Parse()
}

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

func TestAssignWithOr(t *testing.T) {
	prog := parse(`name = env("APP") or "default"`)

	require.Len(t, prog.Statements, 1)
	assign := prog.Statements[0].(*ast.Assignment)
	assert.Equal(t, "name", assign.Name)

	orExpr, ok := assign.Value.(*ast.OrExpr)
	require.True(t, ok, "expected OrExpr")

	call, ok := orExpr.Expr.(*ast.FuncCall)
	require.True(t, ok, "expected FuncCall in or expr")
	assert.Equal(t, "env", call.Name)

	fallback, ok := orExpr.Fallback.(*ast.StringLiteral)
	require.True(t, ok, "expected StringLiteral fallback")
	assert.Equal(t, "default", fallback.Value)
}

func TestMultipleStatements(t *testing.T) {
	prog := parse("x = 1\ny = 2\nprint(x)")

	require.Len(t, prog.Statements, 3)
}

func TestIfStatement(t *testing.T) {
	prog := parse(`if x > 10 { print("big") }`)

	require.Len(t, prog.Statements, 1)
	ifStmt, ok := prog.Statements[0].(*ast.IfStmt)
	require.True(t, ok, "expected IfStmt")

	cond, ok := ifStmt.Condition.(*ast.BinaryExpr)
	require.True(t, ok, "expected BinaryExpr condition")
	assert.Equal(t, ">", cond.Op)

	require.Len(t, ifStmt.Body, 1)
	assert.IsType(t, &ast.FuncCall{}, ifStmt.Body[0])
	assert.Nil(t, ifStmt.ElseBody)
}

func TestIfElse(t *testing.T) {
	prog := parse(`if ok { print("yes") } else { print("no") }`)

	require.Len(t, prog.Statements, 1)
	ifStmt := prog.Statements[0].(*ast.IfStmt)

	require.Len(t, ifStmt.Body, 1)
	require.Len(t, ifStmt.ElseBody, 1)
}

func TestForLoop(t *testing.T) {
	prog := parse(`for f in files { print(f) }`)

	require.Len(t, prog.Statements, 1)
	forStmt, ok := prog.Statements[0].(*ast.ForStmt)
	require.True(t, ok, "expected ForStmt")
	assert.Equal(t, "f", forStmt.Var)

	collection, ok := forStmt.Collection.(*ast.Identifier)
	require.True(t, ok, "expected Identifier collection")
	assert.Equal(t, "files", collection.Name)

	require.Len(t, forStmt.Body, 1)
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

func TestNegation(t *testing.T) {
	prog := parse(`if !ok { print("fail") }`)

	require.Len(t, prog.Statements, 1)
	ifStmt := prog.Statements[0].(*ast.IfStmt)

	unary, ok := ifStmt.Condition.(*ast.UnaryExpr)
	require.True(t, ok, "expected UnaryExpr")
	assert.Equal(t, "!", unary.Op)
}
