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

	assert.Contains(t, output, `rm -f "temp.txt"`)
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

func TestRangeForLoop(t *testing.T) {
	output := body(compile(`for i in range(0, 10) { print(i) }`))

	assert.Contains(t, output, `for i in $(seq 0 10); do`)
	assert.Contains(t, output, `echo "$i"`)
	assert.Contains(t, output, "done")
}

func TestArgsBuiltin(t *testing.T) {
	output := body(compile(`params = args()`))

	assert.Contains(t, output, `params=("$@")`)
}

func TestOsBuiltin(t *testing.T) {
	output := body(compile(`platform = os()`))

	assert.Contains(t, output, `platform=$(uname -s | tr '[:upper:]' '[:lower:]')`)
}

func TestOrWithContinueCodegen(t *testing.T) {
	output := body(compile(`content = read(f) or continue`))

	assert.Contains(t, output, "continue")
}

func TestEqualityComparison(t *testing.T) {
	output := body(compile(`if x == 10 { print("eq") }`))

	assert.Contains(t, output, `if [ "$x" = 10 ]; then`)
}

func TestNotEqualComparison(t *testing.T) {
	output := body(compile(`if x != 10 { print("ne") }`))

	assert.Contains(t, output, `if [ "$x" != 10 ]; then`)
}

func TestLessThanComparison(t *testing.T) {
	output := body(compile(`if x < 10 { print("lt") }`))

	assert.Contains(t, output, `if [ "$x" -lt 10 ]; then`)
}

func TestGreaterEqualComparison(t *testing.T) {
	output := body(compile(`if x >= 10 { print("ge") }`))

	assert.Contains(t, output, `if [ "$x" -ge 10 ]; then`)
}

func TestLessEqualComparison(t *testing.T) {
	output := body(compile(`if x <= 10 { print("le") }`))

	assert.Contains(t, output, `if [ "$x" -le 10 ]; then`)
}

func TestArithmeticExpression(t *testing.T) {
	output := body(compile(`result = a + b`))

	assert.Contains(t, output, `result=$((a + b))`)
}

func TestWhileLoop(t *testing.T) {
	output := body(compile(`while x > 0 { print(x) }`))

	assert.Contains(t, output, `while [ "$x" -gt 0 ]; do`)
	assert.Contains(t, output, `echo "$x"`)
	assert.Contains(t, output, "done")
}

func TestBreakStatement(t *testing.T) {
	output := body(compile(`break`))

	assert.Equal(t, "break", output)
}

func TestLogicalAnd(t *testing.T) {
	output := body(compile(`if a and b { print("both") }`))

	assert.Contains(t, output, `if [ "$a" = true ] && [ "$b" = true ]; then`)
}

func TestSleepBuiltin(t *testing.T) {
	output := body(compile(`sleep(5)`))

	assert.Contains(t, output, "sleep 5")
}

func TestAppendBuiltin(t *testing.T) {
	output := body(compile(`append("log.txt", "entry")`))

	assert.Contains(t, output, `echo "entry" >> "log.txt"`)
}

func TestHostnameBuiltin(t *testing.T) {
	output := body(compile(`host = hostname()`))

	assert.Contains(t, output, `host=$(hostname)`)
}

func TestWhoamiBuiltin(t *testing.T) {
	output := body(compile(`user = whoami()`))

	assert.Contains(t, output, `user=$(whoami)`)
}

func TestArchBuiltin(t *testing.T) {
	output := body(compile(`a = arch()`))

	assert.Contains(t, output, `a=$(uname -m)`)
}

func TestDirnameBuiltin(t *testing.T) {
	output := body(compile(`dir = dirname("/path/to/file.txt")`))

	assert.Contains(t, output, `dir=$(dirname "/path/to/file.txt")`)
}

func TestBasenameBuiltin(t *testing.T) {
	output := body(compile(`name = basename("/path/to/file.txt")`))

	assert.Contains(t, output, `name=$(basename "/path/to/file.txt")`)
}

func TestIsFileBuiltin(t *testing.T) {
	output := body(compile(`if is_file("test.txt") { print("file") }`))

	assert.Contains(t, output, `[ -f "test.txt" ]`)
}

func TestIsDirBuiltin(t *testing.T) {
	output := body(compile(`if is_dir("build") { print("dir") }`))

	assert.Contains(t, output, `[ -d "build" ]`)
}

func TestRmdirBuiltin(t *testing.T) {
	output := body(compile(`rmdir("build")`))

	assert.Contains(t, output, `rm -rf "build"`)
}

func TestUpperBuiltin(t *testing.T) {
	output := body(compile(`x = upper("hello")`))

	assert.Contains(t, output, `x=$(echo "hello" | tr '[:lower:]' '[:upper:]')`)
}

func TestLowerBuiltin(t *testing.T) {
	output := body(compile(`x = lower("HELLO")`))

	assert.Contains(t, output, `x=$(echo "HELLO" | tr '[:upper:]' '[:lower:]')`)
}
