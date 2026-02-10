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
// Returns the generated Bash and any codegen errors found.
func Generate(prog *ast.Program) (string, []string) {
	g := &Generator{}
	g.writeln("#!/bin/bash")
	g.writeln("set -euo pipefail")
	g.writeln("")

	for _, stmt := range prog.Statements {
		g.genStatement(stmt)
	}

	output := strings.TrimRight(g.buf.String(), "\n") + "\n"
	return output, findCodegenErrors(output)
}

// findCodegenErrors scans generated Bash for # error: markers
// left by builtins that received invalid arguments.
func findCodegenErrors(output string) []string {
	var errs []string
	for _, line := range strings.Split(output, "\n") {
		trimmed := strings.TrimSpace(line)
		if idx := strings.Index(trimmed, "# error:"); idx >= 0 {
			errs = append(errs, strings.TrimSpace(trimmed[idx+len("# error:"):]))
		}
	}
	return errs
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
