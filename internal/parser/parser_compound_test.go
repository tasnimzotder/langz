package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tasnimzotder/langz/internal/ast"
)

func TestParseCompoundAssign(t *testing.T) {
	tests := []struct {
		input string
		op    string
	}{
		{`x += 5`, "+"},
		{`x -= 3`, "-"},
		{`x *= 2`, "*"},
		{`x /= 4`, "/"},
	}

	for _, tt := range tests {
		t.Run(tt.op, func(t *testing.T) {
			prog := parse(tt.input)
			require.Len(t, prog.Statements, 1)
			assign := prog.Statements[0].(*ast.Assignment)
			assert.Equal(t, "x", assign.Name)

			// Desugared: x = x op value
			bin, ok := assign.Value.(*ast.BinaryExpr)
			require.True(t, ok, "expected BinaryExpr")
			assert.Equal(t, tt.op, bin.Op)

			left, ok := bin.Left.(*ast.Identifier)
			require.True(t, ok)
			assert.Equal(t, "x", left.Name)
		})
	}
}
