package vmcompiler

import (
	"math"

	"github.com/leonardinius/goloxvm/internal/vm/bytecode"
	"github.com/leonardinius/goloxvm/internal/vm/vmchunk"
	"github.com/leonardinius/goloxvm/internal/vm/vmvalue"
	"github.com/leonardinius/goloxvm/internal/vmcompiler/scanner"
	"github.com/leonardinius/goloxvm/internal/vmcompiler/tokens"
)

const (
	MaxConstantCount = math.MaxUint8
	MaxLocalCount    = math.MaxUint8 + 1
)

type Compiler struct {
	Locals     [MaxLocalCount]Local
	LocalCount int
	ScoreDepth int
}

type Local struct {
	Name  scanner.Token
	Depth int
}

func NewCompiler() Compiler {
	c := Compiler{}
	gCurrent = &c
	return c
}

func Compile(source []byte, chunk *vmchunk.Chunk) bool {
	gScanner = scanner.NewScanner(source)
	defer gScanner.Free()
	_ = NewCompiler()
	gCompilingChunk = chunk
	defer endCompiler()

	gParser = NewParser()

	advance()

	for !match(tokens.TokenEOF) {
		declaration()
	}

	return !gParser.hadError
}

func currentChunk() *vmchunk.Chunk {
	return gCompilingChunk
}

func emitOpcode(op bytecode.OpCode) {
	currentChunk().WriteOpcode(op, gParser.previous.Line)
}

func emitOpcodes(op1, op2 bytecode.OpCode) {
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
	if constant > MaxConstantCount {
		errorAtPrev("Too many constants in one chunk.")
		return 0
	}
	return byte(constant)
}

func emitReturn() {
	emitOpcode(bytecode.OpNil)
	emitOpcode(bytecode.OpReturn)
}

func endCompiler() {
	emitReturn()
}

func beginScope() {
	gCurrent.ScoreDepth++
}

func endScope() {
	gCurrent.ScoreDepth--

	for gCurrent.LocalCount > 0 &&
		gCurrent.Locals[gCurrent.LocalCount-1].Depth > gCurrent.ScoreDepth {
		emitOpcode(bytecode.OpPop)
		gCurrent.LocalCount--
	}
}
