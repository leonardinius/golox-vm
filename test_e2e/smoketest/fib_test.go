package smoketest_test

import (
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/leonardinius/goloxvm/internal/cmd"
	"github.com/leonardinius/goloxvm/internal/tests"
)

func TestFibonacci(t *testing.T) {
	dir, err := tests.ProjectDir()
	require.NoError(t, err)

	script := path.Join(dir, "testdata", "fib.lox")
	rc := cmd.Main(script)
	require.Equal(t, 0, rc)
}
