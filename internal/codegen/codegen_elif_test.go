package codegen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestElifCodegen(t *testing.T) {
	output := compile(`
x = 2
if x == 1 {
	print("one")
} else if x == 2 {
	print("two")
}
`)
	assert.Contains(t, output, "elif")
	assert.NotContains(t, output, "else\n  if")
}

func TestElifChainCodegen(t *testing.T) {
	output := compile(`
x = 3
if x == 1 {
	print("one")
} else if x == 2 {
	print("two")
} else if x == 3 {
	print("three")
} else {
	print("other")
}
`)
	assert.Contains(t, output, "elif [ \"$x\" = 2 ]; then")
	assert.Contains(t, output, "elif [ \"$x\" = 3 ]; then")
	assert.Contains(t, output, "else")
	// Should have exactly one fi for the whole chain
	assert.Equal(t, 1, countOccurrences(output, "fi"))
}

func countOccurrences(s, substr string) int {
	count := 0
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			// Make sure it's "fi" not a substring of another word
			if i+len(substr) >= len(s) || s[i+len(substr)] == '\n' || s[i+len(substr)] == ' ' {
				if i == 0 || s[i-1] == '\n' || s[i-1] == ' ' {
					count++
				}
			}
		}
	}
	return count
}
