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
	MaxArity         = math.MaxUint8
	MaxConstantCount = math.MaxUint8 + 1
	MaxLocalCount    = math.MaxUint8 + 1
	MaxUpvalueCount  = math.MaxUint8 + 1
	MaxJump          = math.MaxUint16
)

type FunctionType int

const (
	_ FunctionType = iota
	FunctionTypeScript
	FunctionTypeFunction
	FunctionTypeMethod
	FunctionTypeInitializer
)

type Compiler struct {
	Chunk    vmchunk.Chunk
	Function *vmvalue.ObjFunction
	FnType   FunctionType

	Locals     [MaxLocalCount]Local
	LocalCount int
	ScoreDepth int

	Upvalues [MaxUpvalueCount]Upvalue

	Enclosing *Compiler
}

type ClassCompiler struct {
	Enclosing *ClassCompiler
}

type Local struct {
	Name       scanner.Token
	Depth      int
	IsCaptured bool
}

func (l *Local) SetName(name string) {
	l.Name.Source = []byte(name)
	l.Name.Start = 0
	l.Name.Length = len(l.Name.Source)
}

type Upvalue struct {
	Index int
	Local byte
}

func NewCompiler(fnType FunctionType, fnName *vmvalue.ObjString) *Compiler {
	chunk := vmchunk.NewChunk()
	compiler := Compiler{}
	compiler.Chunk = chunk
	compiler.FnType = fnType
	compiler.Function = vmvalue.NewFunction(chunk.AsPtr(), chunk.Free, chunk.Mark)
	compiler.Function.Name = fnName
	compiler.Enclosing = gCurrent
	gCurrent = &compiler

	compiler.LocalCount = 0
	local := &compiler.Locals[compiler.LocalCount]
	compiler.LocalCount++
	local.Depth = 0
	local.IsCaptured = false
	if fnType != FunctionTypeFunction {
		local.SetName("this")
	} else {
		local.SetName("")
	}
	return &compiler
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

func emitByte(b byte) {
	currentChunk().Write(b, gParser.previous.Line)
}

func emitOpByte(op bytecode.OpCode, b byte) {
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
	emitOpByte(bytecode.OpConstant, byte(makeConstant(v)))
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

func makeConstant(v vmvalue.Value) int {
	constant := currentChunk().AddConstant(v)
	if constant >= MaxConstantCount {
		errorAtPrev("Too many constants in one chunk.")
		return 0
	}
	return constant
}

func emitReturn() {
	if gCurrent.FnType == FunctionTypeInitializer {
		emitOpByte(bytecode.OpGetLocal, 0)
	} else {
		emitOpcode(bytecode.OpNil)
	}
	emitOpcode(bytecode.OpReturn)
}

func endCompiler() *vmvalue.ObjFunction {
	emitReturn()
	fn := gCurrent.Function
	gCurrent = gCurrent.Enclosing
	if vmdebug.DebugDisassembler && !gParser.hadError {
		disassembleFunction(fn)
	}
	return fn
}

func beginScope() {
	gCurrent.ScoreDepth++
}

func endScope() {
	gCurrent.ScoreDepth--

	for gCurrent.LocalCount > 0 {
		local := &gCurrent.Locals[gCurrent.LocalCount-1]
		if local.Depth <= gCurrent.ScoreDepth {
			break
		}
		if local.IsCaptured {
			emitOpcode(bytecode.OpCloseUpvalue)
		} else {
			emitOpcode(bytecode.OpPop)
		}
		gCurrent.LocalCount--
	}
}

func disassembleFunction(fn *vmvalue.ObjFunction) {
	fnName := "<script>"
	if fn.Name != nil {
		fnName = string(fn.Name.Chars)
	}
	chunk := vmchunk.FromPtr(fn.Chunk)
	vmdebug.DisassembleChunk(chunk, fnName)
}

func MarkCompilerRoots() {
	compiler := gCurrent
	for compiler != nil {
		vmvalue.MarkObject(compiler.Function)
		compiler = compiler.Enclosing
	}
}
