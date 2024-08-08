package smoketest_test

import (
	"io"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/leonardinius/goloxvm/internal/cmd"
	"github.com/leonardinius/goloxvm/internal/tests"
)

func TestFibonacci(t *testing.T) {
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	t.Cleanup(func() {
		os.Stdout = orig
	})

	// Get the project directory
	dir, err := tests.ProjectDir()
	require.NoError(t, err)
	// Run the scipt
	script := path.Join(dir, "testdata", "fib.lox")
	rc := cmd.Main(script)

	_ = w.Close()
	out, _ := io.ReadAll(r)
	assert.Equal(t, 0, rc)
	// Parse the output
	lines := strings.SplitN(string(out), "\n", 2)
	// Check the output
	assert.Equal(t, "9227465.0", strings.TrimSpace(lines[0]))
}
