package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertTokens(t *testing.T, input string, expected []Token) {
	t.Helper()
	tokens := New(input).Tokenize()
	expected = append(expected, Token{Type: EOF, Value: ""})

	require.Len(t, tokens, len(expected))
	for i, tok := range tokens {
		assert.Equal(t, expected[i].Type, tok.Type, "token[%d] type", i)
		assert.Equal(t, expected[i].Value, tok.Value, "token[%d] value", i)
	}
}

func TestAssignment(t *testing.T) {
	assertTokens(t, `name = "hello"`, []Token{
		{Type: IDENT, Value: "name"},
		{Type: ASSIGN, Value: "="},
		{Type: STRING, Value: "hello"},
	})
}

func TestIntegerLiteral(t *testing.T) {
	assertTokens(t, `count = 42`, []Token{
		{Type: IDENT, Value: "count"},
		{Type: ASSIGN, Value: "="},
		{Type: INT, Value: "42"},
	})
}

func TestKeywords(t *testing.T) {
	assertTokens(t, `if true { return }`, []Token{
		{Type: IF, Value: "if"},
		{Type: TRUE, Value: "true"},
		{Type: LBRACE, Value: "{"},
		{Type: RETURN, Value: "return"},
		{Type: RBRACE, Value: "}"},
	})
}

func TestFunctionDeclaration(t *testing.T) {
	assertTokens(t, `fn greet(name: str) {`, []Token{
		{Type: FN, Value: "fn"},
		{Type: IDENT, Value: "greet"},
		{Type: LPAREN, Value: "("},
		{Type: IDENT, Value: "name"},
		{Type: COLON, Value: ":"},
		{Type: IDENT, Value: "str"},
		{Type: RPAREN, Value: ")"},
		{Type: LBRACE, Value: "{"},
	})
}

func TestForLoop(t *testing.T) {
	assertTokens(t, `for f in files {`, []Token{
		{Type: FOR, Value: "for"},
		{Type: IDENT, Value: "f"},
		{Type: IN, Value: "in"},
		{Type: IDENT, Value: "files"},
		{Type: LBRACE, Value: "{"},
	})
}

func TestComparison(t *testing.T) {
	assertTokens(t, `if x > 10 {`, []Token{
		{Type: IF, Value: "if"},
		{Type: IDENT, Value: "x"},
		{Type: GT, Value: ">"},
		{Type: INT, Value: "10"},
		{Type: LBRACE, Value: "{"},
	})
}

func TestOrErrorHandling(t *testing.T) {
	assertTokens(t, `val = exec("cmd") or "fallback"`, []Token{
		{Type: IDENT, Value: "val"},
		{Type: ASSIGN, Value: "="},
		{Type: IDENT, Value: "exec"},
		{Type: LPAREN, Value: "("},
		{Type: STRING, Value: "cmd"},
		{Type: RPAREN, Value: ")"},
		{Type: OR, Value: "or"},
		{Type: STRING, Value: "fallback"},
	})
}

func TestNegation(t *testing.T) {
	assertTokens(t, `if !success {`, []Token{
		{Type: IF, Value: "if"},
		{Type: BANG, Value: "!"},
		{Type: IDENT, Value: "success"},
		{Type: LBRACE, Value: "{"},
	})
}

func TestCommaAndArrow(t *testing.T) {
	assertTokens(t, `fn add(a: int, b: int) -> int {`, []Token{
		{Type: FN, Value: "fn"},
		{Type: IDENT, Value: "add"},
		{Type: LPAREN, Value: "("},
		{Type: IDENT, Value: "a"},
		{Type: COLON, Value: ":"},
		{Type: IDENT, Value: "int"},
		{Type: COMMA, Value: ","},
		{Type: IDENT, Value: "b"},
		{Type: COLON, Value: ":"},
		{Type: IDENT, Value: "int"},
		{Type: RPAREN, Value: ")"},
		{Type: ARROW, Value: "->"},
		{Type: IDENT, Value: "int"},
		{Type: LBRACE, Value: "{"},
	})
}

func TestMatchStatement(t *testing.T) {
	assertTokens(t, `match status { "ok" => print("good") _ => exit(1) }`, []Token{
		{Type: MATCH, Value: "match"},
		{Type: IDENT, Value: "status"},
		{Type: LBRACE, Value: "{"},
		{Type: STRING, Value: "ok"},
		{Type: FATARROW, Value: "=>"},
		{Type: IDENT, Value: "print"},
		{Type: LPAREN, Value: "("},
		{Type: STRING, Value: "good"},
		{Type: RPAREN, Value: ")"},
		{Type: UNDERSCORE, Value: "_"},
		{Type: FATARROW, Value: "=>"},
		{Type: IDENT, Value: "exit"},
		{Type: LPAREN, Value: "("},
		{Type: INT, Value: "1"},
		{Type: RPAREN, Value: ")"},
		{Type: RBRACE, Value: "}"},
	})
}

func TestElseKeyword(t *testing.T) {
	assertTokens(t, `} else {`, []Token{
		{Type: RBRACE, Value: "}"},
		{Type: ELSE, Value: "else"},
		{Type: LBRACE, Value: "{"},
	})
}

func TestContinueBreak(t *testing.T) {
	assertTokens(t, `continue`, []Token{
		{Type: CONTINUE, Value: "continue"},
	})
}

func TestAllComparisonOperators(t *testing.T) {
	assertTokens(t, `a == b`, []Token{
		{Type: IDENT, Value: "a"},
		{Type: EQ, Value: "=="},
		{Type: IDENT, Value: "b"},
	})
	assertTokens(t, `a != b`, []Token{
		{Type: IDENT, Value: "a"},
		{Type: NEQ, Value: "!="},
		{Type: IDENT, Value: "b"},
	})
	assertTokens(t, `a < b`, []Token{
		{Type: IDENT, Value: "a"},
		{Type: LT, Value: "<"},
		{Type: IDENT, Value: "b"},
	})
	assertTokens(t, `a >= b`, []Token{
		{Type: IDENT, Value: "a"},
		{Type: GTE, Value: ">="},
		{Type: IDENT, Value: "b"},
	})
	assertTokens(t, `a <= b`, []Token{
		{Type: IDENT, Value: "a"},
		{Type: LTE, Value: "<="},
		{Type: IDENT, Value: "b"},
	})
}

func TestArithmeticOperators(t *testing.T) {
	assertTokens(t, `a + b - c * d / e % f`, []Token{
		{Type: IDENT, Value: "a"},
		{Type: PLUS, Value: "+"},
		{Type: IDENT, Value: "b"},
		{Type: MINUS, Value: "-"},
		{Type: IDENT, Value: "c"},
		{Type: STAR, Value: "*"},
		{Type: IDENT, Value: "d"},
		{Type: SLASH, Value: "/"},
		{Type: IDENT, Value: "e"},
		{Type: PERCENT, Value: "%"},
		{Type: IDENT, Value: "f"},
	})
}

func TestModuloOperator(t *testing.T) {
	assertTokens(t, `x = a % b`, []Token{
		{Type: IDENT, Value: "x"},
		{Type: ASSIGN, Value: "="},
		{Type: IDENT, Value: "a"},
		{Type: PERCENT, Value: "%"},
		{Type: IDENT, Value: "b"},
	})
}

func TestParenthesizedExpression(t *testing.T) {
	assertTokens(t, `x = (a + b) * c`, []Token{
		{Type: IDENT, Value: "x"},
		{Type: ASSIGN, Value: "="},
		{Type: LPAREN, Value: "("},
		{Type: IDENT, Value: "a"},
		{Type: PLUS, Value: "+"},
		{Type: IDENT, Value: "b"},
		{Type: RPAREN, Value: ")"},
		{Type: STAR, Value: "*"},
		{Type: IDENT, Value: "c"},
	})
}

func TestWhileAndBreak(t *testing.T) {
	assertTokens(t, `while x > 0 { break }`, []Token{
		{Type: WHILE, Value: "while"},
		{Type: IDENT, Value: "x"},
		{Type: GT, Value: ">"},
		{Type: INT, Value: "0"},
		{Type: LBRACE, Value: "{"},
		{Type: BREAK, Value: "break"},
		{Type: RBRACE, Value: "}"},
	})
}

func TestAndOrKeywords(t *testing.T) {
	assertTokens(t, `if a and b or c {`, []Token{
		{Type: IF, Value: "if"},
		{Type: IDENT, Value: "a"},
		{Type: AND, Value: "and"},
		{Type: IDENT, Value: "b"},
		{Type: OR, Value: "or"},
		{Type: IDENT, Value: "c"},
		{Type: LBRACE, Value: "{"},
	})
}

func TestBooleans(t *testing.T) {
	assertTokens(t, `x = false`, []Token{
		{Type: IDENT, Value: "x"},
		{Type: ASSIGN, Value: "="},
		{Type: FALSE, Value: "false"},
	})
}

func TestDotAccess(t *testing.T) {
	assertTokens(t, `f.name`, []Token{
		{Type: IDENT, Value: "f"},
		{Type: DOT, Value: "."},
		{Type: IDENT, Value: "name"},
	})
}

func TestSingleLineComment(t *testing.T) {
	assertTokens(t, "x = 1 // this is a comment\ny = 2", []Token{
		{Type: IDENT, Value: "x"},
		{Type: ASSIGN, Value: "="},
		{Type: INT, Value: "1"},
		{Type: IDENT, Value: "y"},
		{Type: ASSIGN, Value: "="},
		{Type: INT, Value: "2"},
	})
}

func TestCommentOnlyLine(t *testing.T) {
	assertTokens(t, "// just a comment", []Token{})
}

func TestStringEscapeQuote(t *testing.T) {
	assertTokens(t, `x = "say \"hello\""`, []Token{
		{Type: IDENT, Value: "x"},
		{Type: ASSIGN, Value: "="},
		{Type: STRING, Value: `say "hello"`},
	})
}

func TestStringEscapeNewline(t *testing.T) {
	assertTokens(t, `x = "line1\nline2"`, []Token{
		{Type: IDENT, Value: "x"},
		{Type: ASSIGN, Value: "="},
		{Type: STRING, Value: "line1\nline2"},
	})
}

func TestStringEscapeTab(t *testing.T) {
	assertTokens(t, `x = "col1\tcol2"`, []Token{
		{Type: IDENT, Value: "x"},
		{Type: ASSIGN, Value: "="},
		{Type: STRING, Value: "col1\tcol2"},
	})
}

func TestStringEscapeBackslash(t *testing.T) {
	assertTokens(t, `x = "path\\to\\file"`, []Token{
		{Type: IDENT, Value: "x"},
		{Type: ASSIGN, Value: "="},
		{Type: STRING, Value: `path\to\file`},
	})
}

func TestStringNoEscapes(t *testing.T) {
	// Existing strings without escapes still work
	assertTokens(t, `x = "hello world"`, []Token{
		{Type: IDENT, Value: "x"},
		{Type: ASSIGN, Value: "="},
		{Type: STRING, Value: "hello world"},
	})
}

func TestCompoundAssignTokens(t *testing.T) {
	assertTokens(t, `x += 5`, []Token{
		{Type: IDENT, Value: "x"},
		{Type: PLUS_ASSIGN, Value: "+="},
		{Type: INT, Value: "5"},
	})
	assertTokens(t, `x -= 3`, []Token{
		{Type: IDENT, Value: "x"},
		{Type: MINUS_ASSIGN, Value: "-="},
		{Type: INT, Value: "3"},
	})
	assertTokens(t, `x *= 2`, []Token{
		{Type: IDENT, Value: "x"},
		{Type: STAR_ASSIGN, Value: "*="},
		{Type: INT, Value: "2"},
	})
	assertTokens(t, `x /= 4`, []Token{
		{Type: IDENT, Value: "x"},
		{Type: SLASH_ASSIGN, Value: "/="},
		{Type: INT, Value: "4"},
	})
}

func TestPipeOperator(t *testing.T) {
	assertTokens(t, `data |> upper`, []Token{
		{Type: IDENT, Value: "data"},
		{Type: PIPE, Value: "|>"},
		{Type: IDENT, Value: "upper"},
	})
}

func TestPipeChain(t *testing.T) {
	assertTokens(t, `data |> upper |> trim`, []Token{
		{Type: IDENT, Value: "data"},
		{Type: PIPE, Value: "|>"},
		{Type: IDENT, Value: "upper"},
		{Type: PIPE, Value: "|>"},
		{Type: IDENT, Value: "trim"},
	})
}

func TestTokenPositions(t *testing.T) {
	tokens := New("x = 1\ny = 2").Tokenize()

	// x is at line 1, col 1
	assert.Equal(t, 1, tokens[0].Line, "x line")
	assert.Equal(t, 1, tokens[0].Col, "x col")

	// 1 is at line 1, col 5
	assert.Equal(t, 1, tokens[2].Line, "1 line")
	assert.Equal(t, 5, tokens[2].Col, "1 col")

	// y is at line 2, col 1
	assert.Equal(t, 2, tokens[3].Line, "y line")
	assert.Equal(t, 1, tokens[3].Col, "y col")
}

func TestUnterminatedString(t *testing.T) {
	tokens := New(`x = "unterminated`).Tokenize()
	// Should produce: IDENT, ASSIGN, ILLEGAL, EOF
	require.True(t, len(tokens) >= 4)
	assert.Equal(t, IDENT, tokens[0].Type)
	assert.Equal(t, ASSIGN, tokens[1].Type)
	assert.Equal(t, ILLEGAL, tokens[2].Type)
	assert.Equal(t, "unterminated string", tokens[2].Value)
	assert.Equal(t, EOF, tokens[len(tokens)-1].Type)
}

func TestUnknownCharacter(t *testing.T) {
	tokens := New(`x = @`).Tokenize()
	// Should produce: IDENT, ASSIGN, ILLEGAL("@"), EOF
	require.True(t, len(tokens) >= 4)
	assert.Equal(t, IDENT, tokens[0].Type)
	assert.Equal(t, ASSIGN, tokens[1].Type)
	assert.Equal(t, ILLEGAL, tokens[2].Type)
	assert.Equal(t, "@", tokens[2].Value)
}

func TestUnknownCharacterTilde(t *testing.T) {
	tokens := New(`~`).Tokenize()
	require.Len(t, tokens, 2) // ILLEGAL + EOF
	assert.Equal(t, ILLEGAL, tokens[0].Type)
	assert.Equal(t, "~", tokens[0].Value)
}

func TestUnicodeIdentifiers(t *testing.T) {
	assertTokens(t, `名前 = "hello"`, []Token{
		{Type: IDENT, Value: "名前"},
		{Type: ASSIGN, Value: "="},
		{Type: STRING, Value: "hello"},
	})
}

func TestUnicodeIdentifierWithDigits(t *testing.T) {
	assertTokens(t, `café2 = 42`, []Token{
		{Type: IDENT, Value: "café2"},
		{Type: ASSIGN, Value: "="},
		{Type: INT, Value: "42"},
	})
}

func TestEmptyInput(t *testing.T) {
	tokens := New("").Tokenize()
	require.Len(t, tokens, 1)
	assert.Equal(t, EOF, tokens[0].Type)
}

func TestOnlyWhitespace(t *testing.T) {
	tokens := New("   \t\n  ").Tokenize()
	require.Len(t, tokens, 1)
	assert.Equal(t, EOF, tokens[0].Type)
}

func TestShebangSkipped(t *testing.T) {
	assertTokens(t, "#!/usr/bin/env langz\nprint(\"hi\")", []Token{
		{Type: IDENT, Value: "print"},
		{Type: LPAREN, Value: "("},
		{Type: STRING, Value: "hi"},
		{Type: RPAREN, Value: ")"},
	})
}

func TestShebangOnlyFile(t *testing.T) {
	tokens := New("#!/usr/bin/env langz\n").Tokenize()
	require.Len(t, tokens, 1)
	assert.Equal(t, EOF, tokens[0].Type)
}

func TestNoShebang(t *testing.T) {
	// # without ! is not a shebang — should produce ILLEGAL
	tokens := New("# not a shebang").Tokenize()
	assert.Equal(t, ILLEGAL, tokens[0].Type)
}

func TestBashBlockSimple(t *testing.T) {
	tokens := New(`bash { echo "hello" }`).Tokenize()
	require.True(t, len(tokens) >= 3)
	assert.Equal(t, BASH, tokens[0].Type)
	assert.Equal(t, BASH_CONTENT, tokens[1].Type)
	assert.Equal(t, `echo "hello"`, tokens[1].Value)
	assert.Equal(t, EOF, tokens[2].Type)
}

func TestBashBlockMultiline(t *testing.T) {
	input := "bash {\n    set -euo pipefail\n    trap 'cleanup' EXIT\n}"
	tokens := New(input).Tokenize()
	require.True(t, len(tokens) >= 3)
	assert.Equal(t, BASH, tokens[0].Type)
	assert.Equal(t, BASH_CONTENT, tokens[1].Type)
	assert.Contains(t, tokens[1].Value, "set -euo pipefail")
	assert.Contains(t, tokens[1].Value, "trap 'cleanup' EXIT")
}

func TestBashBlockNestedBraces(t *testing.T) {
	input := "bash {\n    if [ -f file ]; then\n        echo \"found\"\n    fi\n}"
	tokens := New(input).Tokenize()
	assert.Equal(t, BASH_CONTENT, tokens[1].Type)
	assert.Contains(t, tokens[1].Value, "if [ -f file ]")
}

func TestImportToken(t *testing.T) {
	assertTokens(t, `import "helpers.lz"`, []Token{
		{Type: IMPORT, Value: "import"},
		{Type: STRING, Value: "helpers.lz"},
	})
}

func TestBashKeyword(t *testing.T) {
	// bash without braces just emits the keyword
	tokens := New(`bash`).Tokenize()
	require.Len(t, tokens, 2)
	assert.Equal(t, BASH, tokens[0].Type)
	assert.Equal(t, EOF, tokens[1].Type)
}
