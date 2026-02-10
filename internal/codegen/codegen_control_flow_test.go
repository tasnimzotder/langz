package codegen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestRangeForLoop(t *testing.T) {
	output := body(compile(`for i in range(0, 10) { print(i) }`))

	assert.Contains(t, output, `for i in $(seq 0 10); do`)
	assert.Contains(t, output, `echo "$i"`)
	assert.Contains(t, output, "done")
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

func TestContinueStatement(t *testing.T) {
	output := body(compile(`continue`))

	assert.Equal(t, "continue", output)
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

func TestLogicalAnd(t *testing.T) {
	output := body(compile(`if a and b { print("both") }`))

	assert.Contains(t, output, `if [ "$a" = true ] && [ "$b" = true ]; then`)
}

func TestLogicalOr(t *testing.T) {
	output := body(compile(`if a or b { print("either") }`))

	assert.Contains(t, output, `if [ "$a" = true ] || [ "$b" = true ]; then`)
}

func TestLogicalAndOr(t *testing.T) {
	output := body(compile(`if a and b or c { print("yes") }`))

	// and binds tighter, so: (a && b) || c
	assert.Contains(t, output, `if [ "$a" = true ] && [ "$b" = true ] || [ "$c" = true ]; then`)
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

func TestModuloExpression(t *testing.T) {
	output := body(compile(`result = a % b`))

	assert.Contains(t, output, `result=$((a % b))`)
}

func TestOperatorPrecedence(t *testing.T) {
	// a + b * c should generate $((a + b * c)) with correct nesting
	output := body(compile(`result = a + b * c`))

	assert.Contains(t, output, `result=$((a + b * c))`)
}

func TestParenthesizedExprCodegen(t *testing.T) {
	output := body(compile(`result = (a + b) * c`))

	assert.Contains(t, output, `result=$(((a + b) * c))`)
}

func TestComplexArithmetic(t *testing.T) {
	output := body(compile(`result = a * b + c * d`))

	assert.Contains(t, output, `result=$((a * b + c * d))`)
}

func TestUnaryNegation(t *testing.T) {
	output := body(compile(`if !ok { print("failed") }`))

	// Should negate the boolean check, not just test string emptiness
	assert.Contains(t, output, `! [`)
	assert.Contains(t, output, `= true`)
}

func TestUnaryNegationComparison(t *testing.T) {
	output := body(compile(`if !(x > 10) { print("small") }`))

	// Should negate the comparison
	assert.Contains(t, output, `! [`)
}
