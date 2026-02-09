package lsp

import (
	"testing"

	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDocumentSymbolsVariables(t *testing.T) {
	symbols := getDocumentSymbols("x = 1\ny = \"hello\"")
	require.Len(t, symbols, 2)
	assert.Equal(t, "x", symbols[0].Name)
	assert.Equal(t, protocol.SymbolKindVariable, symbols[0].Kind)
	assert.Equal(t, "y", symbols[1].Name)
}

func TestDocumentSymbolsFunction(t *testing.T) {
	symbols := getDocumentSymbols("fn deploy() {\n  print(\"done\")\n}")
	require.Len(t, symbols, 1)
	assert.Equal(t, "deploy", symbols[0].Name)
	assert.Equal(t, protocol.SymbolKindFunction, symbols[0].Kind)
}

func TestDocumentSymbolsMixed(t *testing.T) {
	source := "x = 1\nfn greet() { }\nfor i in items { }"
	symbols := getDocumentSymbols(source)
	require.Len(t, symbols, 3)
	assert.Equal(t, "x", symbols[0].Name)
	assert.Equal(t, protocol.SymbolKindVariable, symbols[0].Kind)
	assert.Equal(t, "greet", symbols[1].Name)
	assert.Equal(t, protocol.SymbolKindFunction, symbols[1].Kind)
	assert.Equal(t, "i", symbols[2].Name)
	assert.Equal(t, protocol.SymbolKindVariable, symbols[2].Kind)
}

func TestDocumentSymbolsEmpty(t *testing.T) {
	symbols := getDocumentSymbols("")
	assert.Empty(t, symbols)
}

func TestDocumentSymbolsPositions(t *testing.T) {
	symbols := getDocumentSymbols("x = 1")
	require.Len(t, symbols, 1)
	assert.Equal(t, uint32(0), symbols[0].Range.Start.Line)
	assert.Equal(t, uint32(0), symbols[0].Range.Start.Character)
	assert.Equal(t, uint32(1), symbols[0].Range.End.Character)
}
