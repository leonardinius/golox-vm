package vmdebug

import (
	"github.com/leonardinius/goloxvm/internal/vmchunk"
)

type Disassembler interface {
	DisassembleChunk(chunk *vmchunk.Chunk, name string)
	DisassembleInstruction(chunk *vmchunk.Chunk, offset int) int
}

var dd Disassembler = nil

func DisassembleChunk(chunk *vmchunk.Chunk, name string) {
	if DebugDisassembler {
		dd.DisassembleChunk(chunk, name)
	}
}

func DisassembleInstruction(chunk *vmchunk.Chunk, offset int) int {
	if DebugDisassembler {
		return dd.DisassembleInstruction(chunk, offset)
	}
	return 0
}
