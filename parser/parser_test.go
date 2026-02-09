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

func TestNegation(t *testing.T) {
	prog := parse(`if !ok { print("fail") }`)

	require.Len(t, prog.Statements, 1)
	ifStmt := prog.Statements[0].(*ast.IfStmt)

	unary, ok := ifStmt.Condition.(*ast.UnaryExpr)
	require.True(t, ok, "expected UnaryExpr")
	assert.Equal(t, "!", unary.Op)
}

func TestComparisonOperators(t *testing.T) {
	tests := []struct {
		input string
		op    string
	}{
		{`if x == 10 { print("eq") }`, "=="},
		{`if x != 10 { print("ne") }`, "!="},
		{`if x < 10 { print("lt") }`, "<"},
		{`if x >= 10 { print("ge") }`, ">="},
		{`if x <= 10 { print("le") }`, "<="},
		{`if x > 10 { print("gt") }`, ">"},
	}

	for _, tt := range tests {
		t.Run(tt.op, func(t *testing.T) {
			prog := parse(tt.input)
			require.Len(t, prog.Statements, 1)
			ifStmt := prog.Statements[0].(*ast.IfStmt)
			cond, ok := ifStmt.Condition.(*ast.BinaryExpr)
			require.True(t, ok, "expected BinaryExpr")
			assert.Equal(t, tt.op, cond.Op)
		})
	}
}

func TestArithmeticExpression(t *testing.T) {
	prog := parse(`x = a + b`)
	require.Len(t, prog.Statements, 1)
	assign := prog.Statements[0].(*ast.Assignment)
	bin, ok := assign.Value.(*ast.BinaryExpr)
	require.True(t, ok, "expected BinaryExpr")
	assert.Equal(t, "+", bin.Op)
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

func TestBreakStatement(t *testing.T) {
	prog := parse(`break`)
	require.Len(t, prog.Statements, 1)
	_, ok := prog.Statements[0].(*ast.BreakStmt)
	require.True(t, ok, "expected BreakStmt")
}

func TestLogicalAnd(t *testing.T) {
	prog := parse(`if a and b { print("both") }`)
	require.Len(t, prog.Statements, 1)
	ifStmt := prog.Statements[0].(*ast.IfStmt)
	bin, ok := ifStmt.Condition.(*ast.BinaryExpr)
	require.True(t, ok, "expected BinaryExpr")
	assert.Equal(t, "and", bin.Op)
}
