//go:build debug

package vmdebug

import (
	"fmt"

	"github.com/leonardinius/goloxvm/internal/vm/bytecode"
	"github.com/leonardinius/goloxvm/internal/vm/vmchunk"
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
	case bytecode.OpConstant,
		bytecode.OpGetGlobal,
		bytecode.OpDefineGlobal:
		return s.constantInstruction(instruction, chunk, offset)
	case bytecode.OpNil,
		bytecode.OpTrue,
		bytecode.OpFalse,
		bytecode.OpEqual,
		bytecode.OpGreater,
		bytecode.OpLess,
		bytecode.OpAdd,
		bytecode.OpSubtract,
		bytecode.OpMultiply,
		bytecode.OpDivide,
		bytecode.OpNot,
		bytecode.OpNegate,
		bytecode.OpPop,
		bytecode.OpPrint,
		bytecode.OpReturn:
		return s.simpleInstruction(instruction, offset)
	default:
		fmt.Printf("dd: unknown opcode %s (%d)\n", instruction, instruction)
		return offset + 1
	}
}

func (s *stdoutDisassembler) constantInstruction(op bytecode.OpCode, chunk *vmchunk.Chunk, offset int) int {
	constant := chunk.Code[offset+1]
	fmt.Printf("%-16s %4d '", op, constant)
	DebugValue(chunk.Constants.At(int(constant)))
	fmt.Println("'")
	return offset + 2
}

func (s *stdoutDisassembler) simpleInstruction(op bytecode.OpCode, offset int) int {
	fmt.Println(op.String())
	return offset + 1
}
