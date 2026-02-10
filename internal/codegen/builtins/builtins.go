package builtins

import (
	"github.com/tasnimzotder/langz/internal/ast"
)

// ExprGen is a function that converts an AST node to a Bash expression string.
type ExprGen func(node ast.Node) string

// RawValueGen extracts the raw string value from a node (no quoting).
type RawValueGen func(node ast.Node) string

// StmtResult holds the generated Bash for a statement-level builtin.
type StmtResult struct {
	Code string
	OK   bool
}

// ExprResult holds the generated Bash for an expression-level builtin.
type ExprResult struct {
	Code string
	OK   bool
}

// GenStmt generates a Bash statement for a builtin function call.
// Returns (code, true) if the function is a known builtin, ("", false) otherwise.
func GenStmt(name string, args []ast.Node, kwargs []ast.KeywordArg, genExpr ExprGen, genRaw RawValueGen) StmtResult {
	handler, ok := stmtBuiltins[name]
	if !ok {
		return StmtResult{}
	}
	return StmtResult{Code: handler(args, kwargs, genExpr, genRaw), OK: true}
}

// GenExpr generates a Bash expression for a builtin function call.
// Returns (code, true) if the function is a known builtin, ("", false) otherwise.
func GenExpr(name string, args []ast.Node, kwargs []ast.KeywordArg, genExpr ExprGen, genRaw RawValueGen) ExprResult {
	handler, ok := exprBuiltins[name]
	if !ok {
		return ExprResult{}
	}
	return ExprResult{Code: handler(args, kwargs, genExpr, genRaw), OK: true}
}

// FindKwarg looks up a keyword argument by key name.
func FindKwarg(kwargs []ast.KeywordArg, key string) (ast.Node, bool) {
	for _, kw := range kwargs {
		if kw.Key == key {
			return kw.Value, true
		}
	}
	return nil, false
}

type builtinHandler func(args []ast.Node, kwargs []ast.KeywordArg, genExpr ExprGen, genRaw RawValueGen) string
