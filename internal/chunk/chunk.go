package chunk

import (
	"github.com/leonardinius/goloxvm/internal/mem"
)

type OpCode byte

const (
	OpNil OpCode = iota
	OpReturn
)

type Chunk struct {
	Code  []OpCode
	Count int
}

func NewChunk() Chunk {
	return Chunk{}
}

func (chunk *Chunk) InitChunk() {
	chunk.Code = nil
	chunk.Count = 0
}

func (chunk *Chunk) FreeChunk() {
	chunk.Code = mem.FreeArray(chunk.Code)
	chunk.InitChunk()
}

func (chunk *Chunk) Write(op OpCode) {
	chunk.Code = append(chunk.Code, op)
	if cap(chunk.Code) < chunk.Count+1 {
		capacity := mem.GrowCapacity(cap(chunk.Code))
		chunk.Code = mem.GrowArray(chunk.Code, capacity)
	}
	chunk.Code[chunk.Count] = op
	chunk.Count++
}
