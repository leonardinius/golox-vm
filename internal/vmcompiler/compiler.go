package vmcompiler

import (
	"math"

	"github.com/leonardinius/goloxvm/internal/bytecode"
	"github.com/leonardinius/goloxvm/internal/scanner"
	"github.com/leonardinius/goloxvm/internal/tokens"
	"github.com/leonardinius/goloxvm/internal/vmchunk"
	"github.com/leonardinius/goloxvm/internal/vmvalue"
)

var (
	gScanner        scanner.Scanner
	gParser         Parser
	gCompilingChunk *vmchunk.Chunk
)

func Compile(source []byte, chunk *vmchunk.Chunk) bool {
	gScanner = scanner.NewScanner(source)
	defer gScanner.Free()
	gCompilingChunk = chunk
	defer endCompiler()

	gParser = NewParser()

	advance()
	expression()
	consume(tokens.TokenEOF, "Expect end of expression.")

	return !gParser.hadError
}

func currentChunk() *vmchunk.Chunk {
	return gCompilingChunk
}

func emitOpcode(op bytecode.OpCode) {
	currentChunk().WriteOpcode(op, gParser.previous.Line)
}

func emitOpcodes(op1 bytecode.OpCode, op2 bytecode.OpCode) {
	emitOpcode(op1)
	emitOpcode(op2)
}

func emitBytes(op bytecode.OpCode, b byte) {
	currentChunk().WriteOpcode(op, gParser.previous.Line)
	currentChunk().Write(b, gParser.previous.Line)
}

func emitConstant(v vmvalue.Value) {
	emitBytes(bytecode.OpConstant, makeConstant(v))
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
	emitOpcode(bytecode.OpReturn)
}

func endCompiler() {
	emitReturn()
}
