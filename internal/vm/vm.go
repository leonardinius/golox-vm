package vm

import (
	"fmt"
	"os"

	"github.com/leonardinius/goloxvm/internal/vm/bytecode"
	"github.com/leonardinius/goloxvm/internal/vm/hashtable"
	"github.com/leonardinius/goloxvm/internal/vm/vmchunk"
	"github.com/leonardinius/goloxvm/internal/vm/vmdebug"
	"github.com/leonardinius/goloxvm/internal/vm/vmmem"
	"github.com/leonardinius/goloxvm/internal/vm/vmobject"
	"github.com/leonardinius/goloxvm/internal/vm/vmvalue"
	"github.com/leonardinius/goloxvm/internal/vmcompiler"
)

const StackMax = 256

// VM is the virtual machine.
type VM struct {
	Chunk    *vmchunk.Chunk
	IP       int
	Stack    [StackMax]vmvalue.Value
	StackTop int
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
	resetStack()
	resetRootObjects()
	resetVMChunk()
}

func FreeVM() {
	hashtable.FreeGlobals()
	hashtable.InitInternStrings()
	vmobject.FreeObjects()
	resetStack()
	resetRootObjects()
	resetVMChunk()
}

func resetStack() {
	GlobalVM.StackTop = 0
}

func resetRootObjects() {
	vmobject.GRoots = nil
}

func initVMChunk(chunk *vmchunk.Chunk) {
	GlobalVM.Chunk = chunk
	GlobalVM.IP = 0
}

func resetVMChunk() {
	GlobalVM.Chunk = nil
	GlobalVM.IP = 0
}

func Interpret(script string, code []byte) (vmvalue.Value, error) {
	chunk := vmchunk.NewChunk()
	defer chunk.Free()

	if !vmcompiler.Compile(code, &chunk) {
		return vmvalue.NilValue, InterpretCompileError
	}

	initVMChunk(&chunk)
	defer resetVMChunk()

	vmdebug.DisassembleChunk(GlobalVM.Chunk, script)
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
	vmdebug.DisassembleInstruction(GlobalVM.Chunk, GlobalVM.IP)
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

func SetGlobal(name *vmobject.ObjString, value vmvalue.Value) {
	hashtable.SetGlobal(name, value)
}

func GetGlobal(name *vmobject.ObjString) (vmvalue.Value, bool) {
	return hashtable.GetGlobal(name)
}

func GCObjects() *vmobject.Obj {
	return vmobject.GRoots
}

func Run() (vmvalue.Value, error) {
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
			Push(vmvalue.BoolAsValue(!isFalsey(Pop())))
		case bytecode.OpPop:
			Pop()
		case bytecode.OpPrint:
			PrintlnValue(Pop())
		case bytecode.OpGetGlobal:
			name := readString()
			if value, gok := GetGlobal(name); !gok {
				ok = runtimeError("Undefined variable '%s'.", string(name.Chars))
			} else {
				Push(value)
			}
		case bytecode.OpDefineGlobal:
			name := readString()
			SetGlobal(name, Peek(0))
			Pop()
		case bytecode.OpReturn:
			return Pop(), nil
		default:
			ok = runtimeError("Unexpected instruction")
		}
	}
}

func isFalsey(value vmvalue.Value) bool {
	if vmvalue.IsBool(value) {
		return vmvalue.ValueAsBool(value)
	}
	return !vmvalue.IsNil(value)
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
	b := vmvalue.ValueAsString(Pop())
	a := vmvalue.ValueAsString(Pop())
	chars := vmmem.AllocateSlice[byte](len(a.Chars) + len(b.Chars))
	copy(chars, a.Chars)
	copy(chars[len(a.Chars):], b.Chars)
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

func readByte() byte {
	GlobalVM.IP++
	return GlobalVM.Chunk.Code[GlobalVM.IP-1]
}

func readConstant() vmvalue.Value {
	return GlobalVM.Chunk.Constants.At(int(readByte()))
}

func readString() *vmobject.ObjString {
	return vmvalue.ValueAsString(readConstant())
}

func runtimeError(format string, messageAndArgs ...any) (ok bool) {
	fmt.Fprintf(os.Stderr, format, messageAndArgs...)
	fmt.Fprintln(os.Stderr)

	offset := GlobalVM.IP - 1
	line := GlobalVM.Chunk.Lines.GetLineByOffset(offset)
	fmt.Fprintf(os.Stderr, "[line %d] in script\n", line)
	resetStack()
	return false
}

func PrintlnValue(v vmvalue.Value) {
	vmdebug.PrintlnValue(v)
}
