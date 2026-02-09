package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tasnimzotder/langz/internal/ast"
)

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

func TestModuloExpression(t *testing.T) {
	prog := parse(`x = a % b`)
	require.Len(t, prog.Statements, 1)
	assign := prog.Statements[0].(*ast.Assignment)
	bin, ok := assign.Value.(*ast.BinaryExpr)
	require.True(t, ok, "expected BinaryExpr")
	assert.Equal(t, "%", bin.Op)
}

func TestLogicalAnd(t *testing.T) {
	prog := parse(`if a and b { print("both") }`)
	require.Len(t, prog.Statements, 1)
	ifStmt := prog.Statements[0].(*ast.IfStmt)
	bin, ok := ifStmt.Condition.(*ast.BinaryExpr)
	require.True(t, ok, "expected BinaryExpr")
	assert.Equal(t, "and", bin.Op)
}

func TestLogicalOr(t *testing.T) {
	prog := parse(`if a or b { print("either") }`)
	require.Len(t, prog.Statements, 1)
	ifStmt := prog.Statements[0].(*ast.IfStmt)
	bin, ok := ifStmt.Condition.(*ast.BinaryExpr)
	require.True(t, ok, "expected BinaryExpr")
	assert.Equal(t, "or", bin.Op)

	left, ok := bin.Left.(*ast.Identifier)
	require.True(t, ok, "expected Identifier on left")
	assert.Equal(t, "a", left.Name)

	right, ok := bin.Right.(*ast.Identifier)
	require.True(t, ok, "expected Identifier on right")
	assert.Equal(t, "b", right.Name)
}

func TestLogicalAndOr(t *testing.T) {
	// `a and b or c` should parse as `(a and b) or c`
	// because `and` has higher precedence than `or`
	prog := parse(`if a and b or c { print("yes") }`)
	require.Len(t, prog.Statements, 1)
	ifStmt := prog.Statements[0].(*ast.IfStmt)

	// Top-level should be `or`
	orExpr, ok := ifStmt.Condition.(*ast.BinaryExpr)
	require.True(t, ok, "expected BinaryExpr at top level")
	assert.Equal(t, "or", orExpr.Op)

	// Left of `or` should be `a and b`
	andExpr, ok := orExpr.Left.(*ast.BinaryExpr)
	require.True(t, ok, "expected BinaryExpr on left of or")
	assert.Equal(t, "and", andExpr.Op)

	// Right of `or` should be identifier `c`
	right, ok := orExpr.Right.(*ast.Identifier)
	require.True(t, ok, "expected Identifier on right of or")
	assert.Equal(t, "c", right.Name)
}

func TestLogicalOrInWhile(t *testing.T) {
	prog := parse(`while a or b { print("loop") }`)
	require.Len(t, prog.Statements, 1)
	ws, ok := prog.Statements[0].(*ast.WhileStmt)
	require.True(t, ok, "expected WhileStmt")

	bin, ok := ws.Condition.(*ast.BinaryExpr)
	require.True(t, ok, "expected BinaryExpr")
	assert.Equal(t, "or", bin.Op)
}

func TestOperatorPrecedence(t *testing.T) {
	// a + b * c should parse as a + (b * c)
	prog := parse(`x = a + b * c`)
	require.Len(t, prog.Statements, 1)
	assign := prog.Statements[0].(*ast.Assignment)

	// Top-level should be +
	add, ok := assign.Value.(*ast.BinaryExpr)
	require.True(t, ok, "expected BinaryExpr at top level")
	assert.Equal(t, "+", add.Op)

	// Left of + should be identifier 'a'
	left, ok := add.Left.(*ast.Identifier)
	require.True(t, ok, "expected Identifier on left of +")
	assert.Equal(t, "a", left.Name)

	// Right of + should be b * c
	mul, ok := add.Right.(*ast.BinaryExpr)
	require.True(t, ok, "expected BinaryExpr on right of +")
	assert.Equal(t, "*", mul.Op)
}

func TestOperatorPrecedenceSubtractDivide(t *testing.T) {
	// a - b / c should parse as a - (b / c)
	prog := parse(`x = a - b / c`)
	require.Len(t, prog.Statements, 1)
	assign := prog.Statements[0].(*ast.Assignment)

	sub, ok := assign.Value.(*ast.BinaryExpr)
	require.True(t, ok, "expected BinaryExpr")
	assert.Equal(t, "-", sub.Op)

	div, ok := sub.Right.(*ast.BinaryExpr)
	require.True(t, ok, "expected BinaryExpr on right of -")
	assert.Equal(t, "/", div.Op)
}

func TestNestedArithmetic(t *testing.T) {
	// a * b + c * d should parse as (a * b) + (c * d)
	prog := parse(`x = a * b + c * d`)
	require.Len(t, prog.Statements, 1)
	assign := prog.Statements[0].(*ast.Assignment)

	add, ok := assign.Value.(*ast.BinaryExpr)
	require.True(t, ok, "expected BinaryExpr")
	assert.Equal(t, "+", add.Op)

	leftMul, ok := add.Left.(*ast.BinaryExpr)
	require.True(t, ok, "expected BinaryExpr on left")
	assert.Equal(t, "*", leftMul.Op)

	rightMul, ok := add.Right.(*ast.BinaryExpr)
	require.True(t, ok, "expected BinaryExpr on right")
	assert.Equal(t, "*", rightMul.Op)
}

func TestParenthesizedExpression(t *testing.T) {
	// (a + b) * c should parse as (a + b) * c
	prog := parse(`x = (a + b) * c`)
	require.Len(t, prog.Statements, 1)
	assign := prog.Statements[0].(*ast.Assignment)

	// Top-level should be *
	mul, ok := assign.Value.(*ast.BinaryExpr)
	require.True(t, ok, "expected BinaryExpr at top level")
	assert.Equal(t, "*", mul.Op)

	// Left of * should be a + b (from parentheses)
	add, ok := mul.Left.(*ast.BinaryExpr)
	require.True(t, ok, "expected BinaryExpr on left of *")
	assert.Equal(t, "+", add.Op)

	// Right of * should be identifier 'c'
	right, ok := mul.Right.(*ast.Identifier)
	require.True(t, ok, "expected Identifier on right of *")
	assert.Equal(t, "c", right.Name)
}

func TestArithmeticWithComparison(t *testing.T) {
	// if a + b > c * d — comparison should be at top, arithmetic deeper
	prog := parse(`if a + b > c * d { print("yes") }`)
	require.Len(t, prog.Statements, 1)
	ifStmt := prog.Statements[0].(*ast.IfStmt)

	cond, ok := ifStmt.Condition.(*ast.BinaryExpr)
	require.True(t, ok, "expected BinaryExpr")
	assert.Equal(t, ">", cond.Op)

	left, ok := cond.Left.(*ast.BinaryExpr)
	require.True(t, ok, "expected BinaryExpr on left of >")
	assert.Equal(t, "+", left.Op)

	right, ok := cond.Right.(*ast.BinaryExpr)
	require.True(t, ok, "expected BinaryExpr on right of >")
	assert.Equal(t, "*", right.Op)
}
