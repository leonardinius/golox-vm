//go:build debug

package debug

import (
	"fmt"

	"github.com/leonardinius/goloxvm/internal/chunk"
)

type stdoutDisassembler struct{}

var _ disassembler = (*stdoutDisassembler)(nil)

func init() {
	dd = &stdoutDisassembler{}
}

// disassembleChunk implements disassemblerMethods.
func (s *stdoutDisassembler) disassembleChunk(ch *chunk.Chunk, name string) {
	println("== " + name + " ==")

	for offset := 0; offset < ch.Count; {
		offset = s.disassembleInstruction(ch, offset)
	}
}

// disassembleInstruction implements disassemblerMethods.
func (s *stdoutDisassembler) disassembleInstruction(_ *chunk.Chunk, offset int) int {
	fmt.Printf("%d\n", offset)
	return offset + 1
}
