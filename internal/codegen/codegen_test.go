package codegen

import (
	"strings"

	"github.com/tasnimzotder/langz/internal/lexer"
	"github.com/tasnimzotder/langz/internal/parser"
)

func compile(input string) string {
	tokens := lexer.New(input).Tokenize()
	prog, err := parser.New(tokens).ParseWithErrors()
	if err != nil {
		panic(err.Error())
	}
	output, _ := Generate(prog)
	return output
}

func compileWithErrors(input string) (string, []string) {
	tokens := lexer.New(input).Tokenize()
	prog, err := parser.New(tokens).ParseWithErrors()
	if err != nil {
		panic(err.Error())
	}
	return Generate(prog)
}

func body(output string) string {
	// Strip the preamble (#!/bin/bash and set -euo pipefail)
	lines := strings.Split(output, "\n")
	var result []string
	for _, line := range lines {
		if line == "#!/bin/bash" || line == "set -euo pipefail" || line == "" {
			continue
		}
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}
