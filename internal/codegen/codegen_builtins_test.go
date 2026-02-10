package codegen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecBuiltin(t *testing.T) {
	output := body(compile(`result = exec("ls -la")`))

	assert.Contains(t, output, `result=$(ls -la)`)
}

func TestEnvBuiltin(t *testing.T) {
	output := body(compile(`home = env("HOME")`))

	assert.Contains(t, output, `home="${HOME}"`)
}

func TestExistsBuiltin(t *testing.T) {
	output := body(compile(`if exists("file.txt") { print("found") }`))

	assert.Contains(t, output, `[ -e "file.txt" ]`)
}

func TestReadBuiltin(t *testing.T) {
	output := body(compile(`content = read("file.txt")`))

	assert.Contains(t, output, `content=$(cat "file.txt")`)
}

func TestWriteBuiltin(t *testing.T) {
	output := body(compile(`write("out.txt", "hello")`))

	assert.Contains(t, output, `echo "hello" > "out.txt"`)
}

func TestRmBuiltin(t *testing.T) {
	output := body(compile(`rm("temp.txt")`))

	assert.Contains(t, output, `rm -f "temp.txt"`)
}

func TestMkdirBuiltin(t *testing.T) {
	output := body(compile(`mkdir("build/output")`))

	assert.Contains(t, output, `mkdir -p "build/output"`)
}

func TestCopyBuiltin(t *testing.T) {
	output := body(compile(`copy("a.txt", "b.txt")`))

	assert.Contains(t, output, `cp "a.txt" "b.txt"`)
}

func TestMoveBuiltin(t *testing.T) {
	output := body(compile(`move("old.txt", "new.txt")`))

	assert.Contains(t, output, `mv "old.txt" "new.txt"`)
}

func TestChmodBuiltin(t *testing.T) {
	output := body(compile(`chmod("script.sh", "755")`))

	assert.Contains(t, output, `chmod 755 "script.sh"`)
}

func TestGlobBuiltin(t *testing.T) {
	output := body(compile(`files = glob("*.log")`))

	assert.Contains(t, output, `files=(*.log)`)
}

func TestExitBuiltin(t *testing.T) {
	output := body(compile(`exit(1)`))

	assert.Contains(t, output, "exit 1")
}

func TestFetchBuiltin(t *testing.T) {
	output := body(compile(`res = fetch("https://api.example.com/health")`))

	assert.Contains(t, output, `curl -s -w "%{http_code}"`)
	assert.Contains(t, output, `"https://api.example.com/health"`)
	assert.Contains(t, output, `res="$_body"`)
}

func TestSleepBuiltin(t *testing.T) {
	output := body(compile(`sleep(5)`))

	assert.Contains(t, output, "sleep 5")
}

func TestAppendBuiltin(t *testing.T) {
	output := body(compile(`append("log.txt", "entry")`))

	assert.Contains(t, output, `echo "entry" >> "log.txt"`)
}

func TestHostnameBuiltin(t *testing.T) {
	output := body(compile(`host = hostname()`))

	assert.Contains(t, output, `host=$(hostname)`)
}

func TestWhoamiBuiltin(t *testing.T) {
	output := body(compile(`user = whoami()`))

	assert.Contains(t, output, `user=$(whoami)`)
}

func TestArchBuiltin(t *testing.T) {
	output := body(compile(`a = arch()`))

	assert.Contains(t, output, `a=$(uname -m)`)
}

func TestDirnameBuiltin(t *testing.T) {
	output := body(compile(`dir = dirname("/path/to/file.txt")`))

	assert.Contains(t, output, `dir=$(dirname "/path/to/file.txt")`)
}

func TestBasenameBuiltin(t *testing.T) {
	output := body(compile(`name = basename("/path/to/file.txt")`))

	assert.Contains(t, output, `name=$(basename "/path/to/file.txt")`)
}

func TestIsFileBuiltin(t *testing.T) {
	output := body(compile(`if is_file("test.txt") { print("file") }`))

	assert.Contains(t, output, `[ -f "test.txt" ]`)
}

func TestIsDirBuiltin(t *testing.T) {
	output := body(compile(`if is_dir("build") { print("dir") }`))

	assert.Contains(t, output, `[ -d "build" ]`)
}

func TestRmdirBuiltin(t *testing.T) {
	output := body(compile(`rmdir("build")`))

	assert.Contains(t, output, `rm -rf "build"`)
}

func TestUpperBuiltin(t *testing.T) {
	output := body(compile(`x = upper("hello")`))

	assert.Contains(t, output, `x=$(echo "hello" | tr '[:lower:]' '[:upper:]')`)
}

func TestLowerBuiltin(t *testing.T) {
	output := body(compile(`x = lower("HELLO")`))

	assert.Contains(t, output, `x=$(echo "HELLO" | tr '[:upper:]' '[:lower:]')`)
}

func TestOsBuiltin(t *testing.T) {
	output := body(compile(`platform = os()`))

	assert.Contains(t, output, `platform=$(uname -s | tr '[:upper:]' '[:lower:]')`)
}

func TestArgsBuiltin(t *testing.T) {
	output := body(compile(`params = args()`))

	assert.Contains(t, output, `params=("$@")`)
}
