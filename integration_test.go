package langz_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tasnimzotder/langz/codegen"
	"github.com/tasnimzotder/langz/lexer"
	"github.com/tasnimzotder/langz/parser"
)

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

func TestE2E_HelloWorld(t *testing.T) {
	source := `print("Hello, World!")`
	bash := compileSource(t, source)

	assert.Contains(t, bash, "#!/bin/bash")
	assert.Contains(t, bash, "set -euo pipefail")

	output, code := runBash(t, bash)
	assert.Equal(t, 0, code)
	assert.Equal(t, "Hello, World!", output)
}

func TestE2E_Variables(t *testing.T) {
	source := `
name = "Langz"
version = 42
print("Welcome to {name}")
print(version)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	lines := strings.Split(output, "\n")
	assert.Equal(t, "Welcome to Langz", lines[0])
	assert.Equal(t, "42", lines[1])
}

func TestE2E_Function(t *testing.T) {
	source := `
fn greet(name: str) {
	print("Hello {name}")
}
greet("World")
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "Hello World", output)
}

func TestE2E_IfElse(t *testing.T) {
	source := `
x = 5
if x > 3 {
	print("big")
} else {
	print("small")
}
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "big", output)
}

func TestE2E_ForLoop(t *testing.T) {
	source := `
items = ["a", "b", "c"]
for item in items {
	print(item)
}
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	lines := strings.Split(output, "\n")
	assert.Equal(t, []string{"a", "b", "c"}, lines)
}

func TestE2E_RangeLoop(t *testing.T) {
	source := `
for i in range(1, 3) {
	print(i)
}
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	lines := strings.Split(output, "\n")
	assert.Equal(t, []string{"1", "2", "3"}, lines)
}

func TestE2E_MatchStatement(t *testing.T) {
	source := `
platform = "linux"
match platform {
	"darwin" => print("macOS")
	"linux" => print("Linux")
	_ => print("unknown")
}
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "Linux", output)
}

func TestE2E_ExecBuiltin(t *testing.T) {
	source := `result = exec("echo hello from exec")
print(result)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "hello from exec", output)
}

func TestE2E_EnvWithDefault(t *testing.T) {
	source := `name = env("LANGZ_TEST_UNSET_VAR") or "fallback"
print(name)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "fallback", output)
}

func TestE2E_FileOps(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")

	source := `
mkdir("` + tmpDir + `/sub")
write("` + filePath + `", "hello langz")
content = read("` + filePath + `")
print(content)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "hello langz", output)

	// Verify mkdir worked
	_, err := os.Stat(filepath.Join(tmpDir, "sub"))
	assert.NoError(t, err)
}

func TestE2E_ExitCode(t *testing.T) {
	source := `exit(42)`
	bash := compileSource(t, source)
	_, code := runBash(t, bash)

	assert.Equal(t, 42, code)
}

func TestE2E_Comments(t *testing.T) {
	source := `
// This is a comment
name = "test"
// Another comment
print(name)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "test", output)
}

func TestE2E_OsBuiltin(t *testing.T) {
	source := `platform = os()
print(platform)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Contains(t, []string{"darwin", "linux"}, output)
}

func TestE2E_ComprehensiveScript(t *testing.T) {
	tmpDir := t.TempDir()
	confFile := filepath.Join(tmpDir, "app.conf")

	source := `
// Comprehensive Langz script
app = "myapp"
version = "2.0"

fn setup(dir: str) {
	mkdir(dir)
	print("Setup done for {app}")
}

setup("` + tmpDir + `/build")

write("` + confFile + `", "port=8080")
config = read("` + confFile + `")
print(config)

items = ["alpha", "beta"]
for item in items {
	print(item)
}

mode = "production"
match mode {
	"development" => print("dev mode")
	"production" => print("prod mode")
	_ => print("unknown mode")
}
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)

	lines := strings.Split(output, "\n")
	require.GreaterOrEqual(t, len(lines), 5)
	assert.Equal(t, "Setup done for myapp", lines[0])
	assert.Equal(t, "port=8080", lines[1])
	assert.Equal(t, "alpha", lines[2])
	assert.Equal(t, "beta", lines[3])
	assert.Equal(t, "prod mode", lines[4])
}

func TestE2E_CLIBuild(t *testing.T) {
	// Test the CLI binary itself
	tmpDir := t.TempDir()
	lzFile := filepath.Join(tmpDir, "test.lz")
	shFile := filepath.Join(tmpDir, "test.sh")

	err := os.WriteFile(lzFile, []byte(`print("cli works")`), 0644)
	require.NoError(t, err)

	// Build using go run
	cmd := exec.Command("go", "run", "./cmd/langz", "build", lzFile)
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "build failed: %s", string(out))

	// Verify .sh file was created
	_, err = os.Stat(shFile)
	require.NoError(t, err, "expected %s to exist", shFile)

	// Run the generated script
	output, code := runBash(t, mustReadFile(t, shFile))
	assert.Equal(t, 0, code)
	assert.Equal(t, "cli works", output)
}

func TestE2E_CLIRun(t *testing.T) {
	tmpDir := t.TempDir()
	lzFile := filepath.Join(tmpDir, "test.lz")

	err := os.WriteFile(lzFile, []byte(`print("run works")`), 0644)
	require.NoError(t, err)

	cmd := exec.Command("go", "run", "./cmd/langz", "run", lzFile)
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "run failed: %s", string(out))

	assert.Equal(t, "run works", strings.TrimSpace(string(out)))
}

func TestE2E_FunctionWithMultipleStatements(t *testing.T) {
	source := `
fn deploy(target: str) {
	print("deploying to {target}")
	print("done")
}
deploy("prod")
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	lines := strings.Split(output, "\n")
	assert.Equal(t, "deploying to prod", lines[0])
	assert.Equal(t, "done", lines[1])
}

func TestE2E_NestedIfInFor(t *testing.T) {
	source := `
scores = [90, 45, 80]
for s in scores {
	if s > 50 {
		print(s)
	}
}
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	lines := strings.Split(output, "\n")
	assert.Equal(t, []string{"90", "80"}, lines)
}

func TestE2E_MultipleFunctions(t *testing.T) {
	source := `
fn add_prefix(s: str) {
	print("[INFO] {s}")
}

fn log_items(prefix: str) {
	add_prefix("starting {prefix}")
	add_prefix("done {prefix}")
}

log_items("deploy")
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	lines := strings.Split(output, "\n")
	assert.Equal(t, "[INFO] starting deploy", lines[0])
	assert.Equal(t, "[INFO] done deploy", lines[1])
}

func TestE2E_ForWithRange(t *testing.T) {
	source := `
total = 0
for i in range(1, 5) {
	print(i)
}
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	lines := strings.Split(output, "\n")
	assert.Equal(t, []string{"1", "2", "3", "4", "5"}, lines)
}

func TestE2E_ExecPipeline(t *testing.T) {
	source := `
count = exec("echo -e 'a\nb\nc' | wc -l")
print(count)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Contains(t, output, "3")
}

func TestE2E_MatchWithFunctionCalls(t *testing.T) {
	source := `
action = "build"
match action {
	"test" => print("running tests")
	"build" => {
		print("compiling")
		print("linking")
	}
	_ => print("unknown action")
}
`
	bash := compileSource(t, source)

	// Verify the generated code structure
	assert.Contains(t, bash, `case "$action" in`)
	assert.Contains(t, bash, "build)")
	assert.Contains(t, bash, "esac")
}

func TestE2E_FileWorkflow(t *testing.T) {
	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "source.txt")
	dstFile := filepath.Join(tmpDir, "backup.txt")

	source := `
write("` + srcFile + `", "important data")
copy("` + srcFile + `", "` + dstFile + `")
original = read("` + srcFile + `")
backup = read("` + dstFile + `")
print(original)
print(backup)
rm("` + srcFile + `")
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	lines := strings.Split(output, "\n")
	assert.Equal(t, "important data", lines[0])
	assert.Equal(t, "important data", lines[1])

	// Verify rm worked
	_, err := os.Stat(srcFile)
	assert.True(t, os.IsNotExist(err))

	// Verify backup still exists
	_, err = os.Stat(dstFile)
	assert.NoError(t, err)
}

func TestE2E_EnvOrExitFallback(t *testing.T) {
	source := `
name = env("LANGZ_E2E_TEST_VAR") or "default_val"
print(name)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "default_val", output)
}

func TestE2E_StringInterpolationComplex(t *testing.T) {
	source := `
host = "localhost"
port = 8080
print("Server at {host}:{port}")
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "Server at localhost:8080", output)
}

func TestE2E_BooleanCondition(t *testing.T) {
	source := `
verbose = true
if verbose {
	print("debug on")
}
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "debug on", output)
}

func TestE2E_ChainedExec(t *testing.T) {
	source := `
hostname = exec("echo testhost")
kernel = exec("echo 5.15")
print("Host: {hostname}")
print("Kernel: {kernel}")
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	lines := strings.Split(output, "\n")
	assert.Equal(t, "Host: testhost", lines[0])
	assert.Equal(t, "Kernel: 5.15", lines[1])
}

func TestE2E_DeploymentScript(t *testing.T) {
	// Simulates a real deployment automation script
	tmpDir := t.TempDir()

	source := `
// Deployment script
app = "webapp"
env_name = "staging"

fn log(msg: str) {
	print("[{env_name}] {msg}")
}

// Setup build directory
mkdir("` + tmpDir + `/dist")

log("Starting deploy of {app}")

// Write build artifact
write("` + tmpDir + `/dist/app.js", "console.log('hello')")

// Verify it exists
if exists("` + tmpDir + `/dist/app.js") {
	log("Build artifact ready")
} else {
	log("Build failed")
	exit(1)
}

// Platform-specific message
platform = os()
match platform {
	"darwin" => log("Deploying from macOS")
	"linux" => log("Deploying from Linux")
	_ => log("Deploying from unknown OS")
}

log("Deploy complete")
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	lines := strings.Split(output, "\n")

	assert.Equal(t, "[staging] Starting deploy of webapp", lines[0])
	assert.Equal(t, "[staging] Build artifact ready", lines[1])
	// lines[2] is platform-specific
	assert.Contains(t, lines[2], "[staging] Deploying from")
	assert.Equal(t, "[staging] Deploy complete", lines[3])
}

func TestE2E_BatchFileProcessor(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	for _, name := range []string{"a.txt", "b.txt", "c.txt"} {
		os.WriteFile(filepath.Join(tmpDir, name), []byte("content of "+name), 0644)
	}

	source := `
files = ["a.txt", "b.txt", "c.txt"]
for f in files {
	content = read("` + tmpDir + `/{f}")
	print(content)
}
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	lines := strings.Split(output, "\n")
	assert.Equal(t, "content of a.txt", lines[0])
	assert.Equal(t, "content of b.txt", lines[1])
	assert.Equal(t, "content of c.txt", lines[2])
}

func TestE2E_EqualityComparison(t *testing.T) {
	source := `
mode = "prod"
if mode == "prod" {
	print("production")
}
`
	bash := compileSource(t, source)
	// Verify it generates string comparison
	assert.Contains(t, bash, `[ "$mode" = "prod" ]`)
}

func TestE2E_WhileWithBreak(t *testing.T) {
	source := `
while true {
	print("once")
	break
}
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "once", output)
}

func TestE2E_Arithmetic(t *testing.T) {
	source := `
a = 10
b = 3
result = a + b
print(result)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "13", output)
}

func TestE2E_SleepBuiltin(t *testing.T) {
	source := `sleep(0)`
	bash := compileSource(t, source)
	_, code := runBash(t, bash)

	assert.Equal(t, 0, code)
}

func TestE2E_AppendBuiltin(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "log.txt")

	source := `
write("` + logFile + `", "line1")
append("` + logFile + `", "line2")
content = read("` + logFile + `")
print(content)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Contains(t, output, "line1")
	assert.Contains(t, output, "line2")
}

func TestE2E_SystemInfoBuiltins(t *testing.T) {
	source := `
host = hostname()
user = whoami()
a = arch()
print(host)
print(user)
print(a)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	lines := strings.Split(output, "\n")
	assert.Len(t, lines, 3)
	// All should be non-empty
	for _, line := range lines {
		assert.NotEmpty(t, line)
	}
}

func TestE2E_PathBuiltins(t *testing.T) {
	source := `
dir = dirname("/usr/local/bin/langz")
name = basename("/usr/local/bin/langz")
print(dir)
print(name)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	lines := strings.Split(output, "\n")
	assert.Equal(t, "/usr/local/bin", lines[0])
	assert.Equal(t, "langz", lines[1])
}

func TestE2E_UpperLower(t *testing.T) {
	source := `
up = upper("hello")
low = lower("WORLD")
print(up)
print(low)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	lines := strings.Split(output, "\n")
	assert.Equal(t, "HELLO", lines[0])
	assert.Equal(t, "world", lines[1])
}

func TestE2E_SREMonitoringScript(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "health.log")

	source := `
// SRE health check script
host = hostname()
user = whoami()
platform = os()

write("` + logFile + `", "Health Check Report")
append("` + logFile + `", "Host: {host}")
append("` + logFile + `", "User: {user}")
append("` + logFile + `", "OS: {platform}")

report = read("` + logFile + `")
print(report)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Contains(t, output, "Health Check Report")
	assert.Contains(t, output, "Host:")
	assert.Contains(t, output, "User:")
	assert.Contains(t, output, "OS:")
}

func mustReadFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	return string(data)
}
