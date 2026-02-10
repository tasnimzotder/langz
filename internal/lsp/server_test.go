package lsp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/tasnimzotder/langz/internal/lexer"
)

// --- Diagnostics ---

func TestGetDiagnosticsNoErrors(t *testing.T) {
	diags := getDiagnostics(`x = 1`)
	assert.Empty(t, diags)
}

func TestGetDiagnosticsWithError(t *testing.T) {
	diags := getDiagnostics(`fn (`)
	require.GreaterOrEqual(t, len(diags), 1)
	assert.Contains(t, diags[0].Message, "expected IDENT")
	assert.Equal(t, protocol.DiagnosticSeverityError, *diags[0].Severity)
}

func TestGetDiagnosticsPositionConversion(t *testing.T) {
	// Parser positions are 1-based, LSP positions are 0-based
	diags := getDiagnostics(`fn (`)
	require.GreaterOrEqual(t, len(diags), 1)
	// "(" is at parser col 4 (1-based) -> LSP character 3 (0-based)
	assert.Equal(t, protocol.UInteger(0), diags[0].Range.Start.Line)
	assert.Equal(t, protocol.UInteger(3), diags[0].Range.Start.Character)
}

func TestGetDiagnosticsValidProgram(t *testing.T) {
	source := `
name = "hello"
print(name)
if x > 10 {
	print("big")
}
`
	diags := getDiagnostics(source)
	assert.Empty(t, diags)
}

// --- Token lookup ---

func TestFindTokenAtBuiltin(t *testing.T) {
	tok := findTokenAt(`print("hello")`, 1, 1)
	require.NotNil(t, tok)
	assert.Equal(t, lexer.IDENT, tok.Type)
	assert.Equal(t, "print", tok.Value)
}

func TestFindTokenAtMiddleOfToken(t *testing.T) {
	tok := findTokenAt(`print("hello")`, 1, 3)
	require.NotNil(t, tok)
	assert.Equal(t, "print", tok.Value)
}

func TestFindTokenAtEndOfToken(t *testing.T) {
	tok := findTokenAt(`print("hello")`, 1, 5)
	require.NotNil(t, tok)
	assert.Equal(t, "print", tok.Value)
}

func TestFindTokenAtNoToken(t *testing.T) {
	tok := findTokenAt(`x = 1`, 5, 1)
	assert.Nil(t, tok)
}

func TestFindTokenAtString(t *testing.T) {
	tok := findTokenAt(`print("hello")`, 1, 7)
	require.NotNil(t, tok)
	assert.Equal(t, lexer.STRING, tok.Type)
	assert.Equal(t, "hello", tok.Value)
}

// --- Builtin docs ---

func TestBuiltinDocsCompleteness(t *testing.T) {
	expected := []string{
		"print", "write", "append", "read", "rm", "rmdir", "mkdir",
		"copy", "move", "chmod", "glob", "exists", "is_file", "is_dir",
		"exec", "exit", "env", "os", "arch", "hostname", "whoami",
		"dirname", "basename", "upper", "lower", "fetch", "args",
		"range", "sleep", "json_get", "len", "trim", "timestamp",
		"date", "chown",
	}
	for _, name := range expected {
		_, ok := builtinDocs[name]
		assert.True(t, ok, "missing docs for builtin: %s", name)
	}
}

func TestBuiltinDocsContainSignature(t *testing.T) {
	// Each doc should contain the function name
	for name, doc := range builtinDocs {
		assert.Contains(t, doc, name, "doc for %s should contain the function name", name)
	}
}

// --- Kwarg docs ---

func TestBuiltinKwargsHasEntries(t *testing.T) {
	kwargs, ok := builtinKwargs["fetch"]
	require.True(t, ok, "fetch should have kwargs")
	assert.Len(t, kwargs, 5)

	names := make([]string, len(kwargs))
	for i, kw := range kwargs {
		names[i] = kw.Name
	}
	assert.Contains(t, names, "method")
	assert.Contains(t, names, "body")
	assert.Contains(t, names, "headers")
	assert.Contains(t, names, "timeout")
	assert.Contains(t, names, "retries")
}

// --- Server ---

func TestNewServerHasEmptyDocuments(t *testing.T) {
	s := NewServer()
	assert.Empty(t, s.documents)
}
