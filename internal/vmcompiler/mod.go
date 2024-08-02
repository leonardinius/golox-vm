package vmcompiler

import (
	"github.com/leonardinius/goloxvm/internal/vmcompiler/scanner"
)

var (
	gScanner scanner.Scanner
	gParser  Parser
	gCurrent *Compiler = nil
)
