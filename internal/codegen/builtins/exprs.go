package builtins

import (
	"fmt"

	"github.com/tasnimzotder/langz/internal/ast"
)

var exprBuiltins = map[string]builtinHandler{
	"print": stmtBuiltins["print"],
	"exec": func(args []ast.Node, _ []ast.KeywordArg, _ ExprGen, genRaw RawValueGen) string {
		if len(args) == 0 {
			return "# error: exec() requires 1 argument"
		}
		return fmt.Sprintf("$(%s)", genRaw(args[0]))
	},
	"env": func(args []ast.Node, _ []ast.KeywordArg, _ ExprGen, genRaw RawValueGen) string {
		if len(args) == 0 {
			return "# error: env() requires 1 argument"
		}
		return fmt.Sprintf(`"${%s}"`, genRaw(args[0]))
	},
	"read": func(args []ast.Node, _ []ast.KeywordArg, genExpr ExprGen, _ RawValueGen) string {
		if len(args) == 0 {
			return "# error: read() requires 1 argument"
		}
		return fmt.Sprintf("$(cat %s)", genExpr(args[0]))
	},
	"glob": func(args []ast.Node, _ []ast.KeywordArg, _ ExprGen, genRaw RawValueGen) string {
		if len(args) == 0 {
			return "# error: glob() requires 1 argument"
		}
		return fmt.Sprintf("(%s)", genRaw(args[0]))
	},
	"exists": func(args []ast.Node, _ []ast.KeywordArg, genExpr ExprGen, _ RawValueGen) string {
		if len(args) == 0 {
			return "# error: exists() requires 1 argument"
		}
		return fmt.Sprintf("[ -e %s ]", genExpr(args[0]))
	},
	"is_file": func(args []ast.Node, _ []ast.KeywordArg, genExpr ExprGen, _ RawValueGen) string {
		if len(args) == 0 {
			return "# error: is_file() requires 1 argument"
		}
		return fmt.Sprintf("[ -f %s ]", genExpr(args[0]))
	},
	"is_dir": func(args []ast.Node, _ []ast.KeywordArg, genExpr ExprGen, _ RawValueGen) string {
		if len(args) == 0 {
			return "# error: is_dir() requires 1 argument"
		}
		return fmt.Sprintf("[ -d %s ]", genExpr(args[0]))
	},
	"range": func(args []ast.Node, _ []ast.KeywordArg, _ ExprGen, genRaw RawValueGen) string {
		if len(args) == 2 {
			return fmt.Sprintf("$(seq %s %s)", genRaw(args[0]), genRaw(args[1]))
		}
		if len(args) == 1 {
			return fmt.Sprintf("$(seq 0 %s)", genRaw(args[0]))
		}
		return "# error: range() requires 1 or 2 arguments"
	},
	"args": func(_ []ast.Node, _ []ast.KeywordArg, _ ExprGen, _ RawValueGen) string {
		return `("$@")`
	},
	"os": func(_ []ast.Node, _ []ast.KeywordArg, _ ExprGen, _ RawValueGen) string {
		return "$(uname -s | tr '[:upper:]' '[:lower:]')"
	},
	"arch": func(_ []ast.Node, _ []ast.KeywordArg, _ ExprGen, _ RawValueGen) string {
		return "$(uname -m)"
	},
	"hostname": func(_ []ast.Node, _ []ast.KeywordArg, _ ExprGen, _ RawValueGen) string {
		return "$(hostname)"
	},
	"whoami": func(_ []ast.Node, _ []ast.KeywordArg, _ ExprGen, _ RawValueGen) string {
		return "$(whoami)"
	},
	"dirname": func(args []ast.Node, _ []ast.KeywordArg, genExpr ExprGen, _ RawValueGen) string {
		if len(args) == 0 {
			return "# error: dirname() requires 1 argument"
		}
		return fmt.Sprintf("$(dirname %s)", genExpr(args[0]))
	},
	"basename": func(args []ast.Node, _ []ast.KeywordArg, genExpr ExprGen, _ RawValueGen) string {
		if len(args) == 0 {
			return "# error: basename() requires 1 argument"
		}
		return fmt.Sprintf("$(basename %s)", genExpr(args[0]))
	},
	"len": func(args []ast.Node, _ []ast.KeywordArg, _ ExprGen, genRaw RawValueGen) string {
		if len(args) == 0 {
			return "# error: len() requires 1 argument"
		}
		return fmt.Sprintf("${#%s[@]}", genRaw(args[0]))
	},
	"trim": func(args []ast.Node, _ []ast.KeywordArg, genExpr ExprGen, _ RawValueGen) string {
		if len(args) == 0 {
			return "# error: trim() requires 1 argument"
		}
		return fmt.Sprintf("$(echo %s | xargs)", genExpr(args[0]))
	},
	"upper": func(args []ast.Node, _ []ast.KeywordArg, genExpr ExprGen, _ RawValueGen) string {
		if len(args) == 0 {
			return "# error: upper() requires 1 argument"
		}
		return fmt.Sprintf("$(echo %s | tr '[:lower:]' '[:upper:]')", genExpr(args[0]))
	},
	"lower": func(args []ast.Node, _ []ast.KeywordArg, genExpr ExprGen, _ RawValueGen) string {
		if len(args) == 0 {
			return "# error: lower() requires 1 argument"
		}
		return fmt.Sprintf("$(echo %s | tr '[:upper:]' '[:lower:]')", genExpr(args[0]))
	},
	"json_get": func(args []ast.Node, _ []ast.KeywordArg, genExpr ExprGen, _ RawValueGen) string {
		if len(args) < 2 {
			return "# error: json_get() requires 2 arguments (data, path)"
		}
		return fmt.Sprintf("$(echo %s | jq -r %s)", genExpr(args[0]), genExpr(args[1]))
	},
	"timestamp": func(_ []ast.Node, _ []ast.KeywordArg, _ ExprGen, _ RawValueGen) string {
		return "$(date +%s)"
	},
	"date": func(_ []ast.Node, _ []ast.KeywordArg, _ ExprGen, _ RawValueGen) string {
		return `$(date +"%Y-%m-%d")`
	},
}
