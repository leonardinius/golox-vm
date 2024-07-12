package vmcompiler

import (
	"github.com/leonardinius/goloxvm/internal/vmchunk"
	"github.com/leonardinius/goloxvm/internal/vmvalue"
)

func Compile(source string) (vmchunk.Chunk, error) {
	chunk := vmchunk.NewChunk()
	chunk.InitChunk()

	constant1 := chunk.AddConstant(1.1)
	chunk.WriteOpcode(vmchunk.OpConstant, 1)
	chunk.Write(byte(constant1), 1)
	constant2 := chunk.AddConstant(1.2)
	chunk.WriteOpcode(vmchunk.OpConstant, 1)
	chunk.Write(byte(constant2), 1)
	chunk.WriteOpcode(vmchunk.OpNegate, 1)
	chunk.WriteOpcode(vmchunk.OpPop, 1)
	chunk.WriteOpcode(vmchunk.OpPop, 1)

	addBinOp(&chunk, vmchunk.OpAdd, 2.0, 3.0)
	chunk.WriteOpcode(vmchunk.OpPop, 1)
	addBinOp(&chunk, vmchunk.OpSubtract, 2.0, 3.0)
	chunk.WriteOpcode(vmchunk.OpPop, 1)
	addBinOp(&chunk, vmchunk.OpMultiply, 2.0, 3.0)
	chunk.WriteOpcode(vmchunk.OpPop, 1)
	addBinOp(&chunk, vmchunk.OpDivide, 2.0, 3.0)
	chunk.WriteOpcode(vmchunk.OpReturn, 1)

	return chunk, nil
}

func addBinOp(chunk *vmchunk.Chunk, op vmchunk.OpCode, a, b vmvalue.Value) {
	aConstant := chunk.AddConstant(a)
	bConstant := chunk.AddConstant(b)

	chunk.WriteOpcode(vmchunk.OpConstant, 1)
	chunk.Write(byte(aConstant), 1)
	chunk.WriteOpcode(vmchunk.OpConstant, 1)
	chunk.Write(byte(bConstant), 1)
	chunk.WriteOpcode(op, 1)
}
