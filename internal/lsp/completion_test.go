package lsp

import (
	"testing"

	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompletionIncludesBuiltins(t *testing.T) {
	items := getCompletionItems("", 0, 0)
	names := completionNames(items)
	assert.Contains(t, names, "print")
	assert.Contains(t, names, "exec")
	assert.Contains(t, names, "env")
}

func TestCompletionIncludesKeywords(t *testing.T) {
	items := getCompletionItems("", 0, 0)
	names := completionNames(items)
	assert.Contains(t, names, "if")
	assert.Contains(t, names, "fn")
	assert.Contains(t, names, "while")
	assert.Contains(t, names, "for")
}

func TestCompletionIncludesVariables(t *testing.T) {
	items := getCompletionItems("name = \"hello\"\nage = 42", 0, 0)
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
	items := getCompletionItems("fn deploy() {\n  print(\"done\")\n}", 0, 0)
	names := completionNames(items)
	assert.Contains(t, names, "deploy")
}

func TestCompletionDeduplicates(t *testing.T) {
	items := getCompletionItems("print(\"hi\")", 0, 0)
	count := 0
	for _, item := range items {
		if item.Label == "print" {
			count++
		}
	}
	assert.Equal(t, 1, count, "print should appear exactly once")
}

func TestCompletionKwargsInsideFetch(t *testing.T) {
	// Cursor at line 1, col 20: after the comma inside fetch(...)
	// fetch("url", |)
	//                ^ col 14 (1-based)
	source := `fetch("url", )`
	items := getCompletionItems(source, 1, 14)
	names := completionNames(items)
	assert.Contains(t, names, "method:")
	assert.Contains(t, names, "body:")
	assert.Contains(t, names, "headers:")
	assert.Contains(t, names, "timeout:")
	assert.Contains(t, names, "retries:")
}

func TestCompletionKwargsNotOutsideFetch(t *testing.T) {
	source := `x = 1`
	items := getCompletionItems(source, 1, 6)
	names := completionNames(items)
	assert.NotContains(t, names, "method:")
	assert.NotContains(t, names, "timeout:")
}

func TestCompletionKwargsNotInsidePrint(t *testing.T) {
	source := `print()`
	items := getCompletionItems(source, 1, 7)
	names := completionNames(items)
	assert.NotContains(t, names, "method:")
}

func completionNames(items []protocol.CompletionItem) []string {
	names := make([]string, len(items))
	for i, item := range items {
		names[i] = item.Label
	}
	return names
}
