package tests

import (
	"errors"
	"os"
	"path/filepath"
)

var errNoProjectDir = errors.New("unable to locate project directory")

func ProjectDir() (string, error) {
	return projectDir()
}

func fileExists(filePath string) bool {
	info, err := os.Lstat(filePath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	if info.IsDir() {
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
