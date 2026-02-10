package parser

import (
	"fmt"

	"github.com/tasnimzotder/langz/internal/ast"
	"github.com/tasnimzotder/langz/internal/lexer"
)

// ParseError represents a structured parse error with position.
type ParseError struct {
	Line    int
	Col     int
	Message string
}

// Parser converts a token stream into an AST.
type Parser struct {
	tokens  []lexer.Token
	pos     int
	current lexer.Token
	errors  []ParseError
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
	p.errors = append(p.errors, ParseError{
		Line:    p.current.Line,
		Col:     p.current.Col,
		Message: msg,
	})
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

// ParseWithErrors parses tokens and returns the first error.
func (p *Parser) ParseWithErrors() (prog *ast.Program, err error) {
	defer func() {
		if r := recover(); r != nil {
			if prog == nil {
				prog = &ast.Program{}
			}
			err = fmt.Errorf("internal error: %v", r)
		}
	}()
	prog = &ast.Program{}

	for p.current.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			prog.Statements = append(prog.Statements, stmt)
		}
	}

	if len(p.errors) > 0 {
		e := p.errors[0]
		return prog, fmt.Errorf("line %d, col %d: %s", e.Line, e.Col, e.Message)
	}
	return prog, nil
}

// ParseAllErrors parses tokens and returns ALL structured errors.
// The returned *ast.Program is always non-nil (partial program on errors).
func (p *Parser) ParseAllErrors() (prog *ast.Program, errs []ParseError) {
	defer func() {
		if r := recover(); r != nil {
			if prog == nil {
				prog = &ast.Program{}
			}
			errs = append(errs, ParseError{
				Line:    p.current.Line,
				Col:     p.current.Col,
				Message: fmt.Sprintf("internal error: %v", r),
			})
		}
	}()
	prog = &ast.Program{}

	for p.current.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			prog.Statements = append(prog.Statements, stmt)
		}
	}

	return prog, p.errors
}
