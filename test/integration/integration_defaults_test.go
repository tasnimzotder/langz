package integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestE2E_DefaultParam(t *testing.T) {
	source := `
fn deploy(target: str = "staging") {
	print(target)
}
deploy()
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "staging", output)
}

func TestE2E_DefaultParamOverride(t *testing.T) {
	source := `
fn deploy(target: str = "staging") {
	print(target)
}
deploy("production")
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "production", output)
}

func TestE2E_MixedDefaultParams(t *testing.T) {
	source := `
fn greet(name: str, greeting: str = "Hello") {
	print("{greeting} {name}")
}
greet("World")
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "Hello World", output)
}
