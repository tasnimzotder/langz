package builtins

import (
	"fmt"
	"strings"

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

var stmtBuiltins = map[string]builtinHandler{
	"print": func(args []ast.Node, genExpr ExprGen, _ RawValueGen) string {
		parts := make([]string, len(args))
		for i, arg := range args {
			parts[i] = genExpr(arg)
		}
		return fmt.Sprintf("echo %s", strings.Join(parts, " "))
	},
	"write": func(args []ast.Node, genExpr ExprGen, _ RawValueGen) string {
		if len(args) == 2 {
			return fmt.Sprintf("echo %s > %s", genExpr(args[1]), genExpr(args[0]))
		}
		return ""
	},
	"rm": func(args []ast.Node, genExpr ExprGen, _ RawValueGen) string {
		return fmt.Sprintf("rm %s", genExpr(args[0]))
	},
	"mkdir": func(args []ast.Node, genExpr ExprGen, _ RawValueGen) string {
		return fmt.Sprintf("mkdir -p %s", genExpr(args[0]))
	},
	"copy": func(args []ast.Node, genExpr ExprGen, _ RawValueGen) string {
		if len(args) == 2 {
			return fmt.Sprintf("cp %s %s", genExpr(args[0]), genExpr(args[1]))
		}
		return ""
	},
	"move": func(args []ast.Node, genExpr ExprGen, _ RawValueGen) string {
		if len(args) == 2 {
			return fmt.Sprintf("mv %s %s", genExpr(args[0]), genExpr(args[1]))
		}
		return ""
	},
	"chmod": func(args []ast.Node, genExpr ExprGen, genRaw RawValueGen) string {
		if len(args) == 2 {
			return fmt.Sprintf("chmod %s %s", genRaw(args[1]), genExpr(args[0]))
		}
		return ""
	},
	"exit": func(args []ast.Node, _ ExprGen, genRaw RawValueGen) string {
		return fmt.Sprintf("exit %s", genRaw(args[0]))
	},
}

var exprBuiltins = map[string]builtinHandler{
	"print": stmtBuiltins["print"],
	"exec": func(args []ast.Node, _ ExprGen, genRaw RawValueGen) string {
		return fmt.Sprintf("$(%s)", genRaw(args[0]))
	},
	"env": func(args []ast.Node, _ ExprGen, genRaw RawValueGen) string {
		return fmt.Sprintf(`"${%s}"`, genRaw(args[0]))
	},
	"read": func(args []ast.Node, genExpr ExprGen, _ RawValueGen) string {
		return fmt.Sprintf("$(cat %s)", genExpr(args[0]))
	},
	"glob": func(args []ast.Node, _ ExprGen, genRaw RawValueGen) string {
		return fmt.Sprintf("(%s)", genRaw(args[0]))
	},
	"exists": func(args []ast.Node, genExpr ExprGen, _ RawValueGen) string {
		return fmt.Sprintf("[ -e %s ]", genExpr(args[0]))
	},
	"fetch": func(args []ast.Node, genExpr ExprGen, _ RawValueGen) string {
		return fmt.Sprintf("$(curl -sf %s)", genExpr(args[0]))
	},
}
