//go:build debug

package vmdebug

import (
	"fmt"

	"github.com/leonardinius/goloxvm/internal/vmchunk"
)

const DebugDisassembler = true

// Disassembler is an interface for disassembling chunks.
type stdoutDisassembler struct{}

var _ Disassembler = (*stdoutDisassembler)(nil)

func init() {
	dd = &stdoutDisassembler{}
}

// DisassembleChunk implements disassemblerMethods.
func (s *stdoutDisassembler) DisassembleChunk(chunk *vmchunk.Chunk, name string) {
	fmt.Println()
	fmt.Println("== '" + name + "' byte code ==")

	for offset := 0; offset < chunk.Count; {
		offset = s.DisassembleInstruction(chunk, offset)
	}
}

func (s *stdoutDisassembler) DisassembleInstruction(chunk *vmchunk.Chunk, offset int) int {
	fmt.Printf("%04d ", offset)
	line := chunk.DebugGetLine(offset)
	if offset > 0 && line == chunk.DebugGetLine(offset-1) {
		fmt.Print("   | ")
	} else {
		fmt.Printf("%4d ", line)
	}

	instruction := vmchunk.OpCode(chunk.Code[offset])
	switch instruction {

	case vmchunk.OpConstant:
		return s.constantInstruction("OP_CONSTANT", chunk, offset)

	case vmchunk.OpAdd:
		return s.simpleInstruction("OP_ADD", offset)

	case vmchunk.OpSubtract:
		return s.simpleInstruction("OP_SUBTRACT", offset)

	case vmchunk.OpMultiply:
		return s.simpleInstruction("OP_MULTIPLY", offset)

	case vmchunk.OpDivide:
		return s.simpleInstruction("OP_DIVIDE", offset)

	case vmchunk.OpNegate:
		return s.simpleInstruction("OP_NEGATE", offset)

	case vmchunk.OpPop:
		return s.simpleInstruction("OP_POP", offset)

	case vmchunk.OpReturn:
		return s.simpleInstruction("OP_RETURN", offset)

	default:
		fmt.Printf("Unknown opcode %d\n", instruction)
		return offset + 1
	}
}

func (s *stdoutDisassembler) constantInstruction(name string, chunk *vmchunk.Chunk, offset int) int {
	constant := chunk.Code[offset+1]
	fmt.Printf("%-16s %4d '", name, constant)
	PrintValue(chunk.Constants.At(int(constant)))
	fmt.Println("'")
	return offset + 2
}

func (s *stdoutDisassembler) simpleInstruction(name string, offset int) int {
	fmt.Println(name)
	return offset + 1
}
