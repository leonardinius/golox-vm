package vm

import (
	"fmt"
	"math"
	"os"
	"unsafe"

	"github.com/leonardinius/goloxvm/internal/vm/bytecode"
	"github.com/leonardinius/goloxvm/internal/vm/hashtable"
	"github.com/leonardinius/goloxvm/internal/vm/vmchunk"
	"github.com/leonardinius/goloxvm/internal/vm/vmdebug"
	"github.com/leonardinius/goloxvm/internal/vm/vmmem"
	"github.com/leonardinius/goloxvm/internal/vm/vmobject"
	"github.com/leonardinius/goloxvm/internal/vm/vmvalue"
	"github.com/leonardinius/goloxvm/internal/vmcompiler"
)

const (
	MaxCallFrames = 64
	MaxStackCount = MaxCallFrames * (math.MaxUint8 + 1)
)

type CallFrame struct {
	Function *vmobject.ObjFunction
	IP       int
	Slots    []vmvalue.Value
}

// VM is the virtual machine.
type VM struct {
	Frames     [MaxCallFrames]CallFrame
	FrameCount int
	Stack      [MaxStackCount]vmvalue.Value
	StackTop   int
}

var GlobalVM VM

type InterpretError int

const (
	_ InterpretError = iota
	InterpretCompileError
	InterpretRuntimeError
)

func (i InterpretError) Error() string {
	var err string
	switch i {
	case InterpretCompileError:
		err = "compile error"
	case InterpretRuntimeError:
		err = "runtime error"
	default:
		err = fmt.Sprintf("unknown error %d", i)
	}

	return err
}

func InitVM() {
	hashtable.InitInternStrings()
	hashtable.InitGlobals()
	resetVM()
}

func FreeVM() {
	hashtable.FreeGlobals()
	hashtable.FreeInternStrings()
	vmobject.FreeObjects()
	resetVM()
}

func resetVM() {
	GlobalVM.StackTop = 0
	vmobject.GRoots = nil
	GlobalVM.FrameCount = 0
}

func Interpret(code []byte) (vmvalue.Value, error) {
	chunk := vmchunk.NewChunk()
	defer chunk.Free()

	var fn *vmobject.ObjFunction
	var ok bool

	if fn, ok = vmcompiler.Compile(code, &chunk); !ok {
		return vmvalue.NilValue, InterpretCompileError
	}

	fnName := "<script>"
	if fn.Name != nil {
		fnName = string(fn.Name.Chars)
	}

	Push(vmvalue.ObjAsValue(fn))
	nextFrame := &GlobalVM.Frames[GlobalVM.FrameCount]
	nextFrame.Function = fn
	nextFrame.IP = 0
	nextFrame.Slots = GlobalVM.Stack[:]
	GlobalVM.FrameCount++

	vmdebug.DisassembleChunk(&chunk, fnName)
	return Run()
}

func debug0() {
	if GlobalVM.StackTop > 0 {
		fmt.Print("          ")
		for i := range GlobalVM.StackTop {
			fmt.Print("[ ")
			vmdebug.PrintValue(GlobalVM.Stack[i])
			fmt.Print(" ]")
		}
		fmt.Println()
	}
	frame, chunk := frameChunk()
	vmdebug.DisassembleInstruction(chunk, frame.IP)
}

func Push(value vmvalue.Value) {
	GlobalVM.Stack[GlobalVM.StackTop] = value
	GlobalVM.StackTop++
}

func Pop() vmvalue.Value {
	GlobalVM.StackTop--
	return GlobalVM.Stack[GlobalVM.StackTop]
}

func Peek(distance int) vmvalue.Value {
	return GlobalVM.Stack[GlobalVM.StackTop-1-distance]
}

func SetGlobal(name *vmobject.ObjString, value vmvalue.Value) bool {
	return hashtable.SetGlobal(name, value)
}

func GetGlobal(name *vmobject.ObjString) (vmvalue.Value, bool) {
	return hashtable.GetGlobal(name)
}

func DeleteGlobal(name *vmobject.ObjString) bool {
	return hashtable.DeleteGlobal(name)
}

func GCObjects() *vmobject.Obj {
	return vmobject.GRoots
}

func Run() (vmvalue.Value, error) { //nolint:gocyclo // expected high complexity in Run switch
	if vmdebug.DebugDisassembler {
		fmt.Println()
		fmt.Println("== trace execution ==")

		defer fmt.Println()
	}

	ok := true
	for {
		if !ok {
			return vmvalue.NilValue, InterpretRuntimeError
		}
		if vmdebug.DebugDisassembler {
			debug0()
		}

		instruction := bytecode.OpCode(readByte())
		switch instruction {
		case bytecode.OpConstant:
			constant := readConstant()
			Push(constant)
		case bytecode.OpNil:
			Push(vmvalue.NilValue)
		case bytecode.OpTrue:
			Push(vmvalue.TrueValue)
		case bytecode.OpFalse:
			Push(vmvalue.FalseValue)
		case bytecode.OpEqual:
			Push(vmvalue.BoolAsValue(vmvalue.IsValuesEqual(Pop(), Pop())))
		case bytecode.OpGreater:
			ok = binaryNumCompareOp(binOpGreater)
		case bytecode.OpLess:
			ok = binaryNumCompareOp(binOpLess)
		case bytecode.OpAdd:
			if vmvalue.IsString(Peek(0)) && vmvalue.IsString(Peek(1)) {
				ok = stringConcat()
			} else if vmvalue.IsNumber(Peek(0)) && vmvalue.IsNumber(Peek(1)) {
				ok = binaryNumMathOp(binOpAdd)
			} else {
				ok = runtimeError("Operands must be two numbers or two strings.")
			}
		case bytecode.OpSubtract:
			ok = binaryNumMathOp(binOpSubtract)
		case bytecode.OpMultiply:
			ok = binaryNumMathOp(binOpMultiply)
		case bytecode.OpDivide:
			ok = binaryNumMathOp(binOpDivide)
		case bytecode.OpNegate:
			ok = opNegate()
		case bytecode.OpNot:
			Push(vmvalue.BoolAsValue(!isTruey(Pop())))
		case bytecode.OpPop:
			Pop()
		case bytecode.OpPrint:
			PrintlnValue(Pop())
		case bytecode.OpGetLocal:
			frame := currentFrame()
			slot := readByte()
			Push(frame.Slots[slot])
		case bytecode.OpSetLocal:
			frame := currentFrame()
			slot := readByte()
			frame.Slots[slot] = Peek(0)
		case bytecode.OpGetGlobal:
			name := readString()
			if value, gok := GetGlobal(name); !gok {
				ok = runtimeError("Undefined variable '%s'.", string(name.Chars))
			} else {
				Push(value)
			}
		case bytecode.OpSetGlobal:
			name := readString()
			if isNewKey := SetGlobal(name, Peek(0)); isNewKey {
				DeleteGlobal(name)
				ok = runtimeError("Undefined variable '%s'.", string(name.Chars))
			}
		case bytecode.OpDefineGlobal:
			name := readString()
			SetGlobal(name, Peek(0))
			Pop()
		case bytecode.OpJump:
			frame := currentFrame()
			offset := readShort()
			frame.IP += int(offset)
		case bytecode.OpJumpIfFalse:
			frame := currentFrame()
			offset := readShort()
			if isFalsey(Peek(0)) {
				frame.IP += int(offset)
			}
		case bytecode.OpLoop:
			frame := currentFrame()
			offset := readShort()
			frame.IP -= int(offset)
		case bytecode.OpReturn:
			return Pop(), nil
		default:
			ok = runtimeError("Unexpected instruction")
		}
	}
}

func isTruey(value vmvalue.Value) bool {
	if vmvalue.IsBool(value) {
		return vmvalue.ValueAsBool(value)
	}
	return !vmvalue.IsNil(value)
}

func isFalsey(value vmvalue.Value) bool {
	return !isTruey(value)
}

func binaryNumOp(op func(vmvalue.Value, vmvalue.Value) vmvalue.Value) (ok bool) {
	if ok = vmvalue.IsNumber(Peek(0)) && vmvalue.IsNumber(Peek(1)); !ok {
		runtimeError("Operands must be numbers.")
		return ok
	}

	b := Pop()
	a := Pop()
	Push(op(a, b))
	return ok
}

func binaryNumMathOp(op func(float64, float64) float64) (ok bool) {
	return binaryNumOp(func(a vmvalue.Value, b vmvalue.Value) vmvalue.Value {
		av := vmvalue.ValueAsNumber(a)
		bv := vmvalue.ValueAsNumber(b)
		return vmvalue.NumberAsValue(op(av, bv))
	})
}

func binaryNumCompareOp(op func(float64, float64) bool) (ok bool) {
	return binaryNumOp(func(a vmvalue.Value, b vmvalue.Value) vmvalue.Value {
		av := vmvalue.ValueAsNumber(a)
		bv := vmvalue.ValueAsNumber(b)
		return vmvalue.BoolAsValue(op(av, bv))
	})
}

func opNegate() (ok bool) {
	if ok = vmvalue.IsNumber(Peek(0)); !ok {
		runtimeError("Operand must be a number.")
		return ok
	}
	Push(vmvalue.NumberAsValue(-vmvalue.ValueAsNumber(Pop())))
	return ok
}

func stringConcat() (ok bool) {
	b := vmvalue.ValueAsStringChars(Pop())
	a := vmvalue.ValueAsStringChars(Pop())
	length := len(a) + len(b)
	chars := vmmem.AllocateSlice[byte](length)
	copy(chars, a)
	copy(chars[len(a):], b)
	str := hashtable.StringInternTake(chars)
	Push(vmvalue.ObjAsValue(str))
	return true
}

func binOpAdd(a, b float64) float64 {
	return a + b
}

func binOpSubtract(a, b float64) float64 {
	return a - b
}

func binOpMultiply(a, b float64) float64 {
	return a * b
}

func binOpDivide(a, b float64) float64 {
	return a / b
}

func binOpGreater(a, b float64) bool {
	return a > b
}

func binOpLess(a, b float64) bool {
	return a < b
}

func currentFrame() *CallFrame {
	return &GlobalVM.Frames[GlobalVM.FrameCount-1]
}

func frameChunk() (*CallFrame, *vmchunk.Chunk) {
	frame := &GlobalVM.Frames[GlobalVM.FrameCount-1]
	ptr := frame.Function.ChunkPtr
	ch := (**vmchunk.Chunk)(unsafe.Pointer(&ptr)) //nolint:gosec // unsafe.Pointer is used here
	return frame, *ch
}

func readByte() byte {
	frame, chunk := frameChunk()
	frame.IP++
	return chunk.Code[frame.IP-1]
}

func readShort() uint16 {
	frame, chunk := frameChunk()
	frame.IP += 2
	return (uint16(chunk.Code[frame.IP-2]) << 8) | uint16(chunk.Code[frame.IP-1])
}

func readConstant() vmvalue.Value {
	frame, chunk := frameChunk()
	frame.IP++
	at := chunk.Code[frame.IP-1]
	return chunk.ConstantAt(int(at))
}

func readString() *vmobject.ObjString {
	return vmvalue.ValueAsString(readConstant())
}

func runtimeError(format string, messageAndArgs ...any) (ok bool) {
	fmt.Fprintf(os.Stderr, format, messageAndArgs...)
	fmt.Fprintln(os.Stderr)

	frame, chunk := frameChunk()
	offset := frame.IP - 1
	line := chunk.Lines.GetLineByOffset(offset)
	fmt.Fprintf(os.Stderr, "[line %d] in script\n", line)
	resetVM()
	return false
}

func PrintlnValue(v vmvalue.Value) {
	vmdebug.PrintlnValue(v)
}
