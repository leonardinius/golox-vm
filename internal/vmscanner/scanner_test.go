package vmscanner_test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonardinius/goloxvm/internal/tests"
	"github.com/leonardinius/goloxvm/internal/vmscanner"
)

func TestScanner(t *testing.T) {
	t.Parallel()
	testcases := tests.LoadFromDir(t, "testdata/scanner")
	for _, tc := range testcases {
		t.Run("=>/"+filepath.Base(tc.Testcase), func(t *testing.T) {
			scanner := vmscanner.NewScanner([]byte(tc.Input))
			outputs := []string{}
			for {
				token := scanner.ScanToken()
				if token.Type == vmscanner.TokenEOF {
					break
				}
				tokenAsText := fmt.Sprintf("%04d [%s] '%s'", token.Line, token.Type, token.Lexeme())
				outputs = append(outputs, tokenAsText)
			}
			actual := strings.Join(outputs, "\n")

			assert.Equalf(t, tc.Expected, actual, "Err '%s'", tc.Testcase)
		})
	}
}
