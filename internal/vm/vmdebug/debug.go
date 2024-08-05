package vmdebug

import "github.com/leonardinius/goloxvm/internal/vm/vmchunk"

type Disassembler interface {
	DisassembleChunk(chunk *vmchunk.Chunk, name string)
	DisassembleInstruction(chunk *vmchunk.Chunk, offset int) int
}

var gDD Disassembler = nil

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
