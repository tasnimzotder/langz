package integration_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestE2E_CLIBuild(t *testing.T) {
	// Test the CLI binary itself
	tmpDir := t.TempDir()
	lzFile := filepath.Join(tmpDir, "test.lz")
	shFile := filepath.Join(tmpDir, "test.sh")

	err := os.WriteFile(lzFile, []byte(`print("cli works")`), 0644)
	require.NoError(t, err)

	// Build using go run
	root := projectRoot(t)
	cmd := exec.Command("go", "run", "./cmd/langz", "build", lzFile)
	cmd.Dir = root
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

	root := projectRoot(t)
	cmd := exec.Command("go", "run", "./cmd/langz", "run", lzFile)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "run failed: %s", string(out))

	assert.Equal(t, "run works", strings.TrimSpace(string(out)))
}
