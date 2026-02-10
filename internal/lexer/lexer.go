package lexer

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Lexer tokenizes Langz source code into a stream of tokens.
type Lexer struct {
	input   string
	pos     int
	current rune
	line    int
	col     int
}

// New creates a new Lexer for the given input source.
func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, col: 1}
	if len(input) > 0 {
		r, _ := utf8.DecodeRuneInString(input)
		l.current = r
	}
	return l
}

func (l *Lexer) advance() {
	if l.pos < len(l.input) && l.current == '\n' {
		l.line++
		l.col = 0
	}
	size := utf8.RuneLen(l.current)
	if size < 1 {
		size = 1
	}
	l.pos += size
	l.col++
	if l.pos < len(l.input) {
		r, _ := utf8.DecodeRuneInString(l.input[l.pos:])
		l.current = r
	} else {
		l.current = 0
	}
}

func (l *Lexer) token(t TokenType, value string, line, col int) Token {
	return Token{Type: t, Value: value, Line: line, Col: col}
}

func (l *Lexer) peekByte() byte {
	if l.pos+utf8.RuneLen(l.current) < len(l.input) {
		return l.input[l.pos+utf8.RuneLen(l.current)]
	}
	if l.pos+1 < len(l.input) {
		return l.input[l.pos+1]
	}
	return 0
}

func (l *Lexer) peekRune() rune {
	size := utf8.RuneLen(l.current)
	if size < 1 {
		size = 1
	}
	nextPos := l.pos + size
	if nextPos < len(l.input) {
		r, _ := utf8.DecodeRuneInString(l.input[nextPos:])
		return r
	}
	return 0
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) {
		if l.current == ' ' || l.current == '\t' || l.current == '\n' || l.current == '\r' {
			l.advance()
		} else if l.current == '/' && l.peekByte() == '/' {
			l.skipComment()
		} else {
			break
		}
	}
}

func (l *Lexer) skipComment() {
	for l.pos < len(l.input) && l.current != '\n' {
		l.advance()
	}
}

// readString reads a string literal. Returns the content and whether the string
// was properly terminated. An unterminated string returns (partial, false).
func (l *Lexer) readString() (string, bool) {
	l.advance() // skip opening "
	var buf []byte
	for l.pos < len(l.input) && l.current != '"' {
		if l.current == '\\' && l.pos+1 < len(l.input) {
			next := l.input[l.pos+1]
			switch next {
			case '"':
				buf = append(buf, '"')
			case 'n':
				buf = append(buf, '\n')
			case 't':
				buf = append(buf, '\t')
			case '\\':
				buf = append(buf, '\\')
			default:
				buf = append(buf, '\\', next)
			}
			l.advance() // skip '\'
			l.advance() // skip escaped char
		} else {
			buf = append(buf, l.input[l.pos])
			l.advance()
		}
	}
	if l.current != '"' {
		return string(buf), false
	}
	l.advance() // skip closing "
	return string(buf), true
}

func (l *Lexer) readIdent() string {
	start := l.pos
	for l.pos < len(l.input) && isAlphanumeric(l.current) {
		l.advance()
	}
	return l.input[start:l.pos]
}

func (l *Lexer) readNumber() string {
	start := l.pos
	for l.pos < len(l.input) && isDigit(l.current) {
		l.advance()
	}
	return l.input[start:l.pos]
}

// readBashContent reads raw content inside a bash { } block,
// tracking brace depth and respecting string literals.
func (l *Lexer) readBashContent() string {
	depth := 1
	start := l.pos
	for l.pos < len(l.input) && depth > 0 {
		switch l.current {
		case '{':
			depth++
			l.advance()
		case '}':
			depth--
			if depth == 0 {
				content := l.input[start:l.pos]
				l.advance() // skip closing }
				return strings.TrimSpace(content)
			}
			l.advance()
		case '"', '\'':
			// Skip string literals to avoid counting braces inside them
			quote := l.current
			l.advance()
			for l.pos < len(l.input) && l.current != quote {
				if l.current == '\\' {
					l.advance() // skip escape
				}
				l.advance()
			}
			if l.pos < len(l.input) {
				l.advance() // skip closing quote
			}
		case '#':
			// Skip comments to end of line
			for l.pos < len(l.input) && l.current != '\n' {
				l.advance()
			}
		default:
			l.advance()
		}
	}
	// Unterminated — return what we have
	return strings.TrimSpace(l.input[start:l.pos])
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch)
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func isAlphanumeric(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_'
}

// Tokenize scans the entire input and returns a slice of tokens.
func (l *Lexer) Tokenize() []Token {
	var tokens []Token

	// Skip shebang line (e.g. #!/usr/bin/env langz)
	if l.pos < len(l.input) && l.current == '#' && l.peekByte() == '!' {
		for l.pos < len(l.input) && l.current != '\n' {
			l.advance()
		}
	}

	for l.pos < len(l.input) {
		l.skipWhitespace()

		if l.pos >= len(l.input) {
			break
		}

		line, col := l.line, l.col

		switch {
		case l.current == '=':
			if l.peekByte() == '>' {
				tokens = append(tokens, l.token(FATARROW, "=>", line, col))
				l.advance()
			} else if l.peekByte() == '=' {
				tokens = append(tokens, l.token(EQ, "==", line, col))
				l.advance()
			} else {
				tokens = append(tokens, l.token(ASSIGN, "=", line, col))
			}
			l.advance()
		case l.current == '!':
			if l.peekByte() == '=' {
				tokens = append(tokens, l.token(NEQ, "!=", line, col))
				l.advance()
			} else {
				tokens = append(tokens, l.token(BANG, "!", line, col))
			}
			l.advance()
		case l.current == '-' && l.peekByte() == '>':
			tokens = append(tokens, l.token(ARROW, "->", line, col))
			l.advance()
			l.advance()
		case l.current == '>':
			if l.peekByte() == '=' {
				tokens = append(tokens, l.token(GTE, ">=", line, col))
				l.advance()
			} else {
				tokens = append(tokens, l.token(GT, ">", line, col))
			}
			l.advance()
		case l.current == '<':
			if l.peekByte() == '=' {
				tokens = append(tokens, l.token(LTE, "<=", line, col))
				l.advance()
			} else {
				tokens = append(tokens, l.token(LT, "<", line, col))
			}
			l.advance()
		case l.current == '+':
			if l.peekByte() == '=' {
				tokens = append(tokens, l.token(PLUS_ASSIGN, "+=", line, col))
				l.advance()
			} else {
				tokens = append(tokens, l.token(PLUS, "+", line, col))
			}
			l.advance()
		case l.current == '-':
			if l.peekByte() == '=' {
				tokens = append(tokens, l.token(MINUS_ASSIGN, "-=", line, col))
				l.advance()
			} else {
				tokens = append(tokens, l.token(MINUS, "-", line, col))
			}
			l.advance()
		case l.current == '*':
			if l.peekByte() == '=' {
				tokens = append(tokens, l.token(STAR_ASSIGN, "*=", line, col))
				l.advance()
			} else {
				tokens = append(tokens, l.token(STAR, "*", line, col))
			}
			l.advance()
		case l.current == '/':
			if l.peekByte() == '=' {
				tokens = append(tokens, l.token(SLASH_ASSIGN, "/=", line, col))
				l.advance()
			} else {
				tokens = append(tokens, l.token(SLASH, "/", line, col))
			}
			l.advance()
		case l.current == '%':
			tokens = append(tokens, l.token(PERCENT, "%", line, col))
			l.advance()
		case l.current == '(':
			tokens = append(tokens, l.token(LPAREN, "(", line, col))
			l.advance()
		case l.current == ')':
			tokens = append(tokens, l.token(RPAREN, ")", line, col))
			l.advance()
		case l.current == '{':
			tokens = append(tokens, l.token(LBRACE, "{", line, col))
			l.advance()
		case l.current == '}':
			tokens = append(tokens, l.token(RBRACE, "}", line, col))
			l.advance()
		case l.current == '[':
			tokens = append(tokens, l.token(LBRACKET, "[", line, col))
			l.advance()
		case l.current == ']':
			tokens = append(tokens, l.token(RBRACKET, "]", line, col))
			l.advance()
		case l.current == ',':
			tokens = append(tokens, l.token(COMMA, ",", line, col))
			l.advance()
		case l.current == ':':
			tokens = append(tokens, l.token(COLON, ":", line, col))
			l.advance()
		case l.current == '|' && l.peekByte() == '>':
			tokens = append(tokens, l.token(PIPE, "|>", line, col))
			l.advance()
			l.advance()
		case l.current == '.':
			tokens = append(tokens, l.token(DOT, ".", line, col))
			l.advance()
		case l.current == '"':
			str, ok := l.readString()
			if ok {
				tokens = append(tokens, l.token(STRING, str, line, col))
			} else {
				tokens = append(tokens, l.token(ILLEGAL, "unterminated string", line, col))
			}
		case isDigit(l.current):
			tokens = append(tokens, l.token(INT, l.readNumber(), line, col))
		case l.current == '_' && !isAlphanumeric(l.peekRune()):
			tokens = append(tokens, l.token(UNDERSCORE, "_", line, col))
			l.advance()
		case isLetter(l.current) || l.current == '_':
			word := l.readIdent()
			if kwType, ok := keywords[word]; ok {
				tokens = append(tokens, l.token(kwType, word, line, col))
				// After BASH keyword, capture raw content between { }
				if kwType == BASH {
					l.skipWhitespace()
					if l.pos < len(l.input) && l.current == '{' {
						l.advance() // skip opening {
						content := l.readBashContent()
						cLine, cCol := l.line, l.col
						tokens = append(tokens, l.token(BASH_CONTENT, content, cLine, cCol))
					}
				}
			} else {
				tokens = append(tokens, l.token(IDENT, word, line, col))
			}
		default:
			tokens = append(tokens, l.token(ILLEGAL, string(l.current), line, col))
			l.advance()
		}
	}

	tokens = append(tokens, l.token(EOF, "", l.line, l.col))
	return tokens
}
