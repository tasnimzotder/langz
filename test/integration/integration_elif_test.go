package integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestE2E_Elif(t *testing.T) {
	source := `
x = 2
if x == 1 {
	print("one")
} else if x == 2 {
	print("two")
}
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "two", output)
}

func TestE2E_ElifChain(t *testing.T) {
	source := `
x = 3
if x == 1 {
	print("one")
} else if x == 2 {
	print("two")
} else if x == 3 {
	print("three")
}
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "three", output)
}

func TestE2E_ElifElse(t *testing.T) {
	source := `
x = 99
if x == 1 {
	print("one")
} else if x == 2 {
	print("two")
} else {
	print("other")
}
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "other", output)
}
