package lexer

type TokenType string

const (
	// Literals
	IDENT  TokenType = "IDENT"
	INT    TokenType = "INT"
	STRING TokenType = "STRING"

	// Operators & punctuation
	ASSIGN     TokenType = "ASSIGN"     // =
	GT         TokenType = "GT"         // >
	BANG       TokenType = "BANG"       // !
	LPAREN     TokenType = "LPAREN"     // (
	RPAREN     TokenType = "RPAREN"     // )
	LBRACE     TokenType = "LBRACE"     // {
	RBRACE     TokenType = "RBRACE"     // }
	LBRACKET   TokenType = "LBRACKET"   // [
	RBRACKET   TokenType = "RBRACKET"   // ]
	COMMA      TokenType = "COMMA"      // ,
	COLON      TokenType = "COLON"      // :
	DOT        TokenType = "DOT"        // .
	ARROW      TokenType = "ARROW"      // ->
	FATARROW   TokenType = "FATARROW"   // =>
	UNDERSCORE TokenType = "UNDERSCORE" // _

	// Keywords
	IF       TokenType = "IF"
	ELSE     TokenType = "ELSE"
	FOR      TokenType = "FOR"
	IN       TokenType = "IN"
	FN       TokenType = "FN"
	RETURN   TokenType = "RETURN"
	OR       TokenType = "OR"
	MATCH    TokenType = "MATCH"
	TRUE     TokenType = "TRUE"
	FALSE    TokenType = "FALSE"
	CONTINUE TokenType = "CONTINUE"

	EOF TokenType = "EOF"
)

var keywords = map[string]TokenType{
	"if":       IF,
	"else":     ELSE,
	"for":      FOR,
	"in":       IN,
	"fn":       FN,
	"return":   RETURN,
	"or":       OR,
	"match":    MATCH,
	"true":     TRUE,
	"false":    FALSE,
	"continue": CONTINUE,
}

type Token struct {
	Type  TokenType
	Value string
	Line  int
	Col   int
}

type Lexer struct {
	input   string
	pos     int
	current byte
	line    int
	col     int
}

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
	for l.pos < len(l.input) && (l.current == ' ' || l.current == '\t' || l.current == '\n' || l.current == '\r') {
		l.advance()
	}
}

func (l *Lexer) readString() string {
	l.advance() // skip opening "
	start := l.pos
	for l.pos < len(l.input) && l.current != '"' {
		l.advance()
	}
	value := l.input[start:l.pos]
	l.advance() // skip closing "
	return value
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
			} else {
				tokens = append(tokens, l.token(ASSIGN, "=", line, col))
			}
			l.advance()
		case l.current == '-' && l.peek() == '>':
			tokens = append(tokens, l.token(ARROW, "->", line, col))
			l.advance()
			l.advance()
		case l.current == '>':
			tokens = append(tokens, l.token(GT, ">", line, col))
			l.advance()
		case l.current == '!':
			tokens = append(tokens, l.token(BANG, "!", line, col))
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
