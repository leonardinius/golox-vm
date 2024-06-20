//go:build debug

package vmdebug

import (
	"fmt"

	"github.com/leonardinius/goloxvm/internal/vmchunk"
	"github.com/leonardinius/goloxvm/internal/vmvalue"
)

const DEBUG = true

// Disassembler is an interface for disassembling chunks.
type stdoutDisassembler struct{}

var _ Disassembler = (*stdoutDisassembler)(nil)

func init() {
	dd = &stdoutDisassembler{}
}

// DisassembleChunk implements disassemblerMethods.
func (s *stdoutDisassembler) DisassembleChunk(chunk *vmchunk.Chunk, name string) {
	println("== " + name + " ==")

	for offset := 0; offset < chunk.Count; {
		offset = s.disassembleInstruction(chunk, offset)
	}
}

func (s *stdoutDisassembler) disassembleInstruction(chunk *vmchunk.Chunk, offset int) int {
	fmt.Printf("%04d ", offset)

	instruction := vmchunk.OpCode(chunk.Code[offset])
	switch instruction {
	case vmchunk.OpConstant:
		return s.constantInstruction("OP_CONSTANT", chunk, offset)
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
	s.printValue(chunk.Constants.At(int(constant)))
	fmt.Println("'")
	return offset + 2
}

func (s *stdoutDisassembler) simpleInstruction(name string, offset int) int {
	fmt.Println(name)
	return offset + 1
}

func (s *stdoutDisassembler) printValue(v vmvalue.Value) {
	fmt.Printf("%g", v)
}
