package debug

import "github.com/leonardinius/goloxvm/internal/chunk"

type disassembler interface {
	disassembleChunk(ch *chunk.Chunk, name string)
}

var dd disassembler

func DisassembleChunk(ch *chunk.Chunk, name string) {
	dd.disassembleChunk(ch, name)
}
