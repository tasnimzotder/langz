package codegen

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/tasnimzotder/langz/internal/ast"
	"github.com/tasnimzotder/langz/internal/codegen/builtins"
)

var interpRegex = regexp.MustCompile(`\{(\w+)\}`)

// interpolate converts Langz string interpolation {var} to Bash ${var}.
func interpolate(s string) string {
	return interpRegex.ReplaceAllString(s, "${$1}")
}

// bashEscape escapes characters that are special inside Bash double quotes.
func bashEscape(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return s
}

func (g *Generator) genExpr(node ast.Node) string {
	switch n := node.(type) {
	case *ast.StringLiteral:
		return fmt.Sprintf(`"%s"`, interpolate(bashEscape(n.Value)))
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
		if n.Op == "|>" {
			return g.genPipeExpr(n)
		}
		if isArithmeticOp(n.Op) {
			prec := arithPrecedence(n.Op)
			return fmt.Sprintf("$((%s %s %s))", g.genArithOperandPrec(n.Left, prec), n.Op, g.genArithOperandPrec(n.Right, prec))
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
	case *ast.IndexExpr:
		return g.genIndexExpr(n)
	case *ast.MethodCall:
		return g.genMethodCall(n)
	case *ast.MapLiteral:
		return g.genMapLiteral(n)
	default:
		return ""
	}
}

func (g *Generator) genFuncCallExpr(f *ast.FuncCall) string {
	result := builtins.GenExpr(f.Name, f.Args, f.KwArgs, g.genExpr, g.genRawValue)
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

// genPipeExpr transforms a |> f into f(a) by synthesizing a FuncCall.
func (g *Generator) genPipeExpr(expr *ast.BinaryExpr) string {
	switch right := expr.Right.(type) {
	case *ast.Identifier:
		// data |> upper → upper(data)
		call := &ast.FuncCall{Name: right.Name, Args: []ast.Node{expr.Left}}
		return g.genFuncCallExpr(call)
	case *ast.FuncCall:
		// data |> json_get(".name") → json_get(data, ".name")
		args := make([]ast.Node, 0, 1+len(right.Args))
		args = append(args, expr.Left)
		args = append(args, right.Args...)
		call := &ast.FuncCall{Name: right.Name, Args: args, KwArgs: right.KwArgs}
		return g.genFuncCallExpr(call)
	default:
		return "# error: pipe target must be a function"
	}
}

func (g *Generator) genIndexExpr(n *ast.IndexExpr) string {
	obj := g.genVarName(n.Object)
	switch idx := n.Index.(type) {
	case *ast.StringLiteral:
		// Map access with string key: config["host"] → "$config_host"
		return fmt.Sprintf(`"$%s_%s"`, obj, idx.Value)
	default:
		// Array access: items[0] or items[i] → "${items[idx]}"
		return fmt.Sprintf(`"${%s[%s]}"`, obj, g.genRawValue(n.Index))
	}
}

// genVarName extracts the bare variable name from a node (no $ prefix).
func (g *Generator) genVarName(node ast.Node) string {
	if id, ok := node.(*ast.Identifier); ok {
		return id.Name
	}
	return g.genRawValue(node)
}

func (g *Generator) genMethodCall(m *ast.MethodCall) string {
	obj := g.genVarName(m.Object)
	switch m.Method {
	case "replace":
		if len(m.Args) != 2 {
			return "# error: replace() requires 2 arguments (old, new)"
		}
		old := g.genRawValue(m.Args[0])
		newVal := g.genRawValue(m.Args[1])
		return fmt.Sprintf(`"${%s//%s/%s}"`, obj, old, newVal)
	case "contains":
		if len(m.Args) != 1 {
			return "# error: contains() requires 1 argument"
		}
		sub := g.genRawValue(m.Args[0])
		return fmt.Sprintf(`[[ "$%s" == *"%s"* ]]`, obj, sub)
	case "starts_with":
		if len(m.Args) != 1 {
			return "# error: starts_with() requires 1 argument"
		}
		prefix := g.genRawValue(m.Args[0])
		return fmt.Sprintf(`[[ "$%s" == "%s"* ]]`, obj, prefix)
	case "ends_with":
		if len(m.Args) != 1 {
			return "# error: ends_with() requires 1 argument"
		}
		suffix := g.genRawValue(m.Args[0])
		return fmt.Sprintf(`[[ "$%s" == *"%s" ]]`, obj, suffix)
	default:
		return fmt.Sprintf("# error: unknown method %s", m.Method)
	}
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
		return fmt.Sprintf(`"%s"`, interpolate(bashEscape(n.Value)))
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
	return g.genArithOperandPrec(node, 0)
}

// genArithOperandPrec adds parentheses when a lower-precedence op is nested
// inside a higher-precedence context (e.g., (a + b) * c).
func (g *Generator) genArithOperandPrec(node ast.Node, parentPrec int) string {
	switch n := node.(type) {
	case *ast.Identifier:
		return n.Name
	case *ast.IntLiteral:
		return n.Value
	case *ast.BinaryExpr:
		if isArithmeticOp(n.Op) {
			prec := arithPrecedence(n.Op)
			inner := fmt.Sprintf("%s %s %s", g.genArithOperandPrec(n.Left, prec), n.Op, g.genArithOperandPrec(n.Right, prec))
			if parentPrec > prec {
				return fmt.Sprintf("(%s)", inner)
			}
			return inner
		}
		return g.genExpr(node)
	default:
		return g.genRawValue(node)
	}
}

func arithPrecedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/", "%":
		return 2
	default:
		return 0
	}
}

func isArithmeticOp(op string) bool {
	return op == "+" || op == "-" || op == "*" || op == "/" || op == "%"
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
