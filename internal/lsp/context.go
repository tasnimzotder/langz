package lsp

import "github.com/tasnimzotder/langz/internal/lexer"

// findEnclosingFuncCall finds the function call that encloses the given
// cursor position (1-based line and col). Returns the function name and
// the parameter index (0-based, counting commas before cursor).
// Returns ("", -1) if the cursor is not inside any function call.
func findEnclosingFuncCall(source string, line, col int) (string, int) {
	tokens := lexer.New(source).Tokenize()

	// Find the index of the last token before or at the cursor
	cursorIdx := -1
	for i := range tokens {
		t := &tokens[i]
		if t.Type == lexer.EOF {
			break
		}
		if t.Line < line || (t.Line == line && t.Col < col) {
			cursorIdx = i
		}
	}
	if cursorIdx < 0 {
		return "", -1
	}

	// Walk backward from cursor, tracking paren depth
	depth := 0
	commas := 0
	for i := cursorIdx; i >= 0; i-- {
		t := tokens[i]
		switch t.Type {
		case lexer.RPAREN:
			depth++
		case lexer.LPAREN:
			if depth > 0 {
				depth--
			} else {
				// Found unmatched '(' — check if preceded by IDENT
				if i > 0 && tokens[i-1].Type == lexer.IDENT {
					return tokens[i-1].Value, commas
				}
				return "", -1
			}
		case lexer.COMMA:
			if depth == 0 {
				commas++
			}
		}
	}
	return "", -1
}
