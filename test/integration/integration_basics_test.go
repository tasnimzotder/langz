package integration_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestE2E_HelloWorld(t *testing.T) {
	source := `print("Hello, World!")`
	bash := compileSource(t, source)

	assert.Contains(t, bash, "#!/bin/bash")
	assert.Contains(t, bash, "set -euo pipefail")

	output, code := runBash(t, bash)
	assert.Equal(t, 0, code)
	assert.Equal(t, "Hello, World!", output)
}

func TestE2E_Variables(t *testing.T) {
	source := `
name = "Langz"
version = 42
print("Welcome to {name}")
print(version)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	lines := strings.Split(output, "\n")
	assert.Equal(t, "Welcome to Langz", lines[0])
	assert.Equal(t, "42", lines[1])
}

func TestE2E_Function(t *testing.T) {
	source := `
fn greet(name: str) {
	print("Hello {name}")
}
greet("World")
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "Hello World", output)
}

func TestE2E_FunctionWithMultipleStatements(t *testing.T) {
	source := `
fn deploy(target: str) {
	print("deploying to {target}")
	print("done")
}
deploy("prod")
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	lines := strings.Split(output, "\n")
	assert.Equal(t, "deploying to prod", lines[0])
	assert.Equal(t, "done", lines[1])
}

func TestE2E_MultipleFunctions(t *testing.T) {
	source := `
fn add_prefix(s: str) {
	print("[INFO] {s}")
}

fn log_items(prefix: str) {
	add_prefix("starting {prefix}")
	add_prefix("done {prefix}")
}

log_items("deploy")
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	lines := strings.Split(output, "\n")
	assert.Equal(t, "[INFO] starting deploy", lines[0])
	assert.Equal(t, "[INFO] done deploy", lines[1])
}

func TestE2E_Comments(t *testing.T) {
	source := `
// This is a comment
name = "test"
// Another comment
print(name)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "test", output)
}

func TestE2E_StringInterpolationComplex(t *testing.T) {
	source := `
host = "localhost"
port = 8080
print("Server at {host}:{port}")
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "Server at localhost:8080", output)
}

func TestE2E_BooleanCondition(t *testing.T) {
	source := `
verbose = true
if verbose {
	print("debug on")
}
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "debug on", output)
}

func TestE2E_ExitCode(t *testing.T) {
	source := `exit(42)`
	bash := compileSource(t, source)
	_, code := runBash(t, bash)

	assert.Equal(t, 42, code)
}
