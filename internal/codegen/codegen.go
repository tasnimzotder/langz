package codegen

import (
	"strings"

	"github.com/tasnimzotder/langz/internal/ast"
)

// Generator holds state for Bash code generation.
type Generator struct {
	buf    strings.Builder
	indent int
}

// Generate converts an AST program into a Bash script string.
func Generate(prog *ast.Program) string {
	g := &Generator{}
	g.writeln("#!/bin/bash")
	g.writeln("set -euo pipefail")
	g.writeln("")

	for _, stmt := range prog.Statements {
		g.genStatement(stmt)
	}

	return strings.TrimRight(g.buf.String(), "\n") + "\n"
}

func (g *Generator) write(s string) {
	g.buf.WriteString(s)
}

func (g *Generator) writeln(s string) {
	g.writeIndent()
	g.buf.WriteString(s)
	g.buf.WriteString("\n")
}

func (g *Generator) writeIndent() {
	for i := 0; i < g.indent; i++ {
		g.buf.WriteString("  ")
	}
}

func (g *Generator) genBlock(stmts []ast.Node) {
	g.indent++
	for _, stmt := range stmts {
		g.genStatement(stmt)
	}
	g.indent--
}
