package integration_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestE2E_ExecBuiltin(t *testing.T) {
	source := `result = exec("echo hello from exec")
print(result)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "hello from exec", output)
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

func TestE2E_OsBuiltin(t *testing.T) {
	source := `platform = os()
print(platform)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Contains(t, []string{"darwin", "linux"}, output)
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
