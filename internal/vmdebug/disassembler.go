//go:build debug

package vmdebug

import (
	"fmt"

	"github.com/leonardinius/goloxvm/internal/bytecode"
	"github.com/leonardinius/goloxvm/internal/vmchunk"
)

const DebugDisassembler = true

// Disassembler is an interface for disassembling chunks.
type stdoutDisassembler struct{}

var _ Disassembler = (*stdoutDisassembler)(nil)

func init() {
	gDD = &stdoutDisassembler{}
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

	instruction := bytecode.OpCode(chunk.Code[offset])
	switch instruction {

	case bytecode.OpConstant:
		return s.constantInstruction(instruction, chunk, offset)

	case bytecode.OpAdd:
		return s.simpleInstruction(instruction, offset)

	case bytecode.OpSubtract:
		return s.simpleInstruction(instruction, offset)

	case bytecode.OpMultiply:
		return s.simpleInstruction(instruction, offset)

	case bytecode.OpDivide:
		return s.simpleInstruction(instruction, offset)

	case bytecode.OpNegate:
		return s.simpleInstruction(instruction, offset)

	case bytecode.OpPop:
		return s.simpleInstruction(instruction, offset)

	case bytecode.OpReturn:
		return s.simpleInstruction(instruction, offset)

	default:
		fmt.Printf("Unknown opcode %s (%d)\n", instruction, instruction)
		return offset + 1
	}
}

func (s *stdoutDisassembler) constantInstruction(op bytecode.OpCode, chunk *vmchunk.Chunk, offset int) int {
	constant := chunk.Code[offset+1]
	fmt.Printf("%-16s %4d '", op, constant)
	PrintValue(chunk.Constants.At(int(constant)))
	fmt.Println("'")
	return offset + 2
}

func (s *stdoutDisassembler) simpleInstruction(op bytecode.OpCode, offset int) int {
	fmt.Println(op.String())
	return offset + 1
}
