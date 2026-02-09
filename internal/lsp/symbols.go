package lsp

import (
	"github.com/tasnimzotder/langz/internal/lexer"
)

type symbolInfo struct {
	Name string
	Kind string // "variable", "function", "for_var"
	Line int    // 1-based
	Col  int    // 1-based
}

// findSymbols scans the token stream for definition patterns and returns
// all symbols with their positions. This avoids needing AST position info.
func findSymbols(source string) []symbolInfo {
	tokens := lexer.New(source).Tokenize()
	var symbols []symbolInfo

	for i := 0; i < len(tokens); i++ {
		t := tokens[i]
		switch {
		case t.Type == lexer.IDENT && i+1 < len(tokens) && tokens[i+1].Type == lexer.ASSIGN:
			symbols = append(symbols, symbolInfo{
				Name: t.Value,
				Kind: "variable",
				Line: t.Line,
				Col:  t.Col,
			})
		case t.Type == lexer.FN && i+1 < len(tokens) && tokens[i+1].Type == lexer.IDENT:
			next := tokens[i+1]
			symbols = append(symbols, symbolInfo{
				Name: next.Value,
				Kind: "function",
				Line: next.Line,
				Col:  next.Col,
			})
		case t.Type == lexer.FOR && i+1 < len(tokens) && tokens[i+1].Type == lexer.IDENT:
			next := tokens[i+1]
			symbols = append(symbols, symbolInfo{
				Name: next.Value,
				Kind: "for_var",
				Line: next.Line,
				Col:  next.Col,
			})
		}
	}
	return symbols
}
