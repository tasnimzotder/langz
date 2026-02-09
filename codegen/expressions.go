package codegen

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/tasnimzotder/langz/ast"
	"github.com/tasnimzotder/langz/codegen/builtins"
)

var interpRegex = regexp.MustCompile(`\{(\w+)\}`)

// interpolate converts Langz string interpolation {var} to Bash ${var}.
func interpolate(s string) string {
	return interpRegex.ReplaceAllString(s, "${$1}")
}

func (g *Generator) genExpr(node ast.Node) string {
	switch n := node.(type) {
	case *ast.StringLiteral:
		return fmt.Sprintf(`"%s"`, interpolate(n.Value))
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
		obj := g.genExpr(n.Object)
		return fmt.Sprintf("%s.%s", obj, n.Field)
	case *ast.BinaryExpr:
		if isArithmeticOp(n.Op) {
			return fmt.Sprintf("$((%s %s %s))", g.genArithOperand(n.Left), n.Op, g.genArithOperand(n.Right))
		}
		return fmt.Sprintf("%s %s %s", g.genExpr(n.Left), n.Op, g.genExpr(n.Right))
	case *ast.UnaryExpr:
		return fmt.Sprintf("%s%s", n.Op, g.genExpr(n.Operand))
	case *ast.ListLiteral:
		elems := make([]string, len(n.Elements))
		for i, e := range n.Elements {
			elems[i] = g.genExpr(e)
		}
		return fmt.Sprintf("(%s)", strings.Join(elems, " "))
	case *ast.MapLiteral:
		return g.genMapLiteral(n)
	default:
		return ""
	}
}

func (g *Generator) genFuncCallExpr(f *ast.FuncCall) string {
	result := builtins.GenExpr(f.Name, f.Args, g.genExpr, g.genRawValue)
	if result.OK {
		return result.Code
	}

	// User-defined function call
	args := make([]string, len(f.Args))
	for i, arg := range f.Args {
		args[i] = g.genExpr(arg)
	}
	return fmt.Sprintf("%s %s", f.Name, strings.Join(args, " "))
}

// genRawValue extracts the raw value from a node without quoting.
func (g *Generator) genRawValue(node ast.Node) string {
	switch n := node.(type) {
	case *ast.StringLiteral:
		return interpolate(n.Value)
	case *ast.IntLiteral:
		return n.Value
	case *ast.Identifier:
		return fmt.Sprintf("$%s", n.Name)
	default:
		return g.genExpr(node)
	}
}

func (g *Generator) genCondition(node ast.Node) string {
	switch n := node.(type) {
	case *ast.BinaryExpr:
		if n.Op == "and" {
			return fmt.Sprintf("%s && %s", g.genCondition(n.Left), g.genCondition(n.Right))
		}
		if n.Op == "or" {
			return fmt.Sprintf("%s || %s", g.genCondition(n.Left), g.genCondition(n.Right))
		}
		left := g.genConditionOperand(n.Left)
		right := g.genConditionOperand(n.Right)
		op := bashCompareOp(n.Op)
		return fmt.Sprintf("[ %s %s %s ]", left, op, right)
	case *ast.UnaryExpr:
		if n.Op == "!" {
			return fmt.Sprintf("[ ! %s ]", g.genConditionOperand(n.Operand))
		}
		return g.genExpr(node)
	case *ast.FuncCall:
		return g.genFuncCallExpr(n)
	case *ast.Identifier:
		return fmt.Sprintf(`[ "$%s" = true ]`, n.Name)
	default:
		return g.genExpr(node)
	}
}

func (g *Generator) genConditionOperand(node ast.Node) string {
	switch n := node.(type) {
	case *ast.Identifier:
		return fmt.Sprintf(`"$%s"`, n.Name)
	case *ast.IntLiteral:
		return n.Value
	case *ast.StringLiteral:
		return fmt.Sprintf(`"%s"`, interpolate(n.Value))
	default:
		return g.genExpr(node)
	}
}

func (g *Generator) genForCollection(node ast.Node) string {
	switch n := node.(type) {
	case *ast.Identifier:
		return fmt.Sprintf(`"${%s[@]}"`, n.Name)
	case *ast.FuncCall:
		expr := g.genFuncCallExpr(n)
		if strings.HasPrefix(expr, "$(") {
			return expr
		}
		return fmt.Sprintf("$(%s)", expr)
	default:
		return g.genExpr(node)
	}
}

func (g *Generator) genMapLiteral(m *ast.MapLiteral) string {
	// Placeholder — will be implemented as declare -A
	return ""
}

// genArithOperand produces unquoted values suitable for Bash $((...)).
func (g *Generator) genArithOperand(node ast.Node) string {
	switch n := node.(type) {
	case *ast.Identifier:
		return n.Name
	case *ast.IntLiteral:
		return n.Value
	case *ast.BinaryExpr:
		if isArithmeticOp(n.Op) {
			return fmt.Sprintf("%s %s %s", g.genArithOperand(n.Left), n.Op, g.genArithOperand(n.Right))
		}
		return g.genExpr(node)
	default:
		return g.genRawValue(node)
	}
}

func isArithmeticOp(op string) bool {
	return op == "+" || op == "-" || op == "*" || op == "/"
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
