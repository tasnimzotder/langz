package integration_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tasnimzotder/langz/internal/codegen"
	"github.com/tasnimzotder/langz/internal/lexer"
	"github.com/tasnimzotder/langz/internal/parser"
)

func TestE2E_BashBlockEcho(t *testing.T) {
	source := `bash { echo "hello from bash" }`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "hello from bash", output)
}

func TestE2E_BashBlockMultiline(t *testing.T) {
	source := "bash {\n    X=42\n    echo \"val=$X\"\n}"
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "val=42", output)
}

func TestE2E_BashBlockWithLangz(t *testing.T) {
	source := "name = \"world\"\nbash { echo \"inline bash\" }\nprint(\"hello {name}\")"
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	lines := strings.Split(output, "\n")
	assert.Equal(t, "inline bash", lines[0])
	assert.Equal(t, "hello world", lines[1])
}

func TestE2E_BashBlockNestedBraces(t *testing.T) {
	source := "bash {\n    if [ 1 -eq 1 ]; then\n        echo \"nested\"\n    fi\n}"
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "nested", output)
}

func TestE2E_ShebangCompiles(t *testing.T) {
	source := "#!/usr/bin/env langz\nprint(\"shebang works\")"
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "shebang works", output)
}

func TestE2E_ImportParsesCorrectly(t *testing.T) {
	source := `import "helpers.lz"`
	tokens := lexer.New(source).Tokenize()
	prog, errs := parser.New(tokens).ParseAllErrors()

	require.Empty(t, errs)
	require.Len(t, prog.Statements, 1)

	// ImportStmt codegen is a no-op (resolved before codegen)
	_, codegenErrs := codegen.Generate(prog)
	assert.Empty(t, codegenErrs)
}

func TestE2E_ImportViaCLI(t *testing.T) {
	root := projectRoot(t)

	// Create a library file
	libDir := t.TempDir()
	libFile := libDir + "/helpers.lz"
	mainFile := libDir + "/main.lz"

	require.NoError(t, os.WriteFile(libFile, []byte(`fn greet(name: str) {
    print("Hello {name}")
}
`), 0644))

	require.NoError(t, os.WriteFile(mainFile, []byte(`import "helpers.lz"
greet("World")
`), 0644))

	cmd := exec.Command("go", "run", root+"/cmd/langz", "run", mainFile)
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "run failed: %s", string(out))
	assert.Equal(t, "Hello World", strings.TrimSpace(string(out)))
}

func TestE2E_CircularImportDetected(t *testing.T) {
	root := projectRoot(t)
	dir := t.TempDir()
	aFile := dir + "/a.lz"
	bFile := dir + "/b.lz"

	require.NoError(t, os.WriteFile(aFile, []byte(`import "b.lz"
print("a")
`), 0644))
	require.NoError(t, os.WriteFile(bFile, []byte(`import "a.lz"
print("b")
`), 0644))

	cmd := exec.Command("go", "run", root+"/cmd/langz", "run", aFile)
	out, err := cmd.CombinedOutput()
	require.Error(t, err, "should fail on circular import")
	assert.Contains(t, string(out), "circular import")
}

func TestE2E_ShebangAutoDetect(t *testing.T) {
	root := projectRoot(t)
	dir := t.TempDir()
	scriptFile := dir + "/test.lz"

	require.NoError(t, os.WriteFile(scriptFile, []byte("#!/usr/bin/env langz\nprint(\"shebang auto\")\n"), 0644))

	// langz test.lz (auto-detect, no "run" subcommand)
	cmd := exec.Command("go", "run", root+"/cmd/langz", scriptFile)
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "auto-detect failed: %s", string(out))
	assert.Equal(t, "shebang auto", strings.TrimSpace(string(out)))
}
