package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tasnimzotder/langz/internal/ast"
)

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

func TestWhileLoop(t *testing.T) {
	prog := parse(`while x > 0 { print(x) }`)
	require.Len(t, prog.Statements, 1)
	ws, ok := prog.Statements[0].(*ast.WhileStmt)
	require.True(t, ok, "expected WhileStmt")

	cond, ok := ws.Condition.(*ast.BinaryExpr)
	require.True(t, ok, "expected BinaryExpr condition")
	assert.Equal(t, ">", cond.Op)
	require.Len(t, ws.Body, 1)
}

func TestMatchStatement(t *testing.T) {
	input := `match platform {
		"darwin" => print("macOS")
		"linux" => print("Linux")
		_ => print("unknown")
	}`
	prog := parse(input)

	require.Len(t, prog.Statements, 1)
	m, ok := prog.Statements[0].(*ast.MatchStmt)
	require.True(t, ok, "expected MatchStmt")

	ident, ok := m.Expr.(*ast.Identifier)
	require.True(t, ok, "expected Identifier")
	assert.Equal(t, "platform", ident.Name)

	require.Len(t, m.Cases, 3)

	// First case: "darwin"
	pattern0, ok := m.Cases[0].Pattern.(*ast.StringLiteral)
	require.True(t, ok, "expected StringLiteral pattern")
	assert.Equal(t, "darwin", pattern0.Value)
	require.Len(t, m.Cases[0].Body, 1)

	// Wildcard case
	assert.Nil(t, m.Cases[2].Pattern, "wildcard _ should be nil")
	require.Len(t, m.Cases[2].Body, 1)
}

func TestMatchWithBlockArm(t *testing.T) {
	input := `match action {
		"test" => print("testing")
		"build" => {
			print("compiling")
			print("linking")
		}
		_ => print("unknown")
	}`
	prog := parse(input)

	require.Len(t, prog.Statements, 1)
	m, ok := prog.Statements[0].(*ast.MatchStmt)
	require.True(t, ok, "expected MatchStmt")
	require.Len(t, m.Cases, 3)

	// Single-statement arm
	p0, ok := m.Cases[0].Pattern.(*ast.StringLiteral)
	require.True(t, ok)
	assert.Equal(t, "test", p0.Value)
	assert.Len(t, m.Cases[0].Body, 1)

	// Block arm with two statements
	p1, ok := m.Cases[1].Pattern.(*ast.StringLiteral)
	require.True(t, ok)
	assert.Equal(t, "build", p1.Value)
	require.Len(t, m.Cases[1].Body, 2)
	assert.IsType(t, &ast.FuncCall{}, m.Cases[1].Body[0])
	assert.IsType(t, &ast.FuncCall{}, m.Cases[1].Body[1])

	// Wildcard arm
	assert.Nil(t, m.Cases[2].Pattern)
	assert.Len(t, m.Cases[2].Body, 1)
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

func TestOrWithBlock(t *testing.T) {
	prog := parse(`x = exec("cmd") or { print("failed") "fallback" }`)

	require.Len(t, prog.Statements, 1)
	assign := prog.Statements[0].(*ast.Assignment)

	orExpr, ok := assign.Value.(*ast.OrExpr)
	require.True(t, ok, "expected OrExpr")

	block, ok := orExpr.Fallback.(*ast.BlockExpr)
	require.True(t, ok, "expected BlockExpr fallback")
	require.Len(t, block.Statements, 2)

	call, ok := block.Statements[0].(*ast.FuncCall)
	require.True(t, ok, "expected FuncCall in block")
	assert.Equal(t, "print", call.Name)

	str, ok := block.Statements[1].(*ast.StringLiteral)
	require.True(t, ok, "expected StringLiteral as last expr")
	assert.Equal(t, "fallback", str.Value)
}

func TestOrWithExitShortcut(t *testing.T) {
	prog := parse(`data = read("file") or exit(1)`)

	require.Len(t, prog.Statements, 1)
	assign := prog.Statements[0].(*ast.Assignment)

	orExpr, ok := assign.Value.(*ast.OrExpr)
	require.True(t, ok, "expected OrExpr")

	call, ok := orExpr.Fallback.(*ast.FuncCall)
	require.True(t, ok, "expected FuncCall fallback (exit)")
	assert.Equal(t, "exit", call.Name)
}

func TestOrWithContinue(t *testing.T) {
	prog := parse(`content = read(f) or continue`)

	require.Len(t, prog.Statements, 1)
	assign := prog.Statements[0].(*ast.Assignment)

	orExpr, ok := assign.Value.(*ast.OrExpr)
	require.True(t, ok, "expected OrExpr")

	_, ok = orExpr.Fallback.(*ast.ContinueStmt)
	require.True(t, ok, "expected ContinueStmt fallback")
}
