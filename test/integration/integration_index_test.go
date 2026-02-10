package integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestE2E_ArrayIndex(t *testing.T) {
	source := `
items = ["alpha", "beta", "gamma"]
val = items[1]
print(val)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "beta", output)
}

func TestE2E_MapIndex(t *testing.T) {
	source := `
config = {host: "localhost", port: "8080"}
val = config["host"]
print(val)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "localhost", output)
}

func TestE2E_IndexAssign(t *testing.T) {
	source := `
items = ["a", "b", "c"]
items[1] = "B"
print(items[1])
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "B", output)
}
