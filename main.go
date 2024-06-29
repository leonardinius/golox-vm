package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"

	"github.com/leonardinius/goloxvm/internal/vm"
	"github.com/leonardinius/goloxvm/internal/vmchunk"
	"github.com/leonardinius/goloxvm/internal/vmdebug"
	"github.com/leonardinius/goloxvm/internal/vmvalue"
)

func main() {
	args := os.Args[1:]
	vm.InitVM()
	chunk := vmchunk.NewChunk()
	chunk.InitChunk()

	constant1 := chunk.AddConstant(1.1)
	chunk.WriteOpcode(vmchunk.OpConstant, 1)
	chunk.Write(byte(constant1), 1)
	constant2 := chunk.AddConstant(1.2)
	chunk.WriteOpcode(vmchunk.OpConstant, 1)
	chunk.Write(byte(constant2), 1)
	chunk.WriteOpcode(vmchunk.OpNegate, 1)
	chunk.WriteOpcode(vmchunk.OpPop, 1)
	chunk.WriteOpcode(vmchunk.OpPop, 1)

	addBinOp(&chunk, vmchunk.OpAdd, 2.0, 3.0)
	chunk.WriteOpcode(vmchunk.OpPop, 1)
	addBinOp(&chunk, vmchunk.OpSubtract, 2.0, 3.0)
	chunk.WriteOpcode(vmchunk.OpPop, 1)
	addBinOp(&chunk, vmchunk.OpMultiply, 2.0, 3.0)
	chunk.WriteOpcode(vmchunk.OpPop, 1)
	addBinOp(&chunk, vmchunk.OpDivide, 2.0, 3.0)
	chunk.WriteOpcode(vmchunk.OpReturn, 1)

	vmdebug.DisassembleChunk(&chunk, "test chunk")
	vm.Interpret(&chunk)

	vm.FreeVM()
	chunk.FreeChunk()
	fmt.Println("main > ", strings.Join(args, " "))
	_ = repl()

	os.Exit(0)
}

func addBinOp(chunk *vmchunk.Chunk, op vmchunk.OpCode, a, b vmvalue.Value) {
	aConstant := chunk.AddConstant(a)
	bConstant := chunk.AddConstant(b)

	chunk.WriteOpcode(vmchunk.OpConstant, 1)
	chunk.Write(byte(aConstant), 1)
	chunk.WriteOpcode(vmchunk.OpConstant, 1)
	chunk.Write(byte(bConstant), 1)
	chunk.WriteOpcode(op, 1)
}

func repl() error {
	rl, err := readline.New("> ")
	if err != nil {
		return err
	}
	defer func() { _ = rl.Close() }()

	for {
		line, err := rl.Readline()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Printf("<< %s\n", line)
	}
}
