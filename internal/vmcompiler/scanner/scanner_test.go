package scanner_test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonardinius/goloxvm/internal/tests"
	"github.com/leonardinius/goloxvm/internal/vmcompiler/scanner"
	"github.com/leonardinius/goloxvm/internal/vmcompiler/tokens"
)

func TestScanner(t *testing.T) {
	t.Parallel()
	testcases := tests.LoadFromDir(t, "testdata/scanner")
	for _, tc := range testcases {
		t.Run("=>/"+filepath.Base(tc.Testcase), func(t *testing.T) {
			s := scanner.NewScanner([]byte(tc.Input))
			outputs := []string{}
			for {
				token := s.ScanToken()
				if token.Type == tokens.TokenEOF {
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
