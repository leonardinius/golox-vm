package vm

import (
	"fmt"
	"os"

	"github.com/leonardinius/goloxvm/internal/bytecode"
	"github.com/leonardinius/goloxvm/internal/vmchunk"
	"github.com/leonardinius/goloxvm/internal/vmcompiler"
	"github.com/leonardinius/goloxvm/internal/vmdebug"
	"github.com/leonardinius/goloxvm/internal/vmvalue"
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
	resetStack()
	resetVMChunk()
}

func FreeVM() {
	resetStack()
	resetVMChunk()
}

func resetStack() {
	GlobalVM.StackTop = 0
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
	chunk.InitChunk()
	defer chunk.Free()

	if !vmcompiler.Compile(code, &chunk) {
		return vmvalue.NilValue, InterpretCompileError
	}

	initVMChunk(&chunk)
	defer resetVMChunk()

	vmdebug.DisassembleChunk(GlobalVM.Chunk, script)
	return Run()
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
		case bytecode.OpAdd:
			ok = binaryOp(binOpAdd)
		case bytecode.OpSubtract:
			ok = binaryOp(binOpSubtract)
		case bytecode.OpMultiply:
			ok = binaryOp(binOpMultiply)
		case bytecode.OpDivide:
			ok = binaryOp(binOpDivide)
		case bytecode.OpNegate:
			ok = opNegate()
		case bytecode.OpPop:
			Pop()
		case bytecode.OpReturn:
			value := Pop()
			return value, nil
		default:
			ok = runtimeError("Unexpected instruction")
		}
	}
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

func binaryOp(op func(float64, float64) float64) (ok bool) {
	if ok = vmvalue.IsNumber(Peek(0)) && vmvalue.IsNumber(Peek(1)); !ok {
		runtimeError("Operands must be numbers.")
		return ok
	}

	b := vmvalue.ValueAsNumber(Pop())
	a := vmvalue.ValueAsNumber(Pop())
	Push(vmvalue.NumberValue(op(a, b)))
	return ok
}

func opNegate() (ok bool) {
	if ok = vmvalue.IsNumber(Peek(0)); !ok {
		runtimeError("Operand must be a number.")
		return ok
	}
	Push(vmvalue.NumberValue(-vmvalue.ValueAsNumber(Pop())))
	return ok
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

func readByte() byte {
	GlobalVM.IP++
	return GlobalVM.Chunk.Code[GlobalVM.IP-1]
}

func readConstant() vmvalue.Value {
	return GlobalVM.Chunk.Constants.At(int(readByte()))
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
