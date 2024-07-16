package vm

import (
	"fmt"

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

	for {
		if vmdebug.DebugDisassembler {
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

		instruction := bytecode.OpCode(readByte())
		switch instruction {
		case bytecode.OpConstant:
			constant := readConstant()
			Push(constant)

		case bytecode.OpAdd:
			binaryOpAdd()

		case bytecode.OpSubtract:
			binaryOpSubtract()

		case bytecode.OpMultiply:
			binaryOpMultiply()

		case bytecode.OpDivide:
			binaryOpDivide()

		case bytecode.OpNegate:
			unaryOpNegate()

		case bytecode.OpPop:
			Pop()

		case bytecode.OpReturn:
			value := Pop()
			return value, nil

		default:
			fmt.Printf("Unexpected instruction %d\n", instruction)
			return vmvalue.NilValue, InterpretRuntimeError
		}
	}
}

func Push(value vmvalue.Value) {
	GlobalVM.Stack[GlobalVM.StackTop] = value
	GlobalVM.StackTop++
}

func Pop() vmvalue.Value {
	GlobalVM.StackTop--
	return GlobalVM.Stack[GlobalVM.StackTop]
}

func unaryOpNegate() {
	a := vmvalue.ValueAsNumber(Pop())
	Push(vmvalue.NumberValue(-a))
}

func binaryOpAdd() {
	b := vmvalue.ValueAsNumber(Pop())
	a := vmvalue.ValueAsNumber(Pop())
	Push(vmvalue.NumberValue(a + b))
}

func binaryOpSubtract() {
	b := vmvalue.ValueAsNumber(Pop())
	a := vmvalue.ValueAsNumber(Pop())
	Push(vmvalue.NumberValue(a - b))
}

func binaryOpMultiply() {
	b := vmvalue.ValueAsNumber(Pop())
	a := vmvalue.ValueAsNumber(Pop())
	Push(vmvalue.NumberValue(a * b))
}

func binaryOpDivide() {
	b := vmvalue.ValueAsNumber(Pop())
	a := vmvalue.ValueAsNumber(Pop())
	Push(vmvalue.NumberValue(a / b))
}

func readByte() byte {
	GlobalVM.IP++
	return GlobalVM.Chunk.Code[GlobalVM.IP-1]
}

func readConstant() vmvalue.Value {
	return GlobalVM.Chunk.Constants.At(int(readByte()))
}
