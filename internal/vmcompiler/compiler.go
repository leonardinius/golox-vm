package vmcompiler

import (
	"math"

	"github.com/leonardinius/goloxvm/internal/vm/bytecode"
	"github.com/leonardinius/goloxvm/internal/vm/vmchunk"
	"github.com/leonardinius/goloxvm/internal/vm/vmdebug"
	"github.com/leonardinius/goloxvm/internal/vm/vmvalue"
	"github.com/leonardinius/goloxvm/internal/vmcompiler/scanner"
	"github.com/leonardinius/goloxvm/internal/vmcompiler/tokens"
)

const (
	MaxConstantCount = math.MaxUint8
	MaxLocalCount    = math.MaxUint8 + 1
	MaxJump          = math.MaxUint16
)

type FunctionType int

const (
	_ FunctionType = iota
	FunctionTypeFunction
	FunctionTypeScript
)

type Compiler struct {
	Function *vmvalue.ObjFunction
	FnType   FunctionType

	Locals     [MaxLocalCount]Local
	LocalCount int
	ScoreDepth int

	Enclosing *Compiler
}

type Local struct {
	Name  scanner.Token
	Depth int
}

func (l *Local) SetName(name string) {
	l.Name.Source = []byte(name)
	l.Name.Start = 0
	l.Name.Length = len(l.Name.Source)
}

func NewCompiler(fnType FunctionType, fnName *vmvalue.ObjString) Compiler {
	chunk := vmchunk.NewChunk()
	compiler := Compiler{}
	compiler.FnType = fnType
	compiler.Function = vmvalue.NewFunction(chunk.AsPtr(), chunk.Free)
	compiler.Function.Name = fnName
	compiler.Enclosing = gCurrent
	gCurrent = &compiler

	local := &gCurrent.Locals[gCurrent.LocalCount]
	gCurrent.LocalCount++
	local.Depth = 0
	local.SetName("")
	return compiler
}

func Compile(source []byte) (*vmvalue.ObjFunction, bool) {
	gScanner = scanner.NewScanner(source)
	defer gScanner.Free()
	gParser = NewParser()

	_ = NewCompiler(FunctionTypeScript, nil)

	advance()

	for !match(tokens.TokenEOF) {
		declaration()
	}

	fn := endCompiler()
	return fn, !gParser.hadError
}

func currentChunk() *vmchunk.Chunk {
	return vmchunk.FromPtr(gCurrent.Function.Chunk)
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

func emitJump(op bytecode.OpCode) int {
	emitOpcode(op)
	currentChunk().Write(0xff, gParser.previous.Line)
	currentChunk().Write(0xff, gParser.previous.Line)
	return currentChunk().Count - 2
}

func emitLoop(loopStart int) {
	emitOpcode(bytecode.OpLoop)

	offset := currentChunk().Count - loopStart + 2
	if offset > MaxJump {
		errorAtPrev("Loop body too large.")
	}

	b1 := byte((offset >> 8) & 0xff)
	b2 := byte((offset) & 0xff)
	currentChunk().Write(b1, gParser.previous.Line)
	currentChunk().Write(b2, gParser.previous.Line)
}

func emitConstant(v vmvalue.Value) {
	emitBytes(bytecode.OpConstant, makeConstant(v))
}

func patchJump(offset int) {
	// -2 to adjust for the bytecode for the jump offset itself.
	jump := currentChunk().Count - offset - 2

	if jump > MaxJump {
		errorAtPrev("Too much code to jump over.")
	}

	b1 := byte((jump >> 8) & 0xff)
	b2 := byte((jump) & 0xff)

	currentChunk().Code[offset] = b1
	currentChunk().Code[offset+1] = b2
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

func endCompiler() *vmvalue.ObjFunction {
	emitReturn()
	fn := gCurrent.Function
	gCurrent = gCurrent.Enclosing
	if !gParser.hadError {
		fnName := "<script>"
		if fn.Name != nil {
			fnName = string(fn.Name.Chars)
		}
		chunk := vmchunk.FromPtr(fn.Chunk)
		vmdebug.DisassembleChunk(chunk, fnName)
	}
	return fn
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
