package lsp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefinitionFindsVariable(t *testing.T) {
	source := "x = 1\nprint(x)"
	loc := getDefinition(source, "x")
	require.NotNil(t, loc)
	assert.Equal(t, uint32(0), loc.Range.Start.Line)
	assert.Equal(t, uint32(0), loc.Range.Start.Character)
	assert.Equal(t, uint32(1), loc.Range.End.Character)
}

func TestDefinitionFindsFunction(t *testing.T) {
	source := "fn greet() {\n  print(\"hi\")\n}\ngreet()"
	loc := getDefinition(source, "greet")
	require.NotNil(t, loc)
	assert.Equal(t, uint32(0), loc.Range.Start.Line)
	// "greet" starts at col 4 (after "fn ")
	assert.Equal(t, uint32(3), loc.Range.Start.Character)
	assert.Equal(t, uint32(8), loc.Range.End.Character)
}

func TestDefinitionReturnsNilForUnknown(t *testing.T) {
	source := "x = 1"
	loc := getDefinition(source, "y")
	assert.Nil(t, loc)
}

func TestDefinitionFindsFirstOccurrence(t *testing.T) {
	source := "x = 1\nx = 2"
	loc := getDefinition(source, "x")
	require.NotNil(t, loc)
	assert.Equal(t, uint32(0), loc.Range.Start.Line, "should find first definition")
}
