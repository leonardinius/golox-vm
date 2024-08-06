//go:build debug

package vmdebug

import (
	"fmt"

	"github.com/leonardinius/goloxvm/internal/vm/bytecode"
	"github.com/leonardinius/goloxvm/internal/vm/vmchunk"
)

const (
	DebugDisassembler = true
)

// Disassembler is an interface for disassembling chunks.
type stdoutDisassembler struct{}

var _ Disassembler = (*panicAssert)(nil)

func init() {
	gDD = &panicAssert{}
}

// DisassembleChunk implements disassemblerMethods.
func (s *panicAssert) DisassembleChunk(chunk *vmchunk.Chunk, name string) {
	fmt.Println()
	fmt.Println("== '" + name + "' byte code ==")

	for offset := 0; offset < chunk.Count; {
		offset = s.DisassembleInstruction(chunk, offset)
	}
}

func (s *panicAssert) DisassembleInstruction(chunk *vmchunk.Chunk, offset int) int {
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
		bytecode.OpSetGlobal,
		bytecode.OpDefineGlobal,
		bytecode.OpClosure:
		return s.constantInstruction(instruction, chunk, offset)
	case bytecode.OpGetLocal,
		bytecode.OpSetLocal,
		bytecode.OpCall:
		return s.byteInstruction(instruction, chunk, offset)
	case bytecode.OpJump,
		bytecode.OpJumpIfFalse:
		return s.jumpInstruction(instruction, 1, chunk, offset)
	case bytecode.OpLoop:
		return s.jumpInstruction(instruction, -1, chunk, offset)
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

func (s *panicAssert) constantInstruction(op bytecode.OpCode, chunk *vmchunk.Chunk, offset int) int {
	constant := chunk.Code[offset+1]
	fmt.Printf("%-16s %4d '", op, constant)
	PrintValue(chunk.ConstantAt(int(constant)))
	fmt.Println("'")
	return offset + 2
}

func (s *panicAssert) byteInstruction(op bytecode.OpCode, chunk *vmchunk.Chunk, offset int) int {
	slot := chunk.Code[offset+1]
	fmt.Printf("%-16s %4d\n", op, slot)
	return offset + 2
}

func (s *panicAssert) jumpInstruction(op bytecode.OpCode, sign int, chunk *vmchunk.Chunk, offset int) int {
	jump := int((uint16(chunk.Code[offset+1]) << 8) | uint16(chunk.Code[offset+2]))
	fmt.Printf("%-16s %4d -> %d\n", op, offset, offset+3+sign*jump)
	return offset + 3
}

func (s *panicAssert) simpleInstruction(op bytecode.OpCode, offset int) int {
	fmt.Println(op.String())
	return offset + 1
}
