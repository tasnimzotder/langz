package codegen

import (
	"fmt"
	"strings"

	"github.com/tasnimzotder/langz/ast"
	"github.com/tasnimzotder/langz/codegen/builtins"
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
	case *ast.ReturnStmt:
		g.genReturn(n)
	case *ast.ContinueStmt:
		g.writeln("continue")
	}
}

func (g *Generator) genAssignment(a *ast.Assignment) {
	g.writeIndent()
	value := g.genExpr(a.Value)
	g.write(fmt.Sprintf("%s=%s\n", a.Name, value))
}

func (g *Generator) genFuncCallStmt(f *ast.FuncCall) {
	result := builtins.GenStmt(f.Name, f.Args, g.genExpr, g.genRawValue)
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
		g.writeln(fmt.Sprintf(`local %s="$%d"`, param.Name, i+1))
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

	if len(i.ElseBody) > 0 {
		g.writeln("else")
		g.genBlock(i.ElseBody)
	}

	g.writeln("fi")
}

func (g *Generator) genFor(f *ast.ForStmt) {
	collection := g.genForCollection(f.Collection)
	g.writeln(fmt.Sprintf("for %s in %s; do", f.Var, collection))
	g.genBlock(f.Body)
	g.writeln("done")
}

func (g *Generator) genReturn(r *ast.ReturnStmt) {
	if r.Value == nil {
		g.writeln("return")
		return
	}
	g.writeln(fmt.Sprintf("return %s", g.genExpr(r.Value)))
}
