package parser

import (
	"strings"

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
		if isCompoundAssign(p.peek().Type) {
			return p.parseCompoundAssignment()
		}
		if p.peek().Type == lexer.LBRACKET {
			return p.parseIndexOrExpr()
		}
		if p.peek().Type == lexer.LPAREN {
			return p.parseFuncCall()
		}
		return p.parseExpression()
	case lexer.STRING, lexer.INT, lexer.TRUE, lexer.FALSE, lexer.BANG:
		return p.parseExpression()
	case lexer.ILLEGAL:
		p.addError(p.current.Value)
		p.advance()
		return nil
	default:
		p.addError("unexpected token: " + string(p.current.Type))
		p.advance()
		return nil
	}
}

func (p *Parser) parseAssignment() *ast.Assignment {
	name := p.expect(lexer.IDENT)
	p.expect(lexer.ASSIGN)

	value := p.parsePipeExpr()

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

	condition := p.parseCondition()
	body := p.parseBlock()

	var elseBody []ast.Node
	if p.current.Type == lexer.ELSE {
		p.advance()
		if p.current.Type == lexer.IF {
			// elif: recursively parse as IfStmt
			elif := p.parseIf()
			elseBody = []ast.Node{elif}
		} else {
			elseBody = p.parseBlock()
		}
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

	condition := p.parseCondition()
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

		var defaultVal ast.Node
		if p.current.Type == lexer.ASSIGN {
			p.advance()
			defaultVal = p.parsePrimary()
		}

		params = append(params, ast.Param{Name: paramName.Value, Type: paramType.Value, Default: defaultVal})
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

		var body []ast.Node
		if p.current.Type == lexer.LBRACE {
			// Block arm: { stmt; stmt; ... }
			body = p.parseBlock()
		} else {
			// Single-statement arm: collect until next pattern or closing brace
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
		}

		cases = append(cases, ast.MatchCase{Pattern: pattern, Body: body})
	}

	p.expect(lexer.RBRACE)

	return &ast.MatchStmt{Expr: expr, Cases: cases}
}

func (p *Parser) parseIndexOrExpr() ast.Node {
	name := p.expect(lexer.IDENT)
	p.expect(lexer.LBRACKET)
	index := p.parseExpression()
	p.expect(lexer.RBRACKET)

	if p.current.Type == lexer.ASSIGN {
		// Index assignment: arr[0] = value
		p.advance()
		value := p.parsePipeExpr()
		return &ast.IndexAssignment{Object: name.Value, Index: index, Value: value}
	}

	// Index expression used as statement (rare but valid)
	return &ast.IndexExpr{Object: &ast.Identifier{Name: name.Value}, Index: index}
}

func isCompoundAssign(t lexer.TokenType) bool {
	return t == lexer.PLUS_ASSIGN || t == lexer.MINUS_ASSIGN ||
		t == lexer.STAR_ASSIGN || t == lexer.SLASH_ASSIGN
}

func (p *Parser) parseCompoundAssignment() *ast.Assignment {
	name := p.expect(lexer.IDENT)
	op := p.current
	p.advance()
	value := p.parsePipeExpr()
	arithOp := strings.TrimSuffix(op.Value, "=")
	return &ast.Assignment{
		Name: name.Value,
		Value: &ast.BinaryExpr{
			Left:  &ast.Identifier{Name: name.Value},
			Op:    arithOp,
			Right: value,
		},
	}
}

func (p *Parser) parseReturn() *ast.ReturnStmt {
	p.expect(lexer.RETURN)

	var value ast.Node
	if p.current.Type != lexer.RBRACE && p.current.Type != lexer.EOF {
		value = p.parseExpression()
	}

	return &ast.ReturnStmt{Value: value}
}
