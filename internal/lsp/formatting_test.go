package lsp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatIndentation(t *testing.T) {
	input := "if x > 1 {\nprint(\"yes\")\n}"
	expected := "if x > 1 {\n    print(\"yes\")\n}"
	assert.Equal(t, expected, formatSource(input, 4, true))
}

func TestFormatNestedBlocks(t *testing.T) {
	input := "if x > 1 {\nif y > 2 {\nprint(\"deep\")\n}\n}"
	expected := "if x > 1 {\n    if y > 2 {\n        print(\"deep\")\n    }\n}"
	assert.Equal(t, expected, formatSource(input, 4, true))
}

func TestFormatAlreadyFormatted(t *testing.T) {
	input := "if x > 1 {\n    print(\"yes\")\n}"
	assert.Equal(t, input, formatSource(input, 4, true))
}

func TestFormatPreservesEmptyLines(t *testing.T) {
	input := "x = 1\n\ny = 2"
	expected := "x = 1\n\ny = 2"
	assert.Equal(t, expected, formatSource(input, 4, true))
}

func TestFormatWithTabs(t *testing.T) {
	input := "if x > 1 {\nprint(\"yes\")\n}"
	expected := "if x > 1 {\n\tprint(\"yes\")\n}"
	assert.Equal(t, expected, formatSource(input, 4, false))
}

func TestFormatClosingBraceWithContent(t *testing.T) {
	input := "fn deploy() {\nexec(\"deploy.sh\")\n}\nprint(\"done\")"
	expected := "fn deploy() {\n    exec(\"deploy.sh\")\n}\nprint(\"done\")"
	assert.Equal(t, expected, formatSource(input, 4, true))
}
