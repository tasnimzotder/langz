package integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tasnimzotder/langz/internal/codegen"
	"github.com/tasnimzotder/langz/internal/lexer"
	"github.com/tasnimzotder/langz/internal/parser"
)

func TestE2E_UnterminatedStringProducesError(t *testing.T) {
	source := `x = "unterminated`
	tokens := lexer.New(source).Tokenize()
	_, err := parser.New(tokens).ParseWithErrors()

	require.Error(t, err, "unterminated string should produce a parse error")
	assert.Contains(t, err.Error(), "unterminated string")
}

func TestE2E_UnknownCharacterProducesError(t *testing.T) {
	source := `x = @value`
	tokens := lexer.New(source).Tokenize()
	_, err := parser.New(tokens).ParseWithErrors()

	require.Error(t, err, "unknown character should produce a parse error")
}

func TestE2E_CodegenErrorCausesNonEmptyErrors(t *testing.T) {
	// write() with wrong arity should produce a codegen error
	source := `write("file.txt")`
	tokens := lexer.New(source).Tokenize()
	prog, err := parser.New(tokens).ParseWithErrors()
	require.NoError(t, err, "should parse without error")

	_, errs := codegen.Generate(prog)
	require.NotEmpty(t, errs, "codegen should report write() arity error")
	assert.Contains(t, errs[0], "write()")
}

func TestE2E_EmptyFileCompiles(t *testing.T) {
	source := ""
	tokens := lexer.New(source).Tokenize()
	prog, err := parser.New(tokens).ParseWithErrors()
	require.NoError(t, err, "empty file should parse without error")

	output, errs := codegen.Generate(prog)
	assert.Empty(t, errs, "empty file should not produce codegen errors")
	assert.Contains(t, output, "#!/bin/bash")
}

func TestE2E_CommentOnlyFileCompiles(t *testing.T) {
	source := "// this is a comment\n// another comment"
	tokens := lexer.New(source).Tokenize()
	prog, err := parser.New(tokens).ParseWithErrors()
	require.NoError(t, err, "comment-only file should parse without error")

	output, errs := codegen.Generate(prog)
	assert.Empty(t, errs, "comment-only file should not produce codegen errors")
	assert.Contains(t, output, "#!/bin/bash")
}

func TestE2E_EmptyFileRunsSuccessfully(t *testing.T) {
	source := ""
	bash := compileSource(t, source)
	_, code := runBash(t, bash)
	assert.Equal(t, 0, code, "empty compiled script should exit 0")
}

func TestE2E_CommentOnlyRunsSuccessfully(t *testing.T) {
	source := "// just comments\n// nothing else"
	bash := compileSource(t, source)
	_, code := runBash(t, bash)
	assert.Equal(t, 0, code, "comment-only compiled script should exit 0")
}
