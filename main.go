package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"

	"github.com/leonardinius/goloxvm/internal/chunk"
	"github.com/leonardinius/goloxvm/internal/debug"
)

func main() {
	args := os.Args[1:]
	ch := chunk.NewChunk()
	ch.InitChunk()
	ch.Write(chunk.OpReturn)

	debug.DisassembleChunk(&ch, "test chunk")

	ch.FreeChunk()
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
