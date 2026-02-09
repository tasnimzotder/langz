package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/tasnimzotder/langz/codegen"
	"github.com/tasnimzotder/langz/lexer"
	"github.com/tasnimzotder/langz/parser"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: langz <build|run> <file.lz>")
		os.Exit(1)
	}

	command := os.Args[1]
	inputFile := os.Args[2]

	source, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", inputFile, err)
		os.Exit(1)
	}

	tokens := lexer.New(string(source)).Tokenize()
	prog := parser.New(tokens).Parse()
	output := codegen.Generate(prog)

	switch command {
	case "build":
		outFile := strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + ".sh"
		err := os.WriteFile(outFile, []byte(output), 0755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", outFile, err)
			os.Exit(1)
		}
		fmt.Printf("Built %s -> %s\n", inputFile, outFile)

	case "run":
		tmpFile, err := os.CreateTemp("", "langz-*.sh")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating temp file: %v\n", err)
			os.Exit(1)
		}
		defer os.Remove(tmpFile.Name())

		tmpFile.WriteString(output)
		tmpFile.Close()

		cmd := exec.Command("bash", tmpFile.Name())
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		if err := cmd.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			}
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\nUsage: langz <build|run> <file.lz>\n", command)
		os.Exit(1)
	}
}
