package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tasnimzotder/langz/internal/ast"
)

func TestParsePipeSimple(t *testing.T) {
	prog := parse(`result = data |> upper`)
	require.Len(t, prog.Statements, 1)
	assign := prog.Statements[0].(*ast.Assignment)
	assert.Equal(t, "result", assign.Name)

	pipe, ok := assign.Value.(*ast.BinaryExpr)
	require.True(t, ok, "expected BinaryExpr")
	assert.Equal(t, "|>", pipe.Op)

	left, ok := pipe.Left.(*ast.Identifier)
	require.True(t, ok, "expected Identifier on left")
	assert.Equal(t, "data", left.Name)

	right, ok := pipe.Right.(*ast.Identifier)
	require.True(t, ok, "expected Identifier on right")
	assert.Equal(t, "upper", right.Name)
}

func TestParsePipeChain(t *testing.T) {
	prog := parse(`result = data |> upper |> trim`)
	require.Len(t, prog.Statements, 1)
	assign := prog.Statements[0].(*ast.Assignment)

	// Left-associative: ((data |> upper) |> trim)
	outer, ok := assign.Value.(*ast.BinaryExpr)
	require.True(t, ok, "expected outer BinaryExpr")
	assert.Equal(t, "|>", outer.Op)

	// Right side of outer pipe is "trim"
	right, ok := outer.Right.(*ast.Identifier)
	require.True(t, ok, "expected Identifier on outer right")
	assert.Equal(t, "trim", right.Name)

	// Left side of outer pipe is (data |> upper)
	inner, ok := outer.Left.(*ast.BinaryExpr)
	require.True(t, ok, "expected inner BinaryExpr")
	assert.Equal(t, "|>", inner.Op)

	innerLeft, ok := inner.Left.(*ast.Identifier)
	require.True(t, ok)
	assert.Equal(t, "data", innerLeft.Name)

	innerRight, ok := inner.Right.(*ast.Identifier)
	require.True(t, ok)
	assert.Equal(t, "upper", innerRight.Name)
}

func TestParsePipeAssignment(t *testing.T) {
	prog := parse(`result = name |> upper`)
	require.Len(t, prog.Statements, 1)
	assign, ok := prog.Statements[0].(*ast.Assignment)
	require.True(t, ok)
	assert.Equal(t, "result", assign.Name)

	pipe, ok := assign.Value.(*ast.BinaryExpr)
	require.True(t, ok)
	assert.Equal(t, "|>", pipe.Op)
}

func TestParsePipeWithFuncCall(t *testing.T) {
	prog := parse(`result = data |> json_get(".name")`)
	require.Len(t, prog.Statements, 1)
	assign := prog.Statements[0].(*ast.Assignment)

	pipe, ok := assign.Value.(*ast.BinaryExpr)
	require.True(t, ok)
	assert.Equal(t, "|>", pipe.Op)

	call, ok := pipe.Right.(*ast.FuncCall)
	require.True(t, ok, "expected FuncCall on right")
	assert.Equal(t, "json_get", call.Name)
	require.Len(t, call.Args, 1)
}

func TestParsePipeWithOr(t *testing.T) {
	prog := parse(`result = data |> upper or "default"`)
	require.Len(t, prog.Statements, 1)
	assign := prog.Statements[0].(*ast.Assignment)

	orExpr, ok := assign.Value.(*ast.OrExpr)
	require.True(t, ok, "expected OrExpr")

	pipe, ok := orExpr.Expr.(*ast.BinaryExpr)
	require.True(t, ok, "expected BinaryExpr inside OrExpr")
	assert.Equal(t, "|>", pipe.Op)
}
