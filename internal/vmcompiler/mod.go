package vmcompiler

import (
	"github.com/leonardinius/goloxvm/internal/vm/vmchunk"
	"github.com/leonardinius/goloxvm/internal/vmcompiler/scanner"
)

var (
	gScanner        scanner.Scanner
	gParser         Parser
	gCurrent        *Compiler
	gCompilingChunk *vmchunk.Chunk
)
