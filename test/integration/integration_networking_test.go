package integration_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestE2E_JsonGet(t *testing.T) {
	// Write JSON to a temp file — Langz strings can't contain escaped quotes
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "data.json")
	err := os.WriteFile(jsonFile, []byte(`{"name": "Alice", "age": 30}`), 0644)
	assert.NoError(t, err)

	source := `
data = read("` + jsonFile + `")
name = json_get(data, ".name")
print(name)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "Alice", output)
}

func TestE2E_JsonGetNested(t *testing.T) {
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "data.json")
	err := os.WriteFile(jsonFile, []byte(`{"user": {"city": "Berlin"}}`), 0644)
	assert.NoError(t, err)

	source := `
data = read("` + jsonFile + `")
city = json_get(data, ".user.city")
print(city)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "Berlin", output)
}

func TestE2E_FetchConventionVars(t *testing.T) {
	// Test that fetch sets _status — use a bogus port so curl fails fast
	source := `
fetch("http://localhost:1", timeout: 1)
print(_status)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	// Script should complete (|| true prevents set -e from killing it)
	assert.Equal(t, 0, code)
	// _status should be set (even if 000 for connection refused)
	assert.NotEmpty(t, output)
}

func TestE2E_FetchOrFallback(t *testing.T) {
	source := `
data = fetch("http://localhost:1", timeout: 1) or "fallback_value"
print(data)
`
	bash := compileSource(t, source)
	output, code := runBash(t, bash)

	assert.Equal(t, 0, code)
	assert.Equal(t, "fallback_value", output)
}

func TestE2E_FetchGeneratesValidBash(t *testing.T) {
	// Single-line fetch with all kwargs — verify it generates multi-line curl
	source := `resp = fetch("https://example.com", method: "POST", body: "data", timeout: 5, retries: 2)`

	bash := compileSource(t, source)
	lines := strings.Split(bash, "\n")
	assert.True(t, len(lines) > 5, "expected multi-line output")
	assert.Contains(t, bash, "curl")
	assert.Contains(t, bash, "_status")
	assert.Contains(t, bash, "_fetch_attempt")
}
