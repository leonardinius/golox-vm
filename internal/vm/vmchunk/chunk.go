package vmchunk

import (
	"github.com/leonardinius/goloxvm/internal/vm/bytecode"
	"github.com/leonardinius/goloxvm/internal/vm/vmmem"
	"github.com/leonardinius/goloxvm/internal/vm/vmvalue"
)

type Chunk struct {
	Code      []uint8
	Count     int
	Constants vmvalue.ValueArray
	Lines     Lines
}

func NewChunk() Chunk {
	chunk := Chunk{}
	chunk.Constants = vmvalue.NewValueArray()
	chunk.resetChunk()
	return chunk
}

func FromPtr(ptr any) *Chunk {
	return ptr.(*Chunk)
}

func (chunk *Chunk) resetChunk() {
	chunk.Code = nil
	chunk.Count = 0
	chunk.Constants.Init()
	chunk.Lines.Init()
}

func (chunk *Chunk) Free() {
	chunk.Code = vmmem.FreeSlice(chunk.Code)
	chunk.Constants.Free()
	chunk.Lines.Free()
	chunk.resetChunk()
}

func (chunk *Chunk) AsPtr() any {
	return any(chunk)
}

func (chunk *Chunk) WriteOpcode(op bytecode.OpCode, line int) {
	chunk.Write(byte(op), line)
}

func (chunk *Chunk) Write(op byte, line int) {
	if len(chunk.Code) < chunk.Count+1 {
		capacity := vmmem.GrowCapacity(cap(chunk.Code))
		chunk.Code = vmmem.GrowSlice(chunk.Code, capacity)
	}
	chunk.Code[chunk.Count] = op
	chunk.Lines.MustWriteOffset(chunk.Count, line)
	chunk.Count++
}

func (chunk *Chunk) DebugGetLine(offset int) int {
	return chunk.Lines.GetLineByOffset(offset)
}

func (chunk *Chunk) AddConstant(v vmvalue.Value) int {
	return chunk.Constants.Write(v)
}

func (chunk *Chunk) ConstantAt(at int) vmvalue.Value {
	return chunk.Constants.At(at)
}
