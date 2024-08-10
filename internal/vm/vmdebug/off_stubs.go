//go:build !debug

package vmdebug

import (
	"github.com/leonardinius/goloxvm/internal/vm/vmchunk"
)

const (
	DebugDisassembler = false
	DebugAssert       = false
)

func Assertf(condition bool, message string, args ...any) {}

func DisassembleChunk(chunk *vmchunk.Chunk, name string) {}

func DisassembleInstruction(chunk *vmchunk.Chunk, offset int) {}
""