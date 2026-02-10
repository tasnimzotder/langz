package integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestE2E_PipeToUpper(t *testing.T) {
	source := `
name = "hello"
result = name |> upper
print(result)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "HELLO", output)
}

func TestE2E_PipeToLower(t *testing.T) {
	source := `
name = "HELLO"
result = name |> lower
print(result)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "hello", output)
}

func TestE2E_PipeChain(t *testing.T) {
	source := `
name = "  hello  "
result = name |> trim |> upper
print(result)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "HELLO", output)
}

func TestE2E_PipeLiteralInput(t *testing.T) {
	source := `
result = "world" |> upper
print(result)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "WORLD", output)
}
