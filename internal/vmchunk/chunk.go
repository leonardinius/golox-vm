package vmchunk

import (
	"github.com/leonardinius/goloxvm/internal/vmmem"
	"github.com/leonardinius/goloxvm/internal/vmvalue"
)

type OpCode byte

const (
	OpConstant OpCode = iota
	OpReturn
)

type Chunk struct {
	Code      []uint8
	Count     int
	Constants vmvalue.ValueArray
}

func NewChunk() Chunk {
	return Chunk{}
}

func (chunk *Chunk) InitChunk() {
	chunk.Code = nil
	chunk.Count = 0
	vmvalue.InitValueArray(&chunk.Constants)
}

func (chunk *Chunk) FreeChunk() {
	chunk.Code = vmmem.FreeArray(chunk.Code)
	vmvalue.FreeValueArray(&chunk.Constants)
	chunk.InitChunk()
}

func (chunk *Chunk) Write(op OpCode) {
	chunk.Write1(byte(op))
}

func (chunk *Chunk) Write1(op byte) {
	chunk.Code = append(chunk.Code, op)
	if cap(chunk.Code) < chunk.Count+1 {
		capacity := vmmem.GrowCapacity(cap(chunk.Code))
		chunk.Code = vmmem.GrowArray(chunk.Code, capacity)
	}
	chunk.Code[chunk.Count] = op
	chunk.Count++
}

func (chunk *Chunk) AddConstant(v vmvalue.Value) int {
	return chunk.Constants.Write(v)
}
