package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tasnimzotder/langz/internal/ast"
)

func TestParseElif(t *testing.T) {
	prog := parse(`
if x == 1 {
	print("one")
} else if x == 2 {
	print("two")
}
`)
	require.Len(t, prog.Statements, 1)
	ifStmt := prog.Statements[0].(*ast.IfStmt)

	// ElseBody should contain a single IfStmt (the elif)
	require.Len(t, ifStmt.ElseBody, 1)
	elif, ok := ifStmt.ElseBody[0].(*ast.IfStmt)
	require.True(t, ok, "expected IfStmt in ElseBody")
	assert.NotNil(t, elif.Condition)
	assert.Len(t, elif.Body, 1)
}

func TestParseElifChain(t *testing.T) {
	prog := parse(`
if x == 1 {
	print("one")
} else if x == 2 {
	print("two")
} else if x == 3 {
	print("three")
}
`)
	require.Len(t, prog.Statements, 1)
	ifStmt := prog.Statements[0].(*ast.IfStmt)

	// First elif
	require.Len(t, ifStmt.ElseBody, 1)
	elif1 := ifStmt.ElseBody[0].(*ast.IfStmt)

	// Second elif
	require.Len(t, elif1.ElseBody, 1)
	_, ok := elif1.ElseBody[0].(*ast.IfStmt)
	require.True(t, ok, "expected second elif")
}

func TestParseElifElse(t *testing.T) {
	prog := parse(`
if x == 1 {
	print("one")
} else if x == 2 {
	print("two")
} else {
	print("other")
}
`)
	require.Len(t, prog.Statements, 1)
	ifStmt := prog.Statements[0].(*ast.IfStmt)

	// First elif
	require.Len(t, ifStmt.ElseBody, 1)
	elif, ok := ifStmt.ElseBody[0].(*ast.IfStmt)
	require.True(t, ok)

	// Final else
	require.True(t, len(elif.ElseBody) > 0, "expected else body")
}
