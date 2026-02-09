package parser

import (
	"github.com/tasnimzotder/langz/internal/ast"
	"github.com/tasnimzotder/langz/internal/lexer"
)

func parse(input string) *ast.Program {
	tokens := lexer.New(input).Tokenize()
	return New(tokens).Parse()
}
