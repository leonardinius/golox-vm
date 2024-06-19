//go:build !debug

package disassembler

import "github.com/leonardinius/goloxvm/internal/chunk"

type noOpDisassembler struct{}

var _ disassemblerMethods = (*noOpDisassembler)(nil)

func init() {
	dd = &noOpDisassembler{}
}

// disassembleChunk implements disassemblerMethods.
func (s *noOpDisassembler) disassembleChunk(ch *chunk.Chunk, name string) {
}
