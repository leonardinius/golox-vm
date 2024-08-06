package vmdebug

import "github.com/leonardinius/goloxvm/internal/vm/vmchunk"

type Assert interface {
	Assertf(condition bool, message string, args ...any)
}

type Disassembler interface {
	DisassembleChunk(chunk *vmchunk.Chunk, name string)
	DisassembleInstruction(chunk *vmchunk.Chunk, offset int) int
}

var (
	gDD     Disassembler = nil
	gAssert Assert       = nil
)

func Assertf(condition bool, message string, args ...any) {
	if DebugAssert {
		gAssert.Assertf(condition, message, args...)
	}
}

func DisassembleChunk(chunk *vmchunk.Chunk, name string) {
	if DebugDisassembler {
		gDD.DisassembleChunk(chunk, name)
	}
}

func DisassembleInstruction(chunk *vmchunk.Chunk, offset int) {
	if DebugDisassembler {
		gDD.DisassembleInstruction(chunk, offset)
	}
}
