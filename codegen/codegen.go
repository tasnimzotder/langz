package codegen

import (
	"fmt"
	"strings"

	"github.com/tasnimzotder/langz/ast"
)

type generator struct {
	buf    strings.Builder
	indent int
}

func Generate(prog *ast.Program) string {
	g := &generator{}
	g.writeln("#!/bin/bash")
	g.writeln("set -euo pipefail")
	g.writeln("")

	for _, stmt := range prog.Statements {
		g.genStatement(stmt)
	}

	return strings.TrimRight(g.buf.String(), "\n") + "\n"
}

func (g *generator) write(s string) {
	g.buf.WriteString(s)
}

func (g *generator) writeln(s string) {
	g.writeIndent()
	g.buf.WriteString(s)
	g.buf.WriteString("\n")
}

func (g *generator) writeIndent() {
	for i := 0; i < g.indent; i++ {
		g.buf.WriteString("  ")
	}
}

func (g *generator) genStatement(node ast.Node) {
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

func (g *generator) genAssignment(a *ast.Assignment) {
	g.writeIndent()
	g.write(fmt.Sprintf("%s=%s\n", a.Name, g.genExpr(a.Value)))
}

func (g *generator) genExpr(node ast.Node) string {
	switch n := node.(type) {
	case *ast.StringLiteral:
		return fmt.Sprintf(`"%s"`, n.Value)
	case *ast.IntLiteral:
		return n.Value
	case *ast.BoolLiteral:
		if n.Value {
			return "true"
		}
		return "false"
	case *ast.Identifier:
		return fmt.Sprintf(`"$%s"`, n.Name)
	case *ast.FuncCall:
		return g.genFuncCallExpr(n)
	case *ast.DotExpr:
		return g.genDotExpr(n)
	case *ast.BinaryExpr:
		return g.genBinaryExpr(n)
	case *ast.UnaryExpr:
		return g.genUnaryExpr(n)
	default:
		return ""
	}
}

func (g *generator) genFuncCallStmt(f *ast.FuncCall) {
	if f.Name == "print" {
		args := make([]string, len(f.Args))
		for i, arg := range f.Args {
			args[i] = g.genExpr(arg)
		}
		g.writeln(fmt.Sprintf("echo %s", strings.Join(args, " ")))
		return
	}

	g.writeln(g.genFuncCallExpr(f))
}

func (g *generator) genFuncCallExpr(f *ast.FuncCall) string {
	if f.Name == "print" {
		args := make([]string, len(f.Args))
		for i, arg := range f.Args {
			args[i] = g.genExpr(arg)
		}
		return fmt.Sprintf("echo %s", strings.Join(args, " "))
	}

	args := make([]string, len(f.Args))
	for i, arg := range f.Args {
		args[i] = g.genExpr(arg)
	}
	return fmt.Sprintf("%s %s", f.Name, strings.Join(args, " "))
}

func (g *generator) genFuncDecl(f *ast.FuncDecl) {
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

func (g *generator) genIf(i *ast.IfStmt) {
	g.writeln(fmt.Sprintf("if %s; then", g.genCondition(i.Condition)))
	g.indent++

	for _, stmt := range i.Body {
		g.genStatement(stmt)
	}

	g.indent--

	if len(i.ElseBody) > 0 {
		g.writeln("else")
		g.indent++
		for _, stmt := range i.ElseBody {
			g.genStatement(stmt)
		}
		g.indent--
	}

	g.writeln("fi")
}

func (g *generator) genCondition(node ast.Node) string {
	switch n := node.(type) {
	case *ast.BinaryExpr:
		left := g.genConditionOperand(n.Left)
		right := g.genConditionOperand(n.Right)
		op := bashCompareOp(n.Op)
		return fmt.Sprintf("[ %s %s %s ]", left, op, right)
	case *ast.UnaryExpr:
		if n.Op == "!" {
			operand := g.genConditionOperand(n.Operand)
			return fmt.Sprintf("[ ! %s ]", operand)
		}
		return g.genExpr(node)
	case *ast.Identifier:
		return fmt.Sprintf(`[ "$%s" = true ]`, n.Name)
	default:
		return g.genExpr(node)
	}
}

func (g *generator) genConditionOperand(node ast.Node) string {
	switch n := node.(type) {
	case *ast.Identifier:
		return fmt.Sprintf(`"$%s"`, n.Name)
	case *ast.IntLiteral:
		return n.Value
	case *ast.StringLiteral:
		return fmt.Sprintf(`"%s"`, n.Value)
	default:
		return g.genExpr(node)
	}
}

func bashCompareOp(op string) string {
	switch op {
	case ">":
		return "-gt"
	case "<":
		return "-lt"
	case ">=":
		return "-ge"
	case "<=":
		return "-le"
	case "==":
		return "="
	case "!=":
		return "!="
	default:
		return op
	}
}

func (g *generator) genFor(f *ast.ForStmt) {
	collection := g.genForCollection(f.Collection)
	g.writeln(fmt.Sprintf("for %s in %s; do", f.Var, collection))
	g.indent++

	for _, stmt := range f.Body {
		g.genStatement(stmt)
	}

	g.indent--
	g.writeln("done")
}

func (g *generator) genForCollection(node ast.Node) string {
	switch n := node.(type) {
	case *ast.Identifier:
		return fmt.Sprintf(`"${%s[@]}"`, n.Name)
	case *ast.FuncCall:
		return fmt.Sprintf("$(%s)", g.genFuncCallExpr(n))
	default:
		return g.genExpr(node)
	}
}

func (g *generator) genReturn(r *ast.ReturnStmt) {
	if r.Value == nil {
		g.writeln("return")
		return
	}
	g.writeln(fmt.Sprintf("return %s", g.genExpr(r.Value)))
}

func (g *generator) genDotExpr(d *ast.DotExpr) string {
	obj := g.genExpr(d.Object)
	return fmt.Sprintf("%s.%s", obj, d.Field)
}

func (g *generator) genBinaryExpr(b *ast.BinaryExpr) string {
	left := g.genExpr(b.Left)
	right := g.genExpr(b.Right)
	return fmt.Sprintf("%s %s %s", left, b.Op, right)
}

func (g *generator) genUnaryExpr(u *ast.UnaryExpr) string {
	operand := g.genExpr(u.Operand)
	return fmt.Sprintf("%s%s", u.Op, operand)
}
