package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/leonardinius/goloxvm/internal/chunk"
	"github.com/leonardinius/goloxvm/internal/disassembler"
)

func main() {
	args := os.Args[1:]
	ch := chunk.NewChunk()
	ch.InitChunk()
	ch.Write(chunk.OpReturn)

	disassembler.DisassembleChunk(&ch, "test chunk")

	ch.FreeChunk()
	fmt.Println("main > ", strings.Join(args, " "))
	os.Exit(0)
}
