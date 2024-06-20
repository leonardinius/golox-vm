package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"

	"github.com/leonardinius/goloxvm/internal/vmchunk"
	"github.com/leonardinius/goloxvm/internal/vmdebug"
)

func main() {
	args := os.Args[1:]
	chunk := vmchunk.NewChunk()
	chunk.InitChunk()

	constant := chunk.AddConstant(1.2)
	chunk.Write(vmchunk.OpConstant)
	chunk.Write1(byte(constant))
	chunk.Write(vmchunk.OpReturn)

	vmdebug.DisassembleChunk(&chunk, "test chunk")

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
