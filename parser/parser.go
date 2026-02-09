package parser

import (
	"fmt"

	"github.com/tasnimzotder/langz/ast"
	"github.com/tasnimzotder/langz/lexer"
)

type Parser struct {
	tokens  []lexer.Token
	pos     int
	current lexer.Token
	errors  []string
}

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
	if p.pos+1 < len(p.tokens) {
		return p.tokens[p.pos+1]
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

func (p *Parser) parseStatement() ast.Node {
	switch p.current.Type {
	case lexer.IF:
		return p.parseIf()
	case lexer.FOR:
		return p.parseFor()
	case lexer.MATCH:
		return p.parseMatch()
	case lexer.FN:
		return p.parseFuncDecl()
	case lexer.RETURN:
		return p.parseReturn()
	case lexer.CONTINUE:
		p.advance()
		return &ast.ContinueStmt{}
	case lexer.IDENT:
		if p.peek().Type == lexer.ASSIGN {
			return p.parseAssignment()
		}
		if p.peek().Type == lexer.LPAREN {
			return p.parseFuncCall()
		}
		return p.parseExpression()
	case lexer.STRING, lexer.INT, lexer.TRUE, lexer.FALSE, lexer.BANG:
		return p.parseExpression()
	default:
		p.advance()
		return nil
	}
}

func (p *Parser) parseAssignment() *ast.Assignment {
	name := p.expect(lexer.IDENT)
	p.expect(lexer.ASSIGN)

	value := p.parseExpression()

	if p.current.Type == lexer.OR {
		p.advance()
		fallback := p.parseOrFallback()
		value = &ast.OrExpr{Expr: value, Fallback: fallback}
	}

	return &ast.Assignment{Name: name.Value, Value: value}
}

func (p *Parser) parseExpression() ast.Node {
	left := p.parsePrimary()

	// Handle dot access: expr.field
	for p.current.Type == lexer.DOT {
		p.advance()
		field := p.expect(lexer.IDENT)
		left = &ast.DotExpr{Object: left, Field: field.Value}
	}

	// Handle binary operators: expr > expr
	if p.current.Type == lexer.GT {
		op := p.current.Value
		p.advance()
		right := p.parsePrimary()
		left = &ast.BinaryExpr{Left: left, Op: op, Right: right}
	}

	return left
}

func (p *Parser) parsePrimary() ast.Node {
	switch p.current.Type {
	case lexer.STRING:
		tok := p.current
		p.advance()
		return &ast.StringLiteral{Value: tok.Value}

	case lexer.INT:
		tok := p.current
		p.advance()
		return &ast.IntLiteral{Value: tok.Value}

	case lexer.TRUE:
		p.advance()
		return &ast.BoolLiteral{Value: true}

	case lexer.FALSE:
		p.advance()
		return &ast.BoolLiteral{Value: false}

	case lexer.BANG:
		p.advance()
		operand := p.parsePrimary()
		return &ast.UnaryExpr{Op: "!", Operand: operand}

	case lexer.IDENT:
		if p.peek().Type == lexer.LPAREN {
			return p.parseFuncCall()
		}
		tok := p.current
		p.advance()
		return &ast.Identifier{Name: tok.Value}

	default:
		p.advance()
		return nil
	}
}

func (p *Parser) parseFuncCall() *ast.FuncCall {
	name := p.expect(lexer.IDENT)
	p.expect(lexer.LPAREN)

	var args []ast.Node
	for p.current.Type != lexer.RPAREN && p.current.Type != lexer.EOF {
		arg := p.parseExpression()
		if arg != nil {
			args = append(args, arg)
		}
		if p.current.Type == lexer.COMMA {
			p.advance()
		}
	}

	p.expect(lexer.RPAREN)
	return &ast.FuncCall{Name: name.Value, Args: args}
}

func (p *Parser) parseBlock() []ast.Node {
	p.expect(lexer.LBRACE)

	var stmts []ast.Node
	for p.current.Type != lexer.RBRACE && p.current.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			stmts = append(stmts, stmt)
		}
	}

	p.expect(lexer.RBRACE)
	return stmts
}

func (p *Parser) parseIf() *ast.IfStmt {
	p.expect(lexer.IF)

	condition := p.parseExpression()
	body := p.parseBlock()

	var elseBody []ast.Node
	if p.current.Type == lexer.ELSE {
		p.advance()
		elseBody = p.parseBlock()
	}

	return &ast.IfStmt{Condition: condition, Body: body, ElseBody: elseBody}
}

func (p *Parser) parseFor() *ast.ForStmt {
	p.expect(lexer.FOR)

	varName := p.expect(lexer.IDENT)
	p.expect(lexer.IN)

	collection := p.parseExpression()
	body := p.parseBlock()

	return &ast.ForStmt{Var: varName.Value, Collection: collection, Body: body}
}

func (p *Parser) parseFuncDecl() *ast.FuncDecl {
	p.expect(lexer.FN)

	name := p.expect(lexer.IDENT)
	p.expect(lexer.LPAREN)

	var params []ast.Param
	for p.current.Type != lexer.RPAREN && p.current.Type != lexer.EOF {
		paramName := p.expect(lexer.IDENT)
		p.expect(lexer.COLON)
		paramType := p.expect(lexer.IDENT)
		params = append(params, ast.Param{Name: paramName.Value, Type: paramType.Value})
		if p.current.Type == lexer.COMMA {
			p.advance()
		}
	}

	p.expect(lexer.RPAREN)

	var returnType string
	if p.current.Type == lexer.ARROW {
		p.advance()
		returnType = p.expect(lexer.IDENT).Value
	}

	body := p.parseBlock()

	return &ast.FuncDecl{
		Name:       name.Value,
		Params:     params,
		ReturnType: returnType,
		Body:       body,
	}
}

func (p *Parser) parseOrFallback() ast.Node {
	// or { block }
	if p.current.Type == lexer.LBRACE {
		stmts := p.parseBlock()
		return &ast.BlockExpr{Statements: stmts}
	}

	// or continue
	if p.current.Type == lexer.CONTINUE {
		p.advance()
		return &ast.ContinueStmt{}
	}

	// or return expr
	if p.current.Type == lexer.RETURN {
		return p.parseReturn()
	}

	// or expr (value or func call like exit(1))
	return p.parseExpression()
}

func (p *Parser) parseMatch() *ast.MatchStmt {
	p.expect(lexer.MATCH)

	expr := p.parseExpression()
	p.expect(lexer.LBRACE)

	var cases []ast.MatchCase
	for p.current.Type != lexer.RBRACE && p.current.Type != lexer.EOF {
		var pattern ast.Node

		if p.current.Type == lexer.UNDERSCORE {
			// Wildcard: _ => ...
			p.advance()
			pattern = nil
		} else {
			pattern = p.parseExpression()
		}

		p.expect(lexer.FATARROW)

		// Collect body statements until the next pattern, wildcard, or closing brace
		var body []ast.Node
		for p.current.Type != lexer.RBRACE &&
			p.current.Type != lexer.UNDERSCORE &&
			p.current.Type != lexer.FATARROW &&
			p.current.Type != lexer.EOF {

			// Peek ahead to see if this is the start of a new case
			if (p.current.Type == lexer.STRING || p.current.Type == lexer.INT ||
				p.current.Type == lexer.TRUE || p.current.Type == lexer.FALSE) &&
				p.peek().Type == lexer.FATARROW {
				break
			}

			stmt := p.parseStatement()
			if stmt != nil {
				body = append(body, stmt)
			}
		}

		cases = append(cases, ast.MatchCase{Pattern: pattern, Body: body})
	}

	p.expect(lexer.RBRACE)

	return &ast.MatchStmt{Expr: expr, Cases: cases}
}

func (p *Parser) parseReturn() *ast.ReturnStmt {
	p.expect(lexer.RETURN)

	var value ast.Node
	if p.current.Type != lexer.RBRACE && p.current.Type != lexer.EOF {
		value = p.parseExpression()
	}

	return &ast.ReturnStmt{Value: value}
}
