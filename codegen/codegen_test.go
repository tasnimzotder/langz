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

func TestExecBuiltin(t *testing.T) {
	output := body(compile(`result = exec("ls -la")`))

	assert.Contains(t, output, `result=$(ls -la)`)
}

func TestEnvBuiltin(t *testing.T) {
	output := body(compile(`home = env("HOME")`))

	assert.Contains(t, output, `home="${HOME}"`)
}

func TestExistsBuiltin(t *testing.T) {
	output := body(compile(`if exists("file.txt") { print("found") }`))

	assert.Contains(t, output, `[ -e "file.txt" ]`)
}

func TestReadBuiltin(t *testing.T) {
	output := body(compile(`content = read("file.txt")`))

	assert.Contains(t, output, `content=$(cat "file.txt")`)
}

func TestWriteBuiltin(t *testing.T) {
	output := body(compile(`write("out.txt", "hello")`))

	assert.Contains(t, output, `echo "hello" > "out.txt"`)
}

func TestRmBuiltin(t *testing.T) {
	output := body(compile(`rm("temp.txt")`))

	assert.Contains(t, output, `rm "temp.txt"`)
}

func TestMkdirBuiltin(t *testing.T) {
	output := body(compile(`mkdir("build/output")`))

	assert.Contains(t, output, `mkdir -p "build/output"`)
}

func TestCopyBuiltin(t *testing.T) {
	output := body(compile(`copy("a.txt", "b.txt")`))

	assert.Contains(t, output, `cp "a.txt" "b.txt"`)
}

func TestMoveBuiltin(t *testing.T) {
	output := body(compile(`move("old.txt", "new.txt")`))

	assert.Contains(t, output, `mv "old.txt" "new.txt"`)
}

func TestChmodBuiltin(t *testing.T) {
	output := body(compile(`chmod("script.sh", "755")`))

	assert.Contains(t, output, `chmod 755 "script.sh"`)
}

func TestGlobBuiltin(t *testing.T) {
	output := body(compile(`files = glob("*.log")`))

	assert.Contains(t, output, `files=(*.log)`)
}

func TestExitBuiltin(t *testing.T) {
	output := body(compile(`exit(1)`))

	assert.Contains(t, output, "exit 1")
}

func TestOrWithValue(t *testing.T) {
	output := body(compile(`name = env("APP") or "default"`))

	assert.Contains(t, output, `name="${APP:-default}"`)
}

func TestOrWithExitCodegen(t *testing.T) {
	output := body(compile(`data = read("f.txt") or exit(1)`))

	assert.Contains(t, output, `cat "f.txt"`)
	assert.Contains(t, output, "exit 1")
}

func TestOrWithBlockCodegen(t *testing.T) {
	output := body(compile(`x = exec("cmd") or { print("failed") "fallback" }`))

	assert.Contains(t, output, "cmd")
	assert.Contains(t, output, `echo "failed"`)
}

func TestMatchCodegen(t *testing.T) {
	input := `match platform {
		"darwin" => print("macOS")
		"linux" => print("Linux")
		_ => print("unknown")
	}`
	output := body(compile(input))

	assert.Contains(t, output, `case "$platform" in`)
	assert.Contains(t, output, `darwin)`)
	assert.Contains(t, output, `echo "macOS"`)
	assert.Contains(t, output, `;;`)
	assert.Contains(t, output, `linux)`)
	assert.Contains(t, output, `*)`)
	assert.Contains(t, output, `esac`)
}

func TestFetchBuiltin(t *testing.T) {
	output := body(compile(`res = fetch("https://api.example.com/health")`))

	assert.Contains(t, output, `curl -sf "https://api.example.com/health"`)
}

func TestOrWithContinueCodegen(t *testing.T) {
	output := body(compile(`content = read(f) or continue`))

	assert.Contains(t, output, "continue")
}
