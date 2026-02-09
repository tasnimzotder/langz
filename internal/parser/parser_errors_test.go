package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tasnimzotder/langz/internal/lexer"
)

func TestParseAllErrorsSingle(t *testing.T) {
	tokens := lexer.New("fn (").Tokenize()
	p := New(tokens)
	prog, errs := p.ParseAllErrors()

	assert.NotNil(t, prog)
	require.GreaterOrEqual(t, len(errs), 1)
	assert.Equal(t, 1, errs[0].Line)
	assert.Contains(t, errs[0].Message, "expected IDENT")
}

func TestParseAllErrorsNone(t *testing.T) {
	tokens := lexer.New(`x = 1`).Tokenize()
	p := New(tokens)
	prog, errs := p.ParseAllErrors()

	require.Len(t, prog.Statements, 1)
	assert.Empty(t, errs)
}

func TestParseWithErrorsBackwardsCompat(t *testing.T) {
	tokens := lexer.New("fn (").Tokenize()
	p := New(tokens)
	_, err := p.ParseWithErrors()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "expected IDENT")
	assert.Contains(t, err.Error(), "line")
	assert.Contains(t, err.Error(), "col")
}

func TestParseErrorPosition(t *testing.T) {
	tokens := lexer.New("fn (").Tokenize()
	p := New(tokens)
	_, errs := p.ParseAllErrors()

	require.GreaterOrEqual(t, len(errs), 1)
	assert.Greater(t, errs[0].Line, 0)
	assert.Greater(t, errs[0].Col, 0)
}
