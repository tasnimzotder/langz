package lsp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignatureHelpFetch(t *testing.T) {
	// fetch(|) — cursor right after '('
	source := `fetch()`
	result := getSignatureHelp(source, 1, 7)
	require.NotNil(t, result)
	require.Len(t, result.Signatures, 1)
	assert.Contains(t, result.Signatures[0].Label, "fetch")
	assert.Equal(t, uint32(0), *result.ActiveParameter)
}

func TestSignatureHelpFetchSecondArg(t *testing.T) {
	// fetch("url", |) — cursor after the comma
	source := `fetch("url", )`
	result := getSignatureHelp(source, 1, 14)
	require.NotNil(t, result)
	require.NotNil(t, result.ActiveParameter)
	assert.Equal(t, uint32(1), *result.ActiveParameter)
}

func TestSignatureHelpJsonGet(t *testing.T) {
	source := `json_get()`
	result := getSignatureHelp(source, 1, 10)
	require.NotNil(t, result)
	require.Len(t, result.Signatures, 1)
	assert.Contains(t, result.Signatures[0].Label, "json_get")
}

func TestSignatureHelpOutsideCall(t *testing.T) {
	source := `x = 1`
	result := getSignatureHelp(source, 1, 6)
	assert.Nil(t, result)
}

func TestSignatureHelpUnknownFunc(t *testing.T) {
	source := `myFunc()`
	result := getSignatureHelp(source, 1, 8)
	assert.Nil(t, result)
}
