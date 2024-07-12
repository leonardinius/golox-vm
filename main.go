package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/chzyer/readline"

	"github.com/leonardinius/goloxvm/internal/vm"
	"github.com/leonardinius/goloxvm/internal/vmdebug"
)

func main() {
	args := os.Args[1:]
	vm.InitVM()

	var err error
	if len(args) == 0 {
		fmt.Println("Welcome to the Lox REPL!")
		err = repl("main")
	} else if len(args) == 1 {
		err = runFile(args[0])
	} else {
		fmt.Printf("Usage: %s [path]\n", filepath.Base(os.Args[0]))
		os.Exit(64)
	}

	vm.FreeVM()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(65)
	}

	os.Exit(0)
}

func repl(welcome string) error {
	rl, err := readline.New(welcome + "> ")
	if err != nil {
		return err
	}
	defer ioClose(rl)

	for {
		line, err := rl.ReadSlice()
		if err != nil {
			return err
		}

		if value, err := vm.Interpret(welcome, line); err == nil {
			vmdebug.PrintlnValue(value)
		} else {
			fmt.Println(err)
		}
	}
}

func runFile(script string) error {
	data, err := os.ReadFile(script) //nolint:gosec // get the data
	if err == nil {
		_, err = vm.Interpret(filepath.Base(script), data)
	}
	return err
}

func ioClose(c io.Closer) {
	if err := c.Close(); err != nil {
		fmt.Printf("[WARN] close: %s\n", err)
	}
}
