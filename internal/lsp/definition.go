package lsp

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/tasnimzotder/langz/internal/lexer"
)

func (s *Server) textDocumentDefinition(ctx *glsp.Context, params *protocol.DefinitionParams) (any, error) {
	uri := params.TextDocument.URI
	content, ok := s.documents[uri]
	if !ok {
		return nil, nil
	}

	line := int(params.Position.Line) + 1
	col := int(params.Position.Character) + 1

	token := findTokenAt(content, line, col)
	if token == nil || token.Type != lexer.IDENT {
		return nil, nil
	}

	// Builtins have no source definition
	if _, isBuiltin := builtinDocs[token.Value]; isBuiltin {
		return nil, nil
	}

	loc := getDefinition(content, token.Value)
	if loc == nil {
		return nil, nil
	}
	loc.URI = uri
	return loc, nil
}

// getDefinition finds the first definition of name in source.
// Returns a Location with 0-based positions (URI left empty for caller to fill).
func getDefinition(source, name string) *protocol.Location {
	for _, sym := range findSymbols(source) {
		if sym.Name == name {
			return &protocol.Location{
				Range: protocol.Range{
					Start: protocol.Position{
						Line:      protocol.UInteger(sym.Line - 1),
						Character: protocol.UInteger(sym.Col - 1),
					},
					End: protocol.Position{
						Line:      protocol.UInteger(sym.Line - 1),
						Character: protocol.UInteger(sym.Col - 1 + len(sym.Name)),
					},
				},
			}
		}
	}
	return nil
}
