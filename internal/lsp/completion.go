package lsp

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/tasnimzotder/langz/internal/lexer"
)

func (s *Server) textDocumentCompletion(ctx *glsp.Context, params *protocol.CompletionParams) (any, error) {
	uri := params.TextDocument.URI
	content := s.documents[uri]
	// LSP positions are 0-based; our token positions are 1-based
	line := int(params.Position.Line) + 1
	col := int(params.Position.Character) + 1
	return getCompletionItems(content, line, col), nil
}

func getCompletionItems(source string, line, col int) []protocol.CompletionItem {
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

	// Context-aware kwargs — only when cursor is inside a function call
	if funcName, _ := findEnclosingFuncCall(source, line, col); funcName != "" {
		if kwargs, ok := builtinKwargs[funcName]; ok {
			for _, kw := range kwargs {
				label := kw.Name + ":"
				if seen[label] {
					continue
				}
				seen[label] = true
				kind := protocol.CompletionItemKindProperty
				items = append(items, protocol.CompletionItem{
					Label:         label,
					Kind:          &kind,
					Documentation: kw.Desc,
				})
			}
		}
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
