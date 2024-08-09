package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/chzyer/readline"

	"github.com/leonardinius/goloxvm/internal/vm"
)

// Main is main entry point for the GoLox-VM
// It takes the command line arguments and calls the appropriate functions
// It also initializes and frees the VM.
func Main(args ...string) int {
	vm.InitVM()

	var err error
	if len(args) == 0 {
		fmt.Println("Welcome to the GoLox-VM REPL!")
		err = repl("repl")
	} else if len(args) == 1 {
		err = runFile(args[0])
	} else {
		fmt.Printf("Usage: %s [path]\n", filepath.Base(os.Args[0]))
		return 64
	}

	vm.FreeVM()

	if err == nil {
		return 0
	}

	// interpreter reports errors to stderr
	switch err {
	case vm.InterpretRuntimeError:
		return 70
	case vm.InterpretCompileError:
		return 65
	default:
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		return 65
	}
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

		if value, err := vm.Interpret(line); err == nil {
			vm.PrintlnValue(value)
		}
		// else {
		// Do nothing
		// interpreter reports errors to stderr
		// fmt.Println(err)
		//}
	}
}

func runFile(script string) error {
	data, err := os.ReadFile(script) //nolint:gosec
	if err == nil {
		_, err = vm.Interpret(data)
	}
	return err
}

func ioClose(c io.Closer) {
	if err := c.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "[WARN ] close: %s\n", err)
	}
}
