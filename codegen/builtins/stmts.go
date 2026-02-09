package builtins

import (
	"fmt"
	"strings"

	"github.com/tasnimzotder/langz/ast"
)

var stmtBuiltins = map[string]builtinHandler{
	"print": func(args []ast.Node, genExpr ExprGen, _ RawValueGen) string {
		if len(args) == 0 {
			return "echo"
		}
		parts := make([]string, len(args))
		for i, arg := range args {
			parts[i] = genExpr(arg)
		}
		return fmt.Sprintf("echo %s", strings.Join(parts, " "))
	},
	"write": func(args []ast.Node, genExpr ExprGen, _ RawValueGen) string {
		if len(args) != 2 {
			return "# error: write() requires 2 arguments (path, content)"
		}
		return fmt.Sprintf("echo %s > %s", genExpr(args[1]), genExpr(args[0]))
	},
	"append": func(args []ast.Node, genExpr ExprGen, _ RawValueGen) string {
		if len(args) != 2 {
			return "# error: append() requires 2 arguments (path, content)"
		}
		return fmt.Sprintf("echo %s >> %s", genExpr(args[1]), genExpr(args[0]))
	},
	"rm": func(args []ast.Node, genExpr ExprGen, _ RawValueGen) string {
		if len(args) == 0 {
			return "# error: rm() requires 1 argument"
		}
		return fmt.Sprintf("rm -f %s", genExpr(args[0]))
	},
	"rmdir": func(args []ast.Node, genExpr ExprGen, _ RawValueGen) string {
		if len(args) == 0 {
			return "# error: rmdir() requires 1 argument"
		}
		return fmt.Sprintf("rm -rf %s", genExpr(args[0]))
	},
	"mkdir": func(args []ast.Node, genExpr ExprGen, _ RawValueGen) string {
		if len(args) == 0 {
			return "# error: mkdir() requires 1 argument"
		}
		return fmt.Sprintf("mkdir -p %s", genExpr(args[0]))
	},
	"copy": func(args []ast.Node, genExpr ExprGen, _ RawValueGen) string {
		if len(args) != 2 {
			return "# error: copy() requires 2 arguments (src, dst)"
		}
		return fmt.Sprintf("cp %s %s", genExpr(args[0]), genExpr(args[1]))
	},
	"move": func(args []ast.Node, genExpr ExprGen, _ RawValueGen) string {
		if len(args) != 2 {
			return "# error: move() requires 2 arguments (src, dst)"
		}
		return fmt.Sprintf("mv %s %s", genExpr(args[0]), genExpr(args[1]))
	},
	"chmod": func(args []ast.Node, genExpr ExprGen, genRaw RawValueGen) string {
		if len(args) != 2 {
			return "# error: chmod() requires 2 arguments (path, mode)"
		}
		return fmt.Sprintf("chmod %s %s", genRaw(args[1]), genExpr(args[0]))
	},
	"exit": func(args []ast.Node, _ ExprGen, genRaw RawValueGen) string {
		if len(args) == 0 {
			return "exit 0"
		}
		return fmt.Sprintf("exit %s", genRaw(args[0]))
	},
	"sleep": func(args []ast.Node, _ ExprGen, genRaw RawValueGen) string {
		if len(args) == 0 {
			return "# error: sleep() requires 1 argument (seconds)"
		}
		return fmt.Sprintf("sleep %s", genRaw(args[0]))
	},
	"chown": func(args []ast.Node, genExpr ExprGen, genRaw RawValueGen) string {
		if len(args) != 2 {
			return "# error: chown() requires 2 arguments (path, owner)"
		}
		return fmt.Sprintf("chown %s %s", genRaw(args[1]), genExpr(args[0]))
	},
}
