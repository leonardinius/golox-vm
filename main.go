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
	chunk.WriteOpcode(vmchunk.OpReturn, 1)

	vmdebug.DisassembleChunk(&chunk, "test chunk")
	vm.Interpret(&chunk)

	vm.FreeVM()
	chunk.FreeChunk()
	fmt.Println("main > ", strings.Join(args, " "))
	_ = repl()

	os.Exit(0)
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
