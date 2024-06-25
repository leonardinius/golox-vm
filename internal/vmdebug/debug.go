package vmdebug

import (
	"github.com/leonardinius/goloxvm/internal/vmchunk"
	"github.com/leonardinius/goloxvm/internal/vmvalue"
)

type Disassembler interface {
	DisassembleChunk(chunk *vmchunk.Chunk, name string)
	DissasembleInstruction(chunk *vmchunk.Chunk, offset int) int
	PrintValue(v vmvalue.Value)
}

var dd Disassembler = nil

func DisassembleChunk(chunk *vmchunk.Chunk, name string) {
	if DEBUG {
		dd.DisassembleChunk(chunk, name)
	}
}

func DissasembleInstruction(chunk *vmchunk.Chunk, offset int) int {
	if DEBUG {
		return dd.DissasembleInstruction(chunk, offset)
	}
	return 0
}

func PrintValue(v vmvalue.Value) {
	if DEBUG {
		dd.PrintValue(v)
	}
}
