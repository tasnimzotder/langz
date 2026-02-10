package integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestE2E_CompoundAdd(t *testing.T) {
	source := `
x = 10
x += 5
print(x)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "15", output)
}

func TestE2E_CompoundSub(t *testing.T) {
	source := `
x = 10
x -= 3
print(x)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "7", output)
}

func TestE2E_CompoundMul(t *testing.T) {
	source := `
x = 5
x *= 3
print(x)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "15", output)
}
