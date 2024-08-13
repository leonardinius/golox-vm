//go:build debug

package vmdebug

import (
	"fmt"

	"github.com/leonardinius/goloxvm/internal/vm/bytecode"
	"github.com/leonardinius/goloxvm/internal/vm/vmchunk"
	"github.com/leonardinius/goloxvm/internal/vm/vmvalue"
)

const (
	DebugDisassembler = true
	DebugAssert       = true
)

func Assertf(condition bool, message string, args ...any) {
	if !condition {
		panic(fmt.Errorf(message, args...))
	}
}

// DisassembleChunk implements disassemblerMethods.
func DisassembleChunk(chunk *vmchunk.Chunk, name string) {
	fmt.Println("== '" + name + "' byte code ==")

	for offset := 0; offset < chunk.Count; {
		offset = DisassembleInstruction(chunk, offset)
	}
}

func DisassembleInstruction(chunk *vmchunk.Chunk, offset int) int {
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
		bytecode.OpDefineGlobal:
		return constantInstruction(instruction, chunk, offset)
	case bytecode.OpClosure:
		return closureInstruction(instruction, chunk, offset)
	case bytecode.OpGetLocal,
		bytecode.OpSetLocal,
		bytecode.OpGetUpvalue,
		bytecode.OpSetUpvalue,
		bytecode.OpCall:
		return byteInstruction(instruction, chunk, offset)
	case bytecode.OpJump,
		bytecode.OpJumpIfFalse:
		return jumpInstruction(instruction, 1, chunk, offset)
	case bytecode.OpLoop:
		return jumpInstruction(instruction, -1, chunk, offset)
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
		bytecode.OpCloseUpvalue,
		bytecode.OpReturn:
		return simpleInstruction(instruction, offset)
	default:
		panic(fmt.Sprintf("dd: unknown opcode (%d)\n", instruction))
	}
}

func constantInstruction(op bytecode.OpCode, chunk *vmchunk.Chunk, offset int) int {
	constant := chunk.Code[offset+1]
	fmt.Printf("%-16s %4d '", op, constant)
	PrintValue(chunk.ConstantAt(int(constant)))
	fmt.Println("'")
	return offset + 2
}

func closureInstruction(op bytecode.OpCode, chunk *vmchunk.Chunk, offset int) int {
	constant := chunk.Code[offset+1]
	value := chunk.ConstantAt(int(constant))
	fmt.Printf("%-16s %4d '", op, constant)
	PrintValue(value)
	fmt.Println("'")
	offset += 2

	fn := vmvalue.ValueAsFunction(value)
	for range fn.UpvalueCount {
		var tag string
		isLocal := chunk.Code[offset] == 1
		index := chunk.Code[offset+1]
		if isLocal {
			tag = "local"
		} else {
			tag = "upvalue"
		}
		fmt.Printf("%04d    | %-20s   %s %d\n", offset, "", tag, index)
		offset += 2
	}

	return offset
}

func byteInstruction(op bytecode.OpCode, chunk *vmchunk.Chunk, offset int) int {
	slot := chunk.Code[offset+1]
	fmt.Printf("%-16s %4d\n", op, slot)
	return offset + 2
}

func jumpInstruction(op bytecode.OpCode, sign int, chunk *vmchunk.Chunk, offset int) int {
	jump := int((uint16(chunk.Code[offset+1]) << 8) | uint16(chunk.Code[offset+2]))
	fmt.Printf("%-16s %4d -> %d\n", op, offset, offset+3+sign*jump)
	return offset + 3
}

func simpleInstruction(op bytecode.OpCode, offset int) int {
	fmt.Println(op.String())
	return offset + 1
}

func PrintValue(v vmvalue.Value) {
	vmvalue.PrintValue(v)
}

func PrintObject[T vmvalue.VMObjectable](o *T) {
	if o == nil {
		fmt.Print("nil")
		return
	}
	vmvalue.PrintAnyObject(o)
}
