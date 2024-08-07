package tests

import (
	"bufio"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type Testcase struct {
	Testcase string
	Input    string
	Expected string
}

var (
	errNoProjectDir = errors.New("unable to locate project directory")
	errNoTestCases  = errors.New("no test cases")
)

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func projectDir() (string, error) {
	d, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// walk back up to find the go.mod file
	for {
		if fileExists(filepath.Join(d, "go.mod")) {
			return d, nil
		}
		if d == "" || d == "." {
			return d, errNoProjectDir
		}
		d = filepath.Dir(d)
	}
}

func loadTestCases(t *testing.T, wd, dir, fileSuffix string) []Testcase {
	t.Helper()
	testcases := []Testcase{}
	err := filepath.WalkDir(filepath.Join(wd, dir), func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		if fileSuffix != "" && !strings.HasSuffix(path, fileSuffix) {
			t.Logf("skipping %s", path)
			return nil
		}

		f, err := os.Open(path) //nolint:gosec
		if err != nil {
			return err
		}

		inputs := []string{}
		expectactions := []string{}
		current := &inputs

		s := bufio.NewScanner(f)
		for s.Scan() {
			line := s.Text()
			// testcase directive

			if strings.HasPrefix(line, "//!#") {
				if strings.HasPrefix(line, "//!# Expect") {
					current = &expectactions
				}
				continue
			}
			*current = append(*current, line)
		}

		for len(inputs) > 0 && strings.TrimSpace(inputs[len(inputs)-1]) == "" {
			inputs = inputs[:len(inputs)-1]
		}
		for len(expectactions) > 0 && strings.TrimSpace(expectactions[len(expectactions)-1]) == "" {
			expectactions = expectactions[:len(expectactions)-1]
		}

		normalizedPath, err := filepath.Rel(wd, path)
		require.NoError(t, err)
		testcase := Testcase{
			Testcase: normalizedPath,
			Input:    strings.Join(inputs, "\n"),
			Expected: strings.Join(expectactions, "\n"),
		}
		testcases = append(testcases, testcase)

		return nil
	})
	require.NoError(t, err, "failed to load testcases")
	if len(testcases) == 0 {
		require.NoErrorf(t, errNoTestCases, "no test cases found in %s", dir)
	}
	return testcases
}

func LoadFromDir(t *testing.T, dir string) []Testcase {
	t.Helper()
	wd, err := projectDir()
	require.NoError(t, err)
	return loadTestCases(t, wd, dir, ".testcase")
}
