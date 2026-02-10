package lexer

// Lexer tokenizes Langz source code into a stream of tokens.
type Lexer struct {
	input   string
	pos     int
	current byte
	line    int
	col     int
}

// New creates a new Lexer for the given input source.
func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, col: 1}
	if len(input) > 0 {
		l.current = input[0]
	}
	return l
}

func (l *Lexer) advance() {
	if l.pos < len(l.input) && l.current == '\n' {
		l.line++
		l.col = 0
	}
	l.pos++
	l.col++
	if l.pos < len(l.input) {
		l.current = l.input[l.pos]
	} else {
		l.current = 0
	}
}

func (l *Lexer) token(t TokenType, value string, line, col int) Token {
	return Token{Type: t, Value: value, Line: line, Col: col}
}

func (l *Lexer) peek() byte {
	if l.pos+1 < len(l.input) {
		return l.input[l.pos+1]
	}
	return 0
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) {
		if l.current == ' ' || l.current == '\t' || l.current == '\n' || l.current == '\r' {
			l.advance()
		} else if l.current == '/' && l.peek() == '/' {
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

func (l *Lexer) readString() string {
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
			buf = append(buf, l.current)
			l.advance()
		}
	}
	l.advance() // skip closing "
	return string(buf)
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

func isLetter(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isAlphanumeric(ch byte) bool {
	return isLetter(ch) || isDigit(ch) || ch == '_'
}

// Tokenize scans the entire input and returns a slice of tokens.
func (l *Lexer) Tokenize() []Token {
	var tokens []Token

	for l.pos < len(l.input) {
		l.skipWhitespace()

		if l.pos >= len(l.input) {
			break
		}

		line, col := l.line, l.col

		switch {
		case l.current == '=':
			if l.peek() == '>' {
				tokens = append(tokens, l.token(FATARROW, "=>", line, col))
				l.advance()
			} else if l.peek() == '=' {
				tokens = append(tokens, l.token(EQ, "==", line, col))
				l.advance()
			} else {
				tokens = append(tokens, l.token(ASSIGN, "=", line, col))
			}
			l.advance()
		case l.current == '!':
			if l.peek() == '=' {
				tokens = append(tokens, l.token(NEQ, "!=", line, col))
				l.advance()
			} else {
				tokens = append(tokens, l.token(BANG, "!", line, col))
			}
			l.advance()
		case l.current == '-' && l.peek() == '>':
			tokens = append(tokens, l.token(ARROW, "->", line, col))
			l.advance()
			l.advance()
		case l.current == '>':
			if l.peek() == '=' {
				tokens = append(tokens, l.token(GTE, ">=", line, col))
				l.advance()
			} else {
				tokens = append(tokens, l.token(GT, ">", line, col))
			}
			l.advance()
		case l.current == '<':
			if l.peek() == '=' {
				tokens = append(tokens, l.token(LTE, "<=", line, col))
				l.advance()
			} else {
				tokens = append(tokens, l.token(LT, "<", line, col))
			}
			l.advance()
		case l.current == '+':
			tokens = append(tokens, l.token(PLUS, "+", line, col))
			l.advance()
		case l.current == '-':
			tokens = append(tokens, l.token(MINUS, "-", line, col))
			l.advance()
		case l.current == '*':
			tokens = append(tokens, l.token(STAR, "*", line, col))
			l.advance()
		case l.current == '/':
			tokens = append(tokens, l.token(SLASH, "/", line, col))
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
		case l.current == '|' && l.peek() == '>':
			tokens = append(tokens, l.token(PIPE, "|>", line, col))
			l.advance()
			l.advance()
		case l.current == '.':
			tokens = append(tokens, l.token(DOT, ".", line, col))
			l.advance()
		case l.current == '"':
			tokens = append(tokens, l.token(STRING, l.readString(), line, col))
		case isDigit(l.current):
			tokens = append(tokens, l.token(INT, l.readNumber(), line, col))
		case l.current == '_' && !isAlphanumeric(l.peek()):
			tokens = append(tokens, l.token(UNDERSCORE, "_", line, col))
			l.advance()
		case isLetter(l.current) || l.current == '_':
			word := l.readIdent()
			if kwType, ok := keywords[word]; ok {
				tokens = append(tokens, l.token(kwType, word, line, col))
			} else {
				tokens = append(tokens, l.token(IDENT, word, line, col))
			}
		default:
			l.advance()
		}
	}

	tokens = append(tokens, l.token(EOF, "", l.line, l.col))
	return tokens
}
