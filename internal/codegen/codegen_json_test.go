package codegen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonGetCodegen(t *testing.T) {
	output := body(compile(`name = json_get(data, ".name")`))

	assert.Contains(t, output, `$(echo "$data" | jq -r ".name")`)
}

func TestJsonGetNestedPath(t *testing.T) {
	output := body(compile(`city = json_get(resp, ".address.city")`))

	assert.Contains(t, output, `$(echo "$resp" | jq -r ".address.city")`)
}

func TestJsonGetErrorNoArgs(t *testing.T) {
	output := body(compile(`x = json_get()`))

	assert.Contains(t, output, `# error: json_get() requires 2 arguments`)
}
