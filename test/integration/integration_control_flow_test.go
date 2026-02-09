package integration_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2E_IfElse(t *testing.T) {
	source := `
x = 5
if x > 3 {
	print("big")
} else {
	print("small")
}
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "big", output)
}

func TestE2E_ForLoop(t *testing.T) {
	source := `
items = ["a", "b", "c"]
for item in items {
	print(item)
}
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	lines := strings.Split(output, "\n")
	assert.Equal(t, []string{"a", "b", "c"}, lines)
}

func TestE2E_RangeLoop(t *testing.T) {
	source := `
for i in range(1, 3) {
	print(i)
}
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	lines := strings.Split(output, "\n")
	assert.Equal(t, []string{"1", "2", "3"}, lines)
}

func TestE2E_ForWithRange(t *testing.T) {
	source := `
total = 0
for i in range(1, 5) {
	print(i)
}
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	lines := strings.Split(output, "\n")
	assert.Equal(t, []string{"1", "2", "3", "4", "5"}, lines)
}

func TestE2E_NestedIfInFor(t *testing.T) {
	source := `
scores = [90, 45, 80]
for s in scores {
	if s > 50 {
		print(s)
	}
}
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	lines := strings.Split(output, "\n")
	assert.Equal(t, []string{"90", "80"}, lines)
}

func TestE2E_WhileWithBreak(t *testing.T) {
	source := `
while true {
	print("once")
	break
}
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "once", output)
}

func TestE2E_MatchStatement(t *testing.T) {
	source := `
platform = "linux"
match platform {
	"darwin" => print("macOS")
	"linux" => print("Linux")
	_ => print("unknown")
}
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "Linux", output)
}

func TestE2E_MatchWithFunctionCalls(t *testing.T) {
	source := `
action = "build"
match action {
	"test" => print("running tests")
	"build" => {
		print("compiling")
		print("linking")
	}
	_ => print("unknown action")
}
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	lines := strings.Split(output, "\n")
	require.Len(t, lines, 2)
	assert.Equal(t, "compiling", lines[0])
	assert.Equal(t, "linking", lines[1])
}

func TestE2E_LogicalOr(t *testing.T) {
	source := `
a = true
b = false
if a or b {
	print("either true")
}

c = false
d = false
if c or d {
	print("should not print")
} else {
	print("both false")
}
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	lines := strings.Split(output, "\n")
	assert.Equal(t, "either true", lines[0])
	assert.Equal(t, "both false", lines[1])
}

func TestE2E_LogicalAndOr(t *testing.T) {
	source := `
a = true
b = false
c = true

// a and b or c => (true and false) or true => false or true => true
if a and b or c {
	print("combined true")
}
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "combined true", output)
}

func TestE2E_OrFallbackStillWorks(t *testing.T) {
	// Verify assignment `or` fallback is not broken
	source := `
name = env("LANGZ_OR_FALLBACK_TEST") or "default_val"
print(name)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "default_val", output)
}

func TestE2E_EnvWithDefault(t *testing.T) {
	source := `name = env("LANGZ_TEST_UNSET_VAR") or "fallback"
print(name)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "fallback", output)
}

func TestE2E_EnvOrExitFallback(t *testing.T) {
	source := `
name = env("LANGZ_E2E_TEST_VAR") or "default_val"
print(name)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "default_val", output)
}

func TestE2E_EqualityComparison(t *testing.T) {
	source := `
mode = "prod"
if mode == "prod" {
	print("production")
}
`
	bash := compileSource(t, source)
	// Verify it generates string comparison
	assert.Contains(t, bash, `[ "$mode" = "prod" ]`)
}

func TestE2E_Arithmetic(t *testing.T) {
	source := `
a = 10
b = 3
result = a + b
print(result)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "13", output)
}

func TestE2E_Modulo(t *testing.T) {
	source := `
a = 10
b = 3
result = a % b
print(result)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "1", strings.TrimSpace(output))
}

func TestE2E_OperatorPrecedence(t *testing.T) {
	source := `
// 2 + 3 * 4 should be 14 (not 20)
result = 2 + 3 * 4
print(result)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "14", strings.TrimSpace(output))
}

func TestE2E_ParenthesizedExpression(t *testing.T) {
	source := `
// (2 + 3) * 4 should be 20
result = (2 + 3) * 4
print(result)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "20", strings.TrimSpace(output))
}

func TestE2E_ComplexMath(t *testing.T) {
	source := `
// Test multiple operations with correct precedence
a = 10
b = 3
c = 2

sum = a + b
diff = a - b
prod = a * b
quot = a / b
rem = a % b
complex = (a + b) * c
nested = a * b + c * b

print(sum)
print(diff)
print(prod)
print(quot)
print(rem)
print(complex)
print(nested)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	lines := strings.Split(strings.TrimSpace(output), "\n")
	require.Len(t, lines, 7)
	assert.Equal(t, "13", lines[0])  // 10 + 3
	assert.Equal(t, "7", lines[1])   // 10 - 3
	assert.Equal(t, "30", lines[2])  // 10 * 3
	assert.Equal(t, "3", lines[3])   // 10 / 3 (integer division)
	assert.Equal(t, "1", lines[4])   // 10 % 3
	assert.Equal(t, "26", lines[5])  // (10 + 3) * 2
	assert.Equal(t, "36", lines[6])  // 10 * 3 + 2 * 3
}

func TestE2E_ArithmeticInCondition(t *testing.T) {
	source := `
a = 5
b = 3
if a + b > 7 {
	print("sum is big")
}
if a * b > 20 {
	print("product is big")
} else {
	print("product is small")
}
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Contains(t, output, "sum is big")
	assert.Contains(t, output, "product is small")
}
