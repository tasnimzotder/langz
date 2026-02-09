package lsp

import (
	"testing"

	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompletionIncludesBuiltins(t *testing.T) {
	items := getCompletionItems("")
	names := completionNames(items)
	assert.Contains(t, names, "print")
	assert.Contains(t, names, "exec")
	assert.Contains(t, names, "env")
}

func TestCompletionIncludesKeywords(t *testing.T) {
	items := getCompletionItems("")
	names := completionNames(items)
	assert.Contains(t, names, "if")
	assert.Contains(t, names, "fn")
	assert.Contains(t, names, "while")
	assert.Contains(t, names, "for")
}

func TestCompletionIncludesVariables(t *testing.T) {
	items := getCompletionItems("name = \"hello\"\nage = 42")
	names := completionNames(items)
	assert.Contains(t, names, "name")
	assert.Contains(t, names, "age")

	// Check kind is Variable
	for _, item := range items {
		if item.Label == "name" {
			require.NotNil(t, item.Kind)
			assert.Equal(t, protocol.CompletionItemKindVariable, *item.Kind)
		}
	}
}

func TestCompletionIncludesFunctions(t *testing.T) {
	items := getCompletionItems("fn deploy() {\n  print(\"done\")\n}")
	names := completionNames(items)
	assert.Contains(t, names, "deploy")
}

func TestCompletionDeduplicates(t *testing.T) {
	items := getCompletionItems("print(\"hi\")")
	count := 0
	for _, item := range items {
		if item.Label == "print" {
			count++
		}
	}
	assert.Equal(t, 1, count, "print should appear exactly once")
}

func completionNames(items []protocol.CompletionItem) []string {
	names := make([]string, len(items))
	for i, item := range items {
		names[i] = item.Label
	}
	return names
}
