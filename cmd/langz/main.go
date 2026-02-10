package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/tasnimzotder/langz/internal/ast"
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

	// Shebang support: if first arg is a .lz file, treat as "run <file>"
	if strings.HasSuffix(command, ".lz") {
		if _, err := os.Stat(command); err == nil {
			os.Args = append([]string{os.Args[0], "run", command}, os.Args[2:]...)
			command = "run"
		}
	}

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
	prog, parseErrs := parser.New(tokens).ParseAllErrors()
	if len(parseErrs) > 0 {
		formatAllParseErrors(string(source), inputFile, parseErrs)
		os.Exit(1)
	}

	// Resolve imports before codegen
	baseDir := filepath.Dir(inputFile)
	absInput, _ := filepath.Abs(inputFile)
	visited := map[string]bool{absInput: true}
	if err := resolveImports(prog, baseDir, visited); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", inputFile, err)
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

// resolveImports walks the AST, finds ImportStmt nodes, reads/parses imported
// files, and prepends their statements. Circular imports are detected via visited.
func resolveImports(prog *ast.Program, baseDir string, visited map[string]bool) error {
	var resolved []ast.Node
	for _, stmt := range prog.Statements {
		imp, ok := stmt.(*ast.ImportStmt)
		if !ok {
			resolved = append(resolved, stmt)
			continue
		}

		importPath := filepath.Join(baseDir, imp.Path)
		absPath, err := filepath.Abs(importPath)
		if err != nil {
			return fmt.Errorf("import %q: %v", imp.Path, err)
		}

		if visited[absPath] {
			return fmt.Errorf("circular import detected: %s", imp.Path)
		}
		visited[absPath] = true

		data, err := os.ReadFile(importPath)
		if err != nil {
			return fmt.Errorf("import %q: %v", imp.Path, err)
		}

		tokens := lexer.New(string(data)).Tokenize()
		importProg, parseErrs := parser.New(tokens).ParseAllErrors()
		if len(parseErrs) > 0 {
			return fmt.Errorf("import %q: %s", imp.Path, parseErrs[0].Message)
		}

		// Recursively resolve imports in the imported file
		importDir := filepath.Dir(importPath)
		if err := resolveImports(importProg, importDir, visited); err != nil {
			return err
		}

		resolved = append(resolved, importProg.Statements...)
	}
	prog.Statements = resolved
	return nil
}

// formatAllParseErrors prints all parse errors with source context.
func formatAllParseErrors(source string, inputFile string, errs []parser.ParseError) {
	lines := strings.Split(source, "\n")
	maxShow := 10
	for i, e := range errs {
		if i >= maxShow {
			remaining := len(errs) - maxShow
			fmt.Fprintf(os.Stderr, "... and %d more error(s)\n", remaining)
			break
		}
		fmt.Fprintf(os.Stderr, "%s:%d:%d: %s\n", inputFile, e.Line, e.Col, e.Message)
		if e.Line >= 1 && e.Line <= len(lines) {
			srcLine := lines[e.Line-1]
			fmt.Fprintf(os.Stderr, "  %s\n", srcLine)
			if e.Col >= 1 && e.Col <= len(srcLine)+1 {
				fmt.Fprintf(os.Stderr, "  %s^\n", strings.Repeat(" ", e.Col-1))
			}
		}
	}
}
