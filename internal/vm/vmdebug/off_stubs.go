//go:build !debug

package vmdebug

import (
	"github.com/leonardinius/goloxvm/internal/vm/vmchunk"
	"github.com/leonardinius/goloxvm/internal/vm/vmvalue"
)

const (
	DebugDisassembler = false
	DebugAssert       = false
)

func Printf(message string, args ...any) {}

func PrintValue(v vmvalue.Value) {}

func PrintObject[T vmvalue.VMObjectable](o *T) {}

func Assertf(condition bool, message string, args ...any) {}

func DisassembleChunk(chunk *vmchunk.Chunk, name string) {}

func DisassembleInstruction(chunk *vmchunk.Chunk, offset int) {}
