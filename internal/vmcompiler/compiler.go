package vmcompiler

import (
	"math"

	"github.com/leonardinius/goloxvm/internal/vmchunk"
	"github.com/leonardinius/goloxvm/internal/vmscanner"
	"github.com/leonardinius/goloxvm/internal/vmvalue"
)

var (
	scanner        vmscanner.Scanner
	parser         Parser
	compilingChunk *vmchunk.Chunk
)

func Compile(source []byte, chunk *vmchunk.Chunk) bool {
	scanner = vmscanner.NewScanner(source)
	defer scanner.Free()
	compilingChunk = chunk
	defer endCompiler()

	parser = NewParser()

	advance()
	expression()
	consume(vmscanner.TokenEOF, "Expect end of expression.")

	return !parser.hadError
}

func currentChunk() *vmchunk.Chunk {
	return compilingChunk
}

func emitCode1(op vmchunk.OpCode) {
	currentChunk().WriteOpcode(op, parser.previous.Line)
}

func emitCode2(op vmchunk.OpCode, b byte) {
	currentChunk().WriteOpcode(op, parser.previous.Line)
	currentChunk().Write(b, parser.previous.Line)
}

func emitConstant(v vmvalue.Value) {
	emitCode2(vmchunk.OpConstant, makeConstant(v))
}

func makeConstant(v vmvalue.Value) byte {
	constant := currentChunk().AddConstant(v)
	if constant > math.MaxUint8 {
		errorAtPrev("Too many constants in one chunk.")
		return 0
	}
	return byte(constant)
}

func emitReturn() {
	emitCode1(vmchunk.OpReturn)
}

func endCompiler() {
	emitReturn()
}
