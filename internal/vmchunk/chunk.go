package vmchunk

import (
	"github.com/leonardinius/goloxvm/internal/vmmem"
	"github.com/leonardinius/goloxvm/internal/vmvalue"
)

type OpCode byte

const (
	OpConstant OpCode = iota
	OpNegate
	OpReturn
)

type Chunk struct {
	Code      []uint8
	Count     int
	Constants vmvalue.ValueArray
	Lines     Lines
}

func NewChunk() Chunk {
	return Chunk{}
}

func (chunk *Chunk) InitChunk() {
	chunk.Code = nil
	chunk.Count = 0
	vmvalue.InitValueArray(&chunk.Constants)
	chunk.Lines.Init()
}

func (chunk *Chunk) FreeChunk() {
	chunk.Code = vmmem.FreeArray(chunk.Code)
	vmvalue.FreeValueArray(&chunk.Constants)
	chunk.Lines.Free()
	chunk.InitChunk()
}

func (chunk *Chunk) WriteOpcode(op OpCode, line int) {
	chunk.Write(byte(op), line)
}

func (chunk *Chunk) Write(op byte, line int) {
	chunk.Code = append(chunk.Code, op)
	if cap(chunk.Code) < chunk.Count+1 {
		capacity := vmmem.GrowCapacity(cap(chunk.Code))
		chunk.Code = vmmem.GrowArray(chunk.Code, capacity)
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
