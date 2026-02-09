package lsp

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/tasnimzotder/langz/internal/lexer"
)

func (s *Server) textDocumentCompletion(ctx *glsp.Context, params *protocol.CompletionParams) (any, error) {
	uri := params.TextDocument.URI
	content := s.documents[uri]
	return getCompletionItems(content), nil
}

func getCompletionItems(source string) []protocol.CompletionItem {
	var items []protocol.CompletionItem
	seen := make(map[string]bool)

	// Builtins
	for name, doc := range builtinDocs {
		seen[name] = true
		kind := protocol.CompletionItemKindFunction
		detail := name + "(...)"
		items = append(items, protocol.CompletionItem{
			Label:         name,
			Kind:          &kind,
			Detail:        &detail,
			Documentation: doc,
		})
	}

	// Keywords
	for _, kw := range lexer.KeywordNames() {
		if seen[kw] {
			continue
		}
		seen[kw] = true
		kind := protocol.CompletionItemKindKeyword
		items = append(items, protocol.CompletionItem{
			Label: kw,
			Kind:  &kind,
		})
	}

	// User-defined symbols
	for _, sym := range findSymbols(source) {
		if seen[sym.Name] {
			continue
		}
		seen[sym.Name] = true
		kind := protocol.CompletionItemKindVariable
		if sym.Kind == "function" {
			kind = protocol.CompletionItemKindFunction
		}
		items = append(items, protocol.CompletionItem{
			Label: sym.Name,
			Kind:  &kind,
		})
	}

	return items
}
