package vm

import (
	"fmt"

	"github.com/leonardinius/goloxvm/internal/vmchunk"
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

const (
	InterpretCompileError InterpretResult = iota
	InterpretRuntimeError
	InterpretSuccess
)

func InitVM() {
	resetStack()
}

func resetStack() {
	GlobalVM.StackTop = 0
}

func FreeVM() {
}

type InterpretResult int

func Interpret(chunk *vmchunk.Chunk) InterpretResult {
	GlobalVM.Chunk = chunk
	GlobalVM.IP = 0
	return Run()
}

func Run() InterpretResult {
	if vmdebug.DEBUG {
		fmt.Printf("")
		fmt.Println("== trace execution ==")
		return run(func() {
			fmt.Print("          ")
			for i := range GlobalVM.StackTop {
				fmt.Print("[ ")
				vmdebug.PrintValue(GlobalVM.Stack[i])
				fmt.Print(" ]")
			}
			fmt.Println()
			vmdebug.DissasembleInstruction(GlobalVM.Chunk, GlobalVM.IP)
		})
	}
	return run(nil)
}

func run(debugFn func()) InterpretResult {
	for {
		if debugFn != nil {
			debugFn()
		}
		instruction := vmchunk.OpCode(readByte())
		switch instruction {
		case vmchunk.OpReturn:
			vmdebug.PrintValue(Pop())
			fmt.Println()
			return InterpretSuccess
		case vmchunk.OpConstant:
			constant := readConstant()
			Push(constant)
		default:
			fmt.Printf("Unexpected instruction %d\n", instruction)
			return InterpretRuntimeError
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

func readByte() byte {
	GlobalVM.IP++
	return GlobalVM.Chunk.Code[GlobalVM.IP-1]
}

func readConstant() vmvalue.Value {
	return GlobalVM.Chunk.Constants.At(int(readByte()))
}
