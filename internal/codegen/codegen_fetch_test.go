package codegen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchSimpleGET(t *testing.T) {
	output := body(compile(`data = fetch("https://api.example.com/health")`))

	assert.Contains(t, output, `_tmp_headers=$(mktemp)`)
	assert.Contains(t, output, `_tmp_body=$(mktemp)`)
	assert.Contains(t, output, `curl -s -w "%{http_code}"`)
	assert.Contains(t, output, `-D "$_tmp_headers"`)
	assert.Contains(t, output, `-o "$_tmp_body"`)
	assert.Contains(t, output, `"https://api.example.com/health"`)
	assert.Contains(t, output, `_body=$(cat "$_tmp_body")`)
	assert.Contains(t, output, `_headers=$(cat "$_tmp_headers")`)
	assert.Contains(t, output, `rm -f "$_tmp_headers" "$_tmp_body"`)
	assert.Contains(t, output, `data="$_body"`)
	// Simple GET should NOT have -X flag
	assert.NotContains(t, output, `-X`)
}

func TestFetchPOSTWithBody(t *testing.T) {
	output := body(compile(`resp = fetch("https://api.com/data", method: "POST", body: "payload")`))

	assert.Contains(t, output, `-X POST`)
	assert.Contains(t, output, `-d "payload"`)
	assert.Contains(t, output, `resp="$_body"`)
}

func TestFetchWithHeaders(t *testing.T) {
	output := body(compile(`resp = fetch("https://api.com", headers: {content_type: "application/json"})`))

	assert.Contains(t, output, `-H "content_type: application/json"`)
}

func TestFetchWithTimeout(t *testing.T) {
	output := body(compile(`resp = fetch("https://api.com", timeout: 30)`))

	assert.Contains(t, output, `--max-time 30`)
}

func TestFetchStandalone(t *testing.T) {
	output := body(compile(`fetch("https://api.com/webhook", method: "POST", body: "event")`))

	assert.Contains(t, output, `curl -s`)
	assert.Contains(t, output, `-X POST`)
	assert.Contains(t, output, `_status=`)
	assert.Contains(t, output, `_body=`)
	// No variable assignment at the end
	assert.NotContains(t, output, `="$_body"`)
}

func TestFetchWithVariableBody(t *testing.T) {
	output := body(compile(`resp = fetch("https://api.com", method: "PUT", body: payload)`))

	assert.Contains(t, output, `-X PUT`)
	assert.Contains(t, output, `-d "$payload"`)
}

func TestFetchSetsConventionVars(t *testing.T) {
	output := body(compile(`data = fetch("https://api.com")`))

	assert.Contains(t, output, `_status=`)
	assert.Contains(t, output, `_body=`)
	assert.Contains(t, output, `_headers=`)
}

func TestFetchOrTrueForSetE(t *testing.T) {
	output := body(compile(`data = fetch("https://api.com")`))

	// curl must end with || true to prevent set -e from killing the script
	assert.Contains(t, output, `|| true`)
}

func TestFetchWithRetries(t *testing.T) {
	output := body(compile(`data = fetch("https://api.com", retries: 3)`))

	assert.Contains(t, output, `_fetch_attempt=0`)
	assert.Contains(t, output, `_fetch_max=3`)
	assert.Contains(t, output, `while [ "$_fetch_attempt" -lt "$_fetch_max" ]; do`)
	assert.Contains(t, output, `_fetch_attempt=$((_fetch_attempt + 1))`)
	assert.Contains(t, output, `break`)
	assert.Contains(t, output, `sleep 1`)
	assert.Contains(t, output, `done`)
	assert.Contains(t, output, `data="$_body"`)
}

func TestFetchOrFallback(t *testing.T) {
	output := body(compile(`data = fetch("https://api.com") or "cached_data"`))

	assert.Contains(t, output, `curl -s`)
	assert.Contains(t, output, `_status=`)
	assert.Contains(t, output, `if [ "$_status" -ge 200 ] && [ "$_status" -lt 300 ]; then`)
	assert.Contains(t, output, `data="cached_data"`)
}

func TestFetchOrExit(t *testing.T) {
	output := body(compile(`data = fetch("https://api.com") or exit(1)`))

	assert.Contains(t, output, `curl -s`)
	assert.Contains(t, output, `exit 1`)
}
