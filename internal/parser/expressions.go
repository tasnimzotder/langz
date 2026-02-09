package parser

import (
	"github.com/tasnimzotder/langz/internal/ast"
	"github.com/tasnimzotder/langz/internal/lexer"
)

// parseCondition handles `or` at the lowest precedence level, used only
// in condition contexts (if/while). This keeps `or` out of parseExpression
// so that assignment fallback (`x = expr or fallback`) still works.
func (p *Parser) parseCondition() ast.Node {
	left := p.parseExpression()
	for p.current.Type == lexer.OR {
		op := p.current.Value
		p.advance()
		right := p.parseExpression()
		left = &ast.BinaryExpr{Left: left, Op: op, Right: right}
	}
	return left
}

func (p *Parser) parseExpression() ast.Node {
	left := p.parseComparison()

	// Handle logical `and` operator
	for p.current.Type == lexer.AND {
		op := p.current.Value
		p.advance()
		right := p.parseComparison()
		left = &ast.BinaryExpr{Left: left, Op: op, Right: right}
	}

	return left
}

func (p *Parser) parseComparison() ast.Node {
	left := p.parseAdditive()

	if isComparisonOp(p.current.Type) {
		op := p.current.Value
		p.advance()
		right := p.parseAdditive()
		left = &ast.BinaryExpr{Left: left, Op: op, Right: right}
	}

	return left
}

func (p *Parser) parseAdditive() ast.Node {
	left := p.parseMultiplicative()

	for p.current.Type == lexer.PLUS || p.current.Type == lexer.MINUS {
		op := p.current.Value
		p.advance()
		right := p.parseMultiplicative()
		left = &ast.BinaryExpr{Left: left, Op: op, Right: right}
	}

	return left
}

func (p *Parser) parseMultiplicative() ast.Node {
	left := p.parseUnary()

	for p.current.Type == lexer.STAR || p.current.Type == lexer.SLASH || p.current.Type == lexer.PERCENT {
		op := p.current.Value
		p.advance()
		right := p.parseUnary()
		left = &ast.BinaryExpr{Left: left, Op: op, Right: right}
	}

	return left
}

func (p *Parser) parseUnary() ast.Node {
	left := p.parsePrimary()

	// Handle dot access: expr.field
	for p.current.Type == lexer.DOT {
		p.advance()
		field := p.expect(lexer.IDENT)
		left = &ast.DotExpr{Object: left, Field: field.Value}
	}

	return left
}

func isComparisonOp(t lexer.TokenType) bool {
	return t == lexer.GT || t == lexer.LT || t == lexer.GTE ||
		t == lexer.LTE || t == lexer.EQ || t == lexer.NEQ
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

	case lexer.LPAREN:
		// Check if this is a function call (handled elsewhere) or grouped expression
		// Grouped expression: (expr)
		p.advance() // skip (
		expr := p.parseExpression()
		p.expect(lexer.RPAREN)
		return expr

	case lexer.LBRACKET:
		return p.parseListLiteral()

	case lexer.LBRACE:
		// Distinguish map literal from block: { ident : ... } is a map
		if p.peek().Type == lexer.IDENT && p.peekAt(2).Type == lexer.COLON {
			return p.parseMapLiteral()
		}
		return nil

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

func (p *Parser) parseListLiteral() *ast.ListLiteral {
	p.expect(lexer.LBRACKET)

	var elements []ast.Node
	for p.current.Type != lexer.RBRACKET && p.current.Type != lexer.EOF {
		elem := p.parseExpression()
		if elem != nil {
			elements = append(elements, elem)
		}
		if p.current.Type == lexer.COMMA {
			p.advance()
		}
	}

	p.expect(lexer.RBRACKET)
	return &ast.ListLiteral{Elements: elements}
}

func (p *Parser) parseMapLiteral() *ast.MapLiteral {
	p.expect(lexer.LBRACE)

	var keys []string
	var values []ast.Node

	for p.current.Type != lexer.RBRACE && p.current.Type != lexer.EOF {
		key := p.expect(lexer.IDENT)
		p.expect(lexer.COLON)
		value := p.parseExpression()

		keys = append(keys, key.Value)
		values = append(values, value)

		if p.current.Type == lexer.COMMA {
			p.advance()
		}
	}

	p.expect(lexer.RBRACE)
	return &ast.MapLiteral{Keys: keys, Values: values}
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
