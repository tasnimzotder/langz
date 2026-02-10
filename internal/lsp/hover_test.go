package lsp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func hoverValue(t *testing.T, hover *protocol.Hover) string {
	t.Helper()
	mc, ok := hover.Contents.(protocol.MarkupContent)
	require.True(t, ok, "expected MarkupContent")
	return mc.Value
}

func TestHoverOnKwargTimeout(t *testing.T) {
	source := `fetch("url", timeout: 5)`
	hover := getHoverAt(source, 1, 14)
	require.NotNil(t, hover)
	val := hoverValue(t, hover)
	assert.Contains(t, val, "timeout")
	assert.Contains(t, val, "seconds")
}

func TestHoverOnKwargMethod(t *testing.T) {
	source := `fetch("url", method: "POST")`
	hover := getHoverAt(source, 1, 14)
	require.NotNil(t, hover)
	val := hoverValue(t, hover)
	assert.Contains(t, val, "HTTP method")
}

func TestHoverOnKwargNotRegularIdent(t *testing.T) {
	source := `x = 1`
	hover := getHoverAt(source, 1, 1)
	assert.Nil(t, hover)
}

func TestHoverOnBuiltinStillWorks(t *testing.T) {
	source := `fetch("url")`
	hover := getHoverAt(source, 1, 1)
	require.NotNil(t, hover)
	val := hoverValue(t, hover)
	assert.Contains(t, val, "fetch")
}
