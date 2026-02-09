package lsp

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/tasnimzotder/langz/internal/lexer"
)

func (s *Server) textDocumentHover(ctx *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	uri := params.TextDocument.URI
	content, ok := s.documents[uri]
	if !ok {
		return nil, nil
	}

	// LSP positions are 0-based; token positions are 1-based
	line := int(params.Position.Line) + 1
	col := int(params.Position.Character) + 1

	token := findTokenAt(content, line, col)
	if token == nil {
		return nil, nil
	}

	if token.Type != lexer.IDENT {
		return nil, nil
	}

	doc, ok := builtinDocs[token.Value]
	if !ok {
		return nil, nil
	}

	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: doc,
		},
		Range: &protocol.Range{
			Start: protocol.Position{
				Line:      protocol.UInteger(token.Line - 1),
				Character: protocol.UInteger(token.Col - 1),
			},
			End: protocol.Position{
				Line:      protocol.UInteger(token.Line - 1),
				Character: protocol.UInteger(token.Col - 1 + len(token.Value)),
			},
		},
	}, nil
}

// findTokenAt returns the token at the given 1-based line and column,
// or nil if no token spans that position.
func findTokenAt(source string, line, col int) *lexer.Token {
	tokens := lexer.New(source).Tokenize()
	for i := range tokens {
		t := &tokens[i]
		if t.Type == lexer.EOF {
			break
		}
		if t.Line == line {
			tokenEnd := t.Col + len(t.Value)
			if t.Type == lexer.STRING {
				tokenEnd = t.Col + len(t.Value) + 2 // account for quotes
			}
			if col >= t.Col && col < tokenEnd {
				return t
			}
		}
	}
	return nil
}
