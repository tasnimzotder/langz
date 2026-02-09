package integration_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tasnimzotder/langz/internal/codegen"
	"github.com/tasnimzotder/langz/internal/lexer"
	"github.com/tasnimzotder/langz/internal/parser"
)

// projectRoot returns the module root by walking up from the test directory.
func projectRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	require.NoError(t, err)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		require.NotEqual(t, dir, parent, "could not find go.mod")
		dir = parent
	}
}

func compileSource(t *testing.T, source string) string {
	t.Helper()
	tokens := lexer.New(source).Tokenize()
	prog, err := parser.New(tokens).ParseWithErrors()
	require.NoError(t, err, "parse error")
	return codegen.Generate(prog)
}

func runBash(t *testing.T, script string) (string, int) {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "langz-test-*.sh")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(script)
	require.NoError(t, err)
	tmpFile.Close()

	cmd := exec.Command("bash", tmpFile.Name())
	out, err := cmd.CombinedOutput()

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
	}

	return strings.TrimSpace(string(out)), exitCode
}

func mustReadFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	return string(data)
}
