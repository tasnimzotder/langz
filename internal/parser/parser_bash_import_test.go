package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tasnimzotder/langz/internal/ast"
)

func TestParseBashBlock(t *testing.T) {
	prog := parse("bash {\n    set -euo pipefail\n}")

	require.Len(t, prog.Statements, 1)
	bb, ok := prog.Statements[0].(*ast.BashBlock)
	require.True(t, ok, "expected BashBlock")
	assert.Contains(t, bb.Content, "set -euo pipefail")
}

func TestParseBashBlockSimple(t *testing.T) {
	prog := parse(`bash { echo "hello" }`)

	require.Len(t, prog.Statements, 1)
	bb, ok := prog.Statements[0].(*ast.BashBlock)
	require.True(t, ok, "expected BashBlock")
	assert.Equal(t, `echo "hello"`, bb.Content)
}

func TestParseImport(t *testing.T) {
	prog := parse(`import "helpers.lz"`)

	require.Len(t, prog.Statements, 1)
	imp, ok := prog.Statements[0].(*ast.ImportStmt)
	require.True(t, ok, "expected ImportStmt")
	assert.Equal(t, "helpers.lz", imp.Path)
}

func TestParseImportAndCode(t *testing.T) {
	prog := parse("import \"lib.lz\"\nprint(\"hello\")")

	require.Len(t, prog.Statements, 2)
	_, ok := prog.Statements[0].(*ast.ImportStmt)
	require.True(t, ok, "expected ImportStmt")
	_, ok = prog.Statements[1].(*ast.FuncCall)
	require.True(t, ok, "expected FuncCall")
}
