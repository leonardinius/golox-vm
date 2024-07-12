package vm

import (
	"errors"
	"fmt"

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

type InterpretResult int

const (
	InterpretRuntimeError InterpretResult = iota
	InterpretSuccess
)

var errRuntimeError = errors.New("runtime error")

func InitVM() {
	resetStack()
	resetChunk()
}

func FreeVM() {
	resetStack()
	resetChunk()
}

func resetStack() {
	GlobalVM.StackTop = 0
}

func initVMChunk(chunk *vmchunk.Chunk) {
	GlobalVM.Chunk = chunk
	GlobalVM.IP = 0
}

func resetChunk() {
	if GlobalVM.Chunk != nil {
		GlobalVM.Chunk.FreeChunk()
		GlobalVM.Chunk = nil
	}
	GlobalVM.IP = 0
}

func Interpret(script string, code []byte) (vmvalue.Value, error) {
	chunk, err := vmcompiler.Compile(code)
	if err != nil {
		return vmvalue.NilValue, fmt.Errorf("compile %s: %w", script, err)
	}

	initVMChunk(&chunk)
	defer resetChunk()

	vmdebug.DisassembleChunk(GlobalVM.Chunk, script)
	if value, result := Run(); result == InterpretRuntimeError {
		return value, errRuntimeError
	} else {
		return value, nil
	}
}

func Run() (vmvalue.Value, InterpretResult) {
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

		instruction := vmchunk.OpCode(readByte())
		switch instruction {
		case vmchunk.OpConstant:
			constant := readConstant()
			Push(constant)

		case vmchunk.OpAdd:
			binaryOpAdd()

		case vmchunk.OpSubtract:
			binaryOpSubtract()

		case vmchunk.OpMultiply:
			binaryOpMultiply()

		case vmchunk.OpDivide:
			binaryOpDivide()

		case vmchunk.OpNegate:
			Push(-Pop())

		case vmchunk.OpPop:
			Pop()

		case vmchunk.OpReturn:
			value := Pop()
			return value, InterpretSuccess

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

func binaryOpAdd() {
	b := Pop()
	a := Pop()
	Push(a + b)
}

func binaryOpSubtract() {
	b := Pop()
	a := Pop()
	Push(a - b)
}

func binaryOpMultiply() {
	b := Pop()
	a := Pop()
	Push(a * b)
}

func binaryOpDivide() {
	b := Pop()
	a := Pop()
	Push(a / b)
}

func readByte() byte {
	GlobalVM.IP++
	return GlobalVM.Chunk.Code[GlobalVM.IP-1]
}

func readConstant() vmvalue.Value {
	return GlobalVM.Chunk.Constants.At(int(readByte()))
}
