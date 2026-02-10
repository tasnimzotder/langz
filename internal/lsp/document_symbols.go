package lsp

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s *Server) textDocumentDocumentSymbol(ctx *glsp.Context, params *protocol.DocumentSymbolParams) (result any, err error) {
	defer recoverErr(&err)
	uri := params.TextDocument.URI
	content, ok := s.documents[uri]
	if !ok {
		return nil, nil
	}
	return getDocumentSymbols(content), nil
}

func getDocumentSymbols(source string) []protocol.DocumentSymbol {
	symbols := findSymbols(source)
	result := make([]protocol.DocumentSymbol, 0, len(symbols))

	for _, sym := range symbols {
		kind := protocol.SymbolKindVariable
		if sym.Kind == "function" {
			kind = protocol.SymbolKindFunction
		}

		nameRange := protocol.Range{
			Start: protocol.Position{
				Line:      protocol.UInteger(sym.Line - 1),
				Character: protocol.UInteger(sym.Col - 1),
			},
			End: protocol.Position{
				Line:      protocol.UInteger(sym.Line - 1),
				Character: protocol.UInteger(sym.Col - 1 + len(sym.Name)),
			},
		}

		result = append(result, protocol.DocumentSymbol{
			Name:           sym.Name,
			Kind:           kind,
			Range:          nameRange,
			SelectionRange: nameRange,
		})
	}

	return result
}
