package codegen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tasnimzotder/langz/internal/ast"
)

func TestPreamble(t *testing.T) {
	output := compile(`x = 1`)

	assert.Contains(t, output, "#!/bin/bash")
	assert.Contains(t, output, "set -euo pipefail")
}

func TestAssignString(t *testing.T) {
	output := body(compile(`name = "hello"`))

	assert.Equal(t, `name="hello"`, output)
}

func TestAssignInt(t *testing.T) {
	output := body(compile(`count = 42`))

	assert.Equal(t, `count=42`, output)
}

func TestPrint(t *testing.T) {
	output := body(compile(`print("hello world")`))

	assert.Equal(t, `echo "hello world"`, output)
}

func TestStringInterpolation(t *testing.T) {
	output := body(compile(`print("Hello {name}")`))

	assert.Equal(t, `echo "Hello ${name}"`, output)
}

func TestStringInterpolationMultiple(t *testing.T) {
	output := body(compile(`print("Deploying {app} to {target}")`))

	assert.Equal(t, `echo "Deploying ${app} to ${target}"`, output)
}

func TestStringNoInterpolation(t *testing.T) {
	output := body(compile(`print("no vars here")`))

	assert.Equal(t, `echo "no vars here"`, output)
}

func TestFuncDecl(t *testing.T) {
	output := body(compile(`fn greet(name: str) { print(name) }`))

	assert.Contains(t, output, "greet() {")
	assert.Contains(t, output, `local name="$1"`)
	assert.Contains(t, output, `echo "$name"`)
	assert.Contains(t, output, "}")
}

func TestReturnStatement(t *testing.T) {
	output := body(compile(`return 0`))

	assert.Equal(t, "return 0", output)
}

func TestListLiteralCodegen(t *testing.T) {
	output := body(compile(`items = ["a", "b", "c"]`))

	assert.Equal(t, `items=("a" "b" "c")`, output)
}

func TestEmptyListCodegen(t *testing.T) {
	output := body(compile(`items = []`))

	assert.Equal(t, `items=()`, output)
}

func TestMapLiteralCodegen(t *testing.T) {
	output := body(compile(`config = {port: 8080, host: "localhost"}`))

	assert.Contains(t, output, `config_port=8080`)
	assert.Contains(t, output, `config_host="localhost"`)
}

func TestStringDollarEscaping(t *testing.T) {
	output := body(compile(`print("Cost: $100")`))

	assert.Equal(t, `echo "Cost: \$100"`, output)
}

func TestStringBacktickEscaping(t *testing.T) {
	output := body(compile("print(\"`whoami`\")"))

	assert.Equal(t, "echo \"\\`whoami\\`\"", output)
}

func TestRawValueEscaping(t *testing.T) {
	output := body(compile(`chmod("file.txt", "755")`))

	assert.Contains(t, output, `chmod 755`)
	assert.NotContains(t, output, `;`)
}

func TestMapKeySanitization(t *testing.T) {
	output := body(compile(`x = config["my-key"]`))

	// Hyphens in map keys should be sanitized to underscores
	assert.Contains(t, output, `config_my_key`)
}

func TestCodegenErrorDetection(t *testing.T) {
	_, errs := compileWithErrors(`write("file.txt")`)

	require.Len(t, errs, 1)
	assert.Contains(t, errs[0], "write()")
}

func TestCodegenNoErrorOnValidCode(t *testing.T) {
	_, errs := compileWithErrors(`print("hello")`)

	assert.Empty(t, errs)
}

func TestUnhandledStatementTypeProducesError(t *testing.T) {
	// ExitCall is in the AST but not handled by genStatement
	prog := &ast.Program{
		Statements: []ast.Node{
			&ast.ExitCall{Code: &ast.IntLiteral{Value: "1"}},
		},
	}
	_, errs := Generate(prog)

	require.Len(t, errs, 1)
	assert.Contains(t, errs[0], "unhandled statement type")
}

func TestUnhandledExpressionTypeProducesError(t *testing.T) {
	// Assign an OrExpr without the special case in genAssignment
	// Use a BlockExpr in assignment value (not a special-cased expr type)
	prog := &ast.Program{
		Statements: []ast.Node{
			&ast.Assignment{
				Name:  "x",
				Value: &ast.BlockExpr{Statements: []ast.Node{}},
			},
		},
	}
	_, errs := Generate(prog)

	require.Len(t, errs, 1)
	assert.Contains(t, errs[0], "unhandled expression type")
}
