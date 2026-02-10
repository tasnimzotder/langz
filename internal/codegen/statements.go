package codegen

import (
	"fmt"
	"strings"

	"github.com/tasnimzotder/langz/internal/ast"
	"github.com/tasnimzotder/langz/internal/codegen/builtins"
)

func (g *Generator) genStatement(node ast.Node) {
	switch n := node.(type) {
	case *ast.Assignment:
		g.genAssignment(n)
	case *ast.FuncCall:
		g.genFuncCallStmt(n)
	case *ast.FuncDecl:
		g.genFuncDecl(n)
	case *ast.IfStmt:
		g.genIf(n)
	case *ast.ForStmt:
		g.genFor(n)
	case *ast.MatchStmt:
		g.genMatch(n)
	case *ast.ReturnStmt:
		g.genReturn(n)
	case *ast.ContinueStmt:
		g.writeln("continue")
	case *ast.BreakStmt:
		g.writeln("break")
	case *ast.IndexAssignment:
		g.genIndexAssignment(n)
	case *ast.WhileStmt:
		g.genWhile(n)
	}
}

func (g *Generator) genAssignment(a *ast.Assignment) {
	if orExpr, ok := a.Value.(*ast.OrExpr); ok {
		g.genOrAssignment(a.Name, orExpr)
		return
	}
	if mapLit, ok := a.Value.(*ast.MapLiteral); ok {
		g.genMapAssignment(a.Name, mapLit)
		return
	}
	if call, ok := a.Value.(*ast.FuncCall); ok && call.Name == "fetch" {
		g.genFetchAssignment(a.Name, call)
		return
	}
	g.writeIndent()
	value := g.genExpr(a.Value)
	g.write(fmt.Sprintf("%s=%s\n", a.Name, value))
}

func (g *Generator) genMapAssignment(name string, m *ast.MapLiteral) {
	for i, key := range m.Keys {
		g.writeln(fmt.Sprintf("%s_%s=%s", name, key, g.genExpr(m.Values[i])))
	}
}

func (g *Generator) genOrAssignment(name string, or *ast.OrExpr) {
	// Special case: env("VAR") or "default" -> var="${VAR:-default}"
	if call, ok := or.Expr.(*ast.FuncCall); ok && call.Name == "env" {
		if strVal, ok := or.Fallback.(*ast.StringLiteral); ok {
			envName := g.genRawValue(call.Args[0])
			g.writeln(fmt.Sprintf(`%s="${%s:-%s}"`, name, envName, interpolate(strVal.Value)))
			return
		}
	}

	// Special case: fetch(url) or fallback -> curl + status check + fallback
	if call, ok := or.Expr.(*ast.FuncCall); ok && call.Name == "fetch" {
		g.genFetchAssignment(name, call)
		g.writeln(`if [ "$_status" -ge 200 ] && [ "$_status" -lt 300 ]; then`)
		g.indent++
		g.writeln("true")
		g.indent--
		g.writeln("else")
		g.genOrFallback(name, or.Fallback)
		g.writeln("fi")
		return
	}

	// General case: if name=$(expr 2>/dev/null); then true; else fallback; fi
	expr := g.genExpr(or.Expr)
	g.writeln(fmt.Sprintf("if %s=$(%s 2>/dev/null); then", name, stripSubshell(expr)))
	g.indent++
	g.writeln("true")
	g.indent--
	g.writeln("else")
	g.genOrFallback(name, or.Fallback)
	g.writeln("fi")
}

func stripSubshell(s string) string {
	// Remove $() wrapper if present
	if len(s) >= 3 && s[0] == '$' && s[1] == '(' && s[len(s)-1] == ')' {
		return s[2 : len(s)-1]
	}
	return s
}

func (g *Generator) genOrFallback(name string, fallback ast.Node) {
	switch n := fallback.(type) {
	case *ast.BlockExpr:
		g.genBlock(n.Statements)
	case *ast.ContinueStmt:
		g.indent++
		g.writeln("continue")
		g.indent--
	case *ast.FuncCall:
		g.indent++
		g.genFuncCallStmt(n)
		g.indent--
	case *ast.ReturnStmt:
		g.indent++
		g.genReturn(n)
		g.indent--
	default:
		g.indent++
		g.writeIndent()
		g.write(fmt.Sprintf("%s=%s\n", name, g.genExpr(fallback)))
		g.indent--
	}
}

func (g *Generator) genFuncCallStmt(f *ast.FuncCall) {
	if f.Name == "fetch" {
		g.genFetchStatement(f)
		return
	}
	result := builtins.GenStmt(f.Name, f.Args, f.KwArgs, g.genExpr, g.genRawValue)
	if result.OK {
		g.writeln(result.Code)
		return
	}

	// User-defined function call
	args := make([]string, len(f.Args))
	for i, arg := range f.Args {
		args[i] = g.genExpr(arg)
	}
	g.writeln(fmt.Sprintf("%s %s", f.Name, strings.Join(args, " ")))
}

func (g *Generator) genFuncDecl(f *ast.FuncDecl) {
	g.writeln(fmt.Sprintf("%s() {", f.Name))
	g.indent++

	for i, param := range f.Params {
		if param.Default != nil {
			def := g.genRawValue(param.Default)
			g.writeln(fmt.Sprintf(`local %s="${%d:-%s}"`, param.Name, i+1, def))
		} else {
			g.writeln(fmt.Sprintf(`local %s="$%d"`, param.Name, i+1))
		}
	}

	for _, stmt := range f.Body {
		g.genStatement(stmt)
	}

	g.indent--
	g.writeln("}")
}

func (g *Generator) genIf(i *ast.IfStmt) {
	g.writeln(fmt.Sprintf("if %s; then", g.genCondition(i.Condition)))
	g.genBlock(i.Body)
	g.genElseChain(i.ElseBody)
	g.writeln("fi")
}

func (g *Generator) genElseChain(elseBody []ast.Node) {
	if len(elseBody) == 1 {
		if elif, ok := elseBody[0].(*ast.IfStmt); ok {
			g.writeln(fmt.Sprintf("elif %s; then", g.genCondition(elif.Condition)))
			g.genBlock(elif.Body)
			g.genElseChain(elif.ElseBody)
			return
		}
	}
	if len(elseBody) > 0 {
		g.writeln("else")
		g.genBlock(elseBody)
	}
}

func (g *Generator) genFor(f *ast.ForStmt) {
	collection := g.genForCollection(f.Collection)
	g.writeln(fmt.Sprintf("for %s in %s; do", f.Var, collection))
	g.genBlock(f.Body)
	g.writeln("done")
}

func (g *Generator) genMatch(m *ast.MatchStmt) {
	expr := g.genConditionOperand(m.Expr)
	g.writeln(fmt.Sprintf("case %s in", expr))
	g.indent++

	for _, c := range m.Cases {
		if c.Pattern == nil {
			g.writeln("*)")
		} else {
			g.writeln(fmt.Sprintf("%s)", g.genRawValue(c.Pattern)))
		}
		g.indent++
		for _, stmt := range c.Body {
			g.genStatement(stmt)
		}
		g.writeln(";;")
		g.indent--
	}

	g.indent--
	g.writeln("esac")
}

func (g *Generator) genWhile(w *ast.WhileStmt) {
	g.writeln(fmt.Sprintf("while %s; do", g.genCondition(w.Condition)))
	g.genBlock(w.Body)
	g.writeln("done")
}

func (g *Generator) genIndexAssignment(n *ast.IndexAssignment) {
	val := g.genExpr(n.Value)
	if strIdx, ok := n.Index.(*ast.StringLiteral); ok {
		// Map assignment: config["host"] = "new" → config_host="new"
		g.writeln(fmt.Sprintf(`%s_%s=%s`, n.Object, strIdx.Value, val))
	} else {
		// Array assignment: items[0] = "new" → items[0]="new"
		idx := g.genRawValue(n.Index)
		g.writeln(fmt.Sprintf(`%s[%s]=%s`, n.Object, idx, val))
	}
}

func (g *Generator) genReturn(r *ast.ReturnStmt) {
	if r.Value == nil {
		g.writeln("return")
		return
	}
	g.writeln(fmt.Sprintf("return %s", g.genExpr(r.Value)))
}
