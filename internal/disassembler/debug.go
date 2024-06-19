package disassembler

import "github.com/leonardinius/goloxvm/internal/chunk"

type disassemblerMethods interface {
	disassembleChunk(ch *chunk.Chunk, name string)
}

var dd disassemblerMethods

func DisassembleChunk(ch *chunk.Chunk, name string) {
	dd.disassembleChunk(ch, name)
}
