//go:build debug

package vmdebug

import (
	"github.com/leonardinius/goloxvm/internal/vm/vmchunk"
)

const (
	DebugDisassembler = true
	DebugAssert       = true
)

var (
	gDD     disassembler = &stdoutDisassembler{}
	gAssert asserts      = &panicAssert{}
)

type asserts interface {
	Assertf(condition bool, message string, args ...any)
}

type disassembler interface {
	DisassembleChunk(chunk *vmchunk.Chunk, name string)
	DisassembleInstruction(chunk *vmchunk.Chunk, offset int) int
}

func Assertf(condition bool, message string, args ...any) {
	gAssert.Assertf(condition, message, args...)
}

func DisassembleChunk(chunk *vmchunk.Chunk, name string) {
	gDD.DisassembleChunk(chunk, name)
}

func DisassembleInstruction(chunk *vmchunk.Chunk, offset int) {
	gDD.DisassembleInstruction(chunk, offset)
}
