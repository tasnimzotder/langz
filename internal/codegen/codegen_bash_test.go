package codegen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBashBlockCodegen(t *testing.T) {
	output := compile(`bash { echo "hello world" }`)
	b := body(output)
	assert.Equal(t, `echo "hello world"`, b)
}

func TestBashBlockMultiline(t *testing.T) {
	output := compile("bash {\n    MY_VAR=1\n    trap 'cleanup' EXIT\n}")
	b := body(output)
	assert.Contains(t, b, "MY_VAR=1")
	assert.Contains(t, b, "trap 'cleanup' EXIT")
}

func TestBashBlockWithSurroundingCode(t *testing.T) {
	source := "x = \"before\"\nbash { echo inline }\nprint(x)"
	output := compile(source)
	b := body(output)
	assert.Contains(t, b, `x="before"`)
	assert.Contains(t, b, "echo inline")
	assert.Contains(t, b, `echo "$x"`)
}
