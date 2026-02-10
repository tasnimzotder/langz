package lsp

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/tasnimzotder/langz/internal/lexer"
	"github.com/tasnimzotder/langz/internal/parser"
)

// publishDiagnostics uses cached tokens to parse and send diagnostics to the client.
func (s *Server) publishDiagnostics(ctx *glsp.Context, uri protocol.DocumentUri, _ string) {
	tokens := s.getTokens(uri)
	diags := getDiagnosticsFromTokens(tokens)
	ctx.Notify(protocol.ServerTextDocumentPublishDiagnostics, &protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diags,
	})
}

// getDiagnostics lexes + parses source, returns LSP diagnostics for all errors.
func getDiagnostics(source string) []protocol.Diagnostic {
	tokens := lexer.New(source).Tokenize()
	return getDiagnosticsFromTokens(tokens)
}

// getDiagnosticsFromTokens parses pre-tokenized input and returns LSP diagnostics.
func getDiagnosticsFromTokens(tokens []lexer.Token) []protocol.Diagnostic {
	_, errs := parser.New(tokens).ParseAllErrors()

	diags := make([]protocol.Diagnostic, 0, len(errs))
	severity := protocol.DiagnosticSeverityError
	sourceName := "langz"

	for _, e := range errs {
		// Parser positions are 1-based; LSP positions are 0-based
		line := protocol.UInteger(0)
		if e.Line > 0 {
			line = protocol.UInteger(e.Line - 1)
		}
		col := protocol.UInteger(0)
		if e.Col > 0 {
			col = protocol.UInteger(e.Col - 1)
		}

		diags = append(diags, protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{Line: line, Character: col},
				End:   protocol.Position{Line: line, Character: col + 1},
			},
			Severity: &severity,
			Source:   &sourceName,
			Message:  e.Message,
		})
	}
	return diags
}
