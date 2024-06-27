package vmdebug

import (
	"github.com/leonardinius/goloxvm/internal/vmchunk"
)

type Disassembler interface {
	DisassembleChunk(chunk *vmchunk.Chunk, name string)
	DissasembleInstruction(chunk *vmchunk.Chunk, offset int) int
}

var dd Disassembler = nil

func DisassembleChunk(chunk *vmchunk.Chunk, name string) {
	if DEBUG_DISASSEMBLER {
		dd.DisassembleChunk(chunk, name)
	}
}

func DissasembleInstruction(chunk *vmchunk.Chunk, offset int) int {
	if DEBUG_DISASSEMBLER {
		return dd.DissasembleInstruction(chunk, offset)
	}
	return 0
}
