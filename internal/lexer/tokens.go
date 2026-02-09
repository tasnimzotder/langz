package lexer

// TokenType represents the type of a lexer token.
type TokenType string

const (
	// Literals
	IDENT  TokenType = "IDENT"
	INT    TokenType = "INT"
	STRING TokenType = "STRING"

	// Operators & punctuation
	ASSIGN     TokenType = "ASSIGN"     // =
	EQ         TokenType = "EQ"         // ==
	NEQ        TokenType = "NEQ"        // !=
	GT         TokenType = "GT"         // >
	GTE        TokenType = "GTE"        // >=
	LT         TokenType = "LT"         // <
	LTE        TokenType = "LTE"        // <=
	BANG       TokenType = "BANG"       // !
	PLUS       TokenType = "PLUS"       // +
	MINUS      TokenType = "MINUS"      // -
	STAR       TokenType = "STAR"       // *
	SLASH      TokenType = "SLASH"      // /
	PERCENT    TokenType = "PERCENT"    // %
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
	AND      TokenType = "AND"
	MATCH    TokenType = "MATCH"
	TRUE     TokenType = "TRUE"
	FALSE    TokenType = "FALSE"
	CONTINUE TokenType = "CONTINUE"
	BREAK    TokenType = "BREAK"
	WHILE    TokenType = "WHILE"

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
	"and":      AND,
	"match":    MATCH,
	"true":     TRUE,
	"false":    FALSE,
	"continue": CONTINUE,
	"break":    BREAK,
	"while":    WHILE,
}

// Token represents a single lexical token with position information.
type Token struct {
	Type  TokenType
	Value string
	Line  int
	Col   int
}
