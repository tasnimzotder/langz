package parser

import (
	"github.com/tasnimzotder/langz/internal/ast"
	"github.com/tasnimzotder/langz/internal/lexer"
)

func parse(input string) *ast.Program {
	tokens := lexer.New(input).Tokenize()
	prog, err := New(tokens).ParseWithErrors()
	if err != nil {
		panic(err.Error())
	}
	return prog
}
