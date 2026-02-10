package lsp

import (
	"fmt"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/tasnimzotder/langz/internal/lexer"
)

func (s *Server) textDocumentHover(ctx *glsp.Context, params *protocol.HoverParams) (result *protocol.Hover, err error) {
	defer recoverErr(&err)
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

	hover := getHoverAt(content, line, col)
	if hover == nil {
		return nil, nil
	}
	return hover, nil
}

// getHoverAt returns hover information for the token at the given 1-based position.
func getHoverAt(source string, line, col int) *protocol.Hover {
	token := findTokenAt(source, line, col)
	if token == nil || token.Type != lexer.IDENT {
		return nil
	}

	// Check builtin docs first
	if doc, ok := builtinDocs[token.Value]; ok {
		return makeHover(doc, token)
	}

	// Check method docs (IDENT preceded by DOT)
	if isMethod(source, token) {
		if doc, ok := methodDocs[token.Value]; ok {
			return makeHover(doc, token)
		}
	}

	// Check if this is a kwarg: IDENT followed by COLON
	if isKwarg(source, token) {
		funcName, _ := findEnclosingFuncCall(source, line, col)
		if kwargs, ok := builtinKwargs[funcName]; ok {
			for _, kw := range kwargs {
				if kw.Name == token.Value {
					doc := fmt.Sprintf("**%s:** %s", kw.Name, kw.Desc)
					return makeHover(doc, token)
				}
			}
		}
	}

	return nil
}

func makeHover(doc string, token *lexer.Token) *protocol.Hover {
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
	}
}

// isMethod checks if the given IDENT token is preceded by a DOT in the token stream.
func isMethod(source string, target *lexer.Token) bool {
	tokens := lexer.New(source).Tokenize()
	for i := range tokens {
		t := &tokens[i]
		if t.Line == target.Line && t.Col == target.Col && t.Type == lexer.IDENT {
			if i > 0 && tokens[i-1].Type == lexer.DOT {
				return true
			}
			return false
		}
	}
	return false
}

// isKwarg checks if the given IDENT token is followed by a COLON in the token stream.
func isKwarg(source string, target *lexer.Token) bool {
	tokens := lexer.New(source).Tokenize()
	for i := range tokens {
		t := &tokens[i]
		if t.Line == target.Line && t.Col == target.Col && t.Type == lexer.IDENT {
			if i+1 < len(tokens) && tokens[i+1].Type == lexer.COLON {
				return true
			}
			return false
		}
	}
	return false
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
