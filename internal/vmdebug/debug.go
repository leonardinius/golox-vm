package vmdebug

import (
	"github.com/leonardinius/goloxvm/internal/vmchunk"
)

type Disassembler interface {
	DisassembleChunk(chunk *vmchunk.Chunk, name string)
}

var dd Disassembler = nil

func DisassembleChunk(chunk *vmchunk.Chunk, name string) {
	if DEBUG {
		dd.DisassembleChunk(chunk, name)
	}
}
