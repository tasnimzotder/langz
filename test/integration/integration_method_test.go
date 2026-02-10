package integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestE2E_Replace(t *testing.T) {
	source := `
name = "hello world"
result = name.replace("world", "langz")
print(result)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "hello langz", output)
}

func TestE2E_Contains(t *testing.T) {
	source := `
path = "/usr/local/bin"
if path.contains("local") {
  print("found")
}
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "found", output)
}

func TestE2E_StartsWith(t *testing.T) {
	source := `
url = "https://example.com"
if url.starts_with("https") {
  print("secure")
}
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "secure", output)
}

func TestE2E_EndsWith(t *testing.T) {
	source := `
file = "deploy.sh"
if file.ends_with(".sh") {
  print("shell")
}
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "shell", output)
}
