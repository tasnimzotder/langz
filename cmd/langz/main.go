package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/tasnimzotder/langz/internal/codegen"
	"github.com/tasnimzotder/langz/internal/lexer"
	"github.com/tasnimzotder/langz/internal/lsp"
	"github.com/tasnimzotder/langz/internal/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: langz <build|run|fmt|lsp> <file.lz>")
		os.Exit(1)
	}

	command := os.Args[1]

	// LSP server needs no file argument
	if command == "lsp" {
		lsp.NewServer().Run()
		return
	}

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: langz <build|run|fmt> <file.lz>")
		os.Exit(1)
	}

	inputFile := os.Args[2]

	source, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", inputFile, err)
		os.Exit(1)
	}

	if command == "fmt" {
		formatted := lsp.FormatSource(string(source), 4, true)
		if formatted == string(source) {
			return
		}
		if err := os.WriteFile(inputFile, []byte(formatted), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", inputFile, err)
			os.Exit(1)
		}
		fmt.Printf("Formatted %s\n", inputFile)
		return
	}

	tokens := lexer.New(string(source)).Tokenize()
	prog, parseErr := parser.New(tokens).ParseWithErrors()
	if parseErr != nil {
		formatParseError(string(source), inputFile, parseErr)
		os.Exit(1)
	}
	output, codegenErrors := codegen.Generate(prog)
	if len(codegenErrors) > 0 {
		for _, e := range codegenErrors {
			fmt.Fprintf(os.Stderr, "%s: %s\n", inputFile, e)
		}
		os.Exit(1)
	}

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

		if _, err := tmpFile.WriteString(output); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing temp file: %v\n", err)
			os.Exit(1)
		}
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
		fmt.Fprintf(os.Stderr, "Unknown command: %s\nUsage: langz <build|run|fmt|lsp> <file.lz>\n", command)
		os.Exit(1)
	}
}

// formatParseError prints a parse error with source context and a ^ pointer.
func formatParseError(source string, inputFile string, parseErr error) {
	// Try to extract line/col from the structured error format
	var line, col int
	var msg string
	_, scanErr := fmt.Sscanf(parseErr.Error(), "line %d, col %d: ", &line, &col)
	if scanErr != nil {
		// Fallback: just print the raw error
		fmt.Fprintf(os.Stderr, "%s: %v\n", inputFile, parseErr)
		return
	}

	// Extract the message after "line X, col Y: "
	errStr := parseErr.Error()
	if idx := strings.Index(errStr, ": "); idx >= 0 {
		msg = errStr[idx+2:]
	}

	lines := strings.Split(source, "\n")
	fmt.Fprintf(os.Stderr, "%s:%d:%d: %s\n", inputFile, line, col, msg)
	if line >= 1 && line <= len(lines) {
		srcLine := lines[line-1]
		fmt.Fprintf(os.Stderr, "  %s\n", srcLine)
		if col >= 1 && col <= len(srcLine)+1 {
			fmt.Fprintf(os.Stderr, "  %s^\n", strings.Repeat(" ", col-1))
		}
	}
}
