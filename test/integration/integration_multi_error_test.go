package integration_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2E_MultipleParseErrors(t *testing.T) {
	root := projectRoot(t)
	dir := t.TempDir()
	errFile := dir + "/errors.lz"

	// Multiple ILLEGAL tokens produce multiple parse errors
	require.NoError(t, os.WriteFile(errFile, []byte("x = @\ny = @\nz = @\n"), 0644))

	cmd := exec.Command("go", "run", root+"/cmd/langz", "run", errFile)
	out, err := cmd.CombinedOutput()
	require.Error(t, err, "should fail on parse errors")

	output := string(out)
	// Should show multiple error lines (one per @)
	errorCount := strings.Count(output, errFile+":")
	assert.GreaterOrEqual(t, errorCount, 2, "should report multiple errors, got:\n%s", output)
}

func TestE2E_ErrorCapAt10(t *testing.T) {
	root := projectRoot(t)
	dir := t.TempDir()
	errFile := dir + "/many_errors.lz"

	// 15 ILLEGAL tokens — should cap at 10 shown
	var source string
	for i := 0; i < 15; i++ {
		source += "@ \n"
	}
	require.NoError(t, os.WriteFile(errFile, []byte(source), 0644))

	cmd := exec.Command("go", "run", root+"/cmd/langz", "run", errFile)
	out, err := cmd.CombinedOutput()
	require.Error(t, err)
	assert.Contains(t, string(out), "... and")
	assert.Contains(t, string(out), "more error")
}
