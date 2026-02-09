package builtins

import (
	"github.com/tasnimzotder/langz/ast"
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
func GenStmt(name string, args []ast.Node, genExpr ExprGen, genRaw RawValueGen) StmtResult {
	handler, ok := stmtBuiltins[name]
	if !ok {
		return StmtResult{}
	}
	return StmtResult{Code: handler(args, genExpr, genRaw), OK: true}
}

// GenExpr generates a Bash expression for a builtin function call.
// Returns (code, true) if the function is a known builtin, ("", false) otherwise.
func GenExpr(name string, args []ast.Node, genExpr ExprGen, genRaw RawValueGen) ExprResult {
	handler, ok := exprBuiltins[name]
	if !ok {
		return ExprResult{}
	}
	return ExprResult{Code: handler(args, genExpr, genRaw), OK: true}
}

type builtinHandler func(args []ast.Node, genExpr ExprGen, genRaw RawValueGen) string
