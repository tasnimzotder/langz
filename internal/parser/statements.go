package parser

import (
	"github.com/tasnimzotder/langz/internal/ast"
	"github.com/tasnimzotder/langz/internal/lexer"
)

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
	case lexer.BREAK:
		p.advance()
		return &ast.BreakStmt{}
	case lexer.WHILE:
		return p.parseWhile()
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

func (p *Parser) parseWhile() *ast.WhileStmt {
	p.expect(lexer.WHILE)

	condition := p.parseExpression()
	body := p.parseBlock()

	return &ast.WhileStmt{Condition: condition, Body: body}
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
