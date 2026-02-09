package parser

import (
	"fmt"

	"github.com/tasnimzotder/langz/ast"
	"github.com/tasnimzotder/langz/lexer"
)

// Parser converts a token stream into an AST.
type Parser struct {
	tokens  []lexer.Token
	pos     int
	current lexer.Token
	errors  []string
}

// New creates a new Parser from a slice of tokens.
func New(tokens []lexer.Token) *Parser {
	p := &Parser{tokens: tokens}
	if len(tokens) > 0 {
		p.current = tokens[0]
	}
	return p
}

func (p *Parser) advance() {
	p.pos++
	if p.pos < len(p.tokens) {
		p.current = p.tokens[p.pos]
	}
}

func (p *Parser) peek() lexer.Token {
	return p.peekAt(1)
}

func (p *Parser) peekAt(offset int) lexer.Token {
	idx := p.pos + offset
	if idx < len(p.tokens) {
		return p.tokens[idx]
	}
	return lexer.Token{Type: lexer.EOF}
}

func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, fmt.Sprintf("line %d, col %d: %s", p.current.Line, p.current.Col, msg))
}

func (p *Parser) expect(t lexer.TokenType) lexer.Token {
	if p.current.Type != t {
		p.addError(fmt.Sprintf("expected %s, got %s", t, p.current.Type))
		return lexer.Token{Type: t}
	}
	tok := p.current
	p.advance()
	return tok
}

// Parse parses tokens into a program. Panics on errors (legacy).
func (p *Parser) Parse() *ast.Program {
	prog, err := p.ParseWithErrors()
	if err != nil {
		panic(err.Error())
	}
	return prog
}

// ParseWithErrors parses tokens and returns errors instead of panicking.
func (p *Parser) ParseWithErrors() (*ast.Program, error) {
	prog := &ast.Program{}

	for p.current.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			prog.Statements = append(prog.Statements, stmt)
		}
	}

	if len(p.errors) > 0 {
		return prog, fmt.Errorf("%s", p.errors[0])
	}
	return prog, nil
}
