package lsp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tasnimzotder/langz/internal/lexer"
)

func TestFindSymbolsVariable(t *testing.T) {
	symbols := findSymbols(`name = "hello"`)
	require.Len(t, symbols, 1)
	assert.Equal(t, "name", symbols[0].Name)
	assert.Equal(t, "variable", symbols[0].Kind)
	assert.Equal(t, 1, symbols[0].Line)
	assert.Equal(t, 1, symbols[0].Col)
}

func TestFindSymbolsFunction(t *testing.T) {
	symbols := findSymbols("fn deploy() {\n  print(\"done\")\n}")
	require.Len(t, symbols, 1)
	assert.Equal(t, "deploy", symbols[0].Name)
	assert.Equal(t, "function", symbols[0].Kind)
}

func TestFindSymbolsForVar(t *testing.T) {
	symbols := findSymbols("for item in items {\n  print(item)\n}")
	require.Len(t, symbols, 1)
	assert.Equal(t, "item", symbols[0].Name)
	assert.Equal(t, "for_var", symbols[0].Kind)
}

func TestFindSymbolsMixed(t *testing.T) {
	source := `x = 1
fn greet(name: string) {
  print(name)
}
for i in items {
  print(i)
}
y = 2`
	symbols := findSymbols(source)
	require.Len(t, symbols, 4)
	assert.Equal(t, "x", symbols[0].Name)
	assert.Equal(t, "variable", symbols[0].Kind)
	assert.Equal(t, "greet", symbols[1].Name)
	assert.Equal(t, "function", symbols[1].Kind)
	assert.Equal(t, "i", symbols[2].Name)
	assert.Equal(t, "for_var", symbols[2].Kind)
	assert.Equal(t, "y", symbols[3].Name)
	assert.Equal(t, "variable", symbols[3].Kind)
}

func TestFindSymbolsEmpty(t *testing.T) {
	symbols := findSymbols("")
	assert.Empty(t, symbols)
}

func TestKeywordNames(t *testing.T) {
	names := lexer.KeywordNames()
	assert.Len(t, names, 16)
	assert.Contains(t, names, "if")
	assert.Contains(t, names, "fn")
	assert.Contains(t, names, "while")
	assert.Contains(t, names, "or")
	assert.Contains(t, names, "bash")
	assert.Contains(t, names, "import")
}
