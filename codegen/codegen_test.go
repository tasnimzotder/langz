package codegen

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tasnimzotder/langz/lexer"
	"github.com/tasnimzotder/langz/parser"
)

func compile(input string) string {
	tokens := lexer.New(input).Tokenize()
	prog := parser.New(tokens).Parse()
	return Generate(prog)
}

func body(output string) string {
	// Strip the preamble (#!/bin/bash and set -euo pipefail)
	lines := strings.Split(output, "\n")
	var result []string
	for _, line := range lines {
		if line == "#!/bin/bash" || line == "set -euo pipefail" || line == "" {
			continue
		}
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}

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

func TestFuncDecl(t *testing.T) {
	output := body(compile(`fn greet(name: str) { print(name) }`))

	assert.Contains(t, output, "greet() {")
	assert.Contains(t, output, `local name="$1"`)
	assert.Contains(t, output, `echo "$name"`)
	assert.Contains(t, output, "}")
}

func TestIfStatement(t *testing.T) {
	output := body(compile(`if x > 10 { print("big") }`))

	assert.Contains(t, output, `if [ "$x" -gt 10 ]; then`)
	assert.Contains(t, output, `echo "big"`)
	assert.Contains(t, output, "fi")
}

func TestIfElse(t *testing.T) {
	output := body(compile(`if ok { print("yes") } else { print("no") }`))

	assert.Contains(t, output, "else")
}

func TestForLoop(t *testing.T) {
	output := body(compile(`for f in files { print(f) }`))

	assert.Contains(t, output, `for f in "${files[@]}"; do`)
	assert.Contains(t, output, `echo "$f"`)
	assert.Contains(t, output, "done")
}

func TestReturnStatement(t *testing.T) {
	output := body(compile(`return 0`))

	assert.Equal(t, "return 0", output)
}

func TestContinueStatement(t *testing.T) {
	output := body(compile(`continue`))

	assert.Equal(t, "continue", output)
}
