package codegen

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

	assert.Contains(t, output, `declare -A config`)
	assert.Contains(t, output, `config[port]=8080`)
	assert.Contains(t, output, `config[host]="localhost"`)
}
