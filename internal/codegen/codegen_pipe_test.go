package codegen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPipeToUpper(t *testing.T) {
	output := body(compile(`
name = "hello"
result = name |> upper
`))
	assert.Contains(t, output, "tr '[:lower:]' '[:upper:]'")
}

func TestPipeToTrim(t *testing.T) {
	output := body(compile(`
name = "hello"
result = name |> trim
`))
	assert.Contains(t, output, "xargs")
}

func TestPipeChain(t *testing.T) {
	output := body(compile(`
name = "hello"
result = name |> upper |> trim
`))
	assert.Contains(t, output, "tr '[:lower:]' '[:upper:]'")
	assert.Contains(t, output, "xargs")
}

func TestPipeWithExtraArg(t *testing.T) {
	output := body(compile(`
data = "test"
result = data |> json_get(".name")
`))
	assert.Contains(t, output, "jq -r")
}

func TestPipeToLower(t *testing.T) {
	output := body(compile(`
name = "HELLO"
result = name |> lower
`))
	assert.Contains(t, output, "tr '[:upper:]' '[:lower:]'")
}
