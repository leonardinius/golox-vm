package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"

	"github.com/chzyer/readline"

	"github.com/leonardinius/goloxvm/internal/vm"
	"github.com/leonardinius/goloxvm/internal/vm/vmdebug"
)

func main() {
	os.Exit(mainApp())
}

func mainApp() int {
	cfgPproff := vmdebug.LoadPprofConfigFromEnv()

	if cfgPproff.On && cfgPproff.CPUProfile != "" {
		f, err := os.Create(cfgPproff.CPUProfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR] could not create CPU profile: %#v", err)
			return 1
		}
		defer ioClose(f)
		if err := pprof.StartCPUProfile(f); err != nil {
			ioClose(f)
			fmt.Fprintf(os.Stderr, "[ERROR] could not start CPU profile: %#v", err)
			return 1
		}
		defer pprof.StopCPUProfile()
	}

	code := mainCli()

	if cfgPproff.On && cfgPproff.MemProfile != "" {
		f, err := os.Create(cfgPproff.MemProfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR] could not create memory profile: %#v", err)
			return code
		}
		defer ioClose(f)
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			ioClose(f)
			fmt.Fprintf(os.Stderr, "[ERROR] could not write memory profile: %#v", err)
			return code
		}
	}

	return code
}

// mainCli is the entry point for the CLI
// It handles the command line arguments and calls the appropriate functions
// It also initializes and frees the VM.
func mainCli() int {
	args := os.Args[1:]
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

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		return 65
	}

	return 0
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
	data, err := os.ReadFile(script) //nolint:gosec // get the data
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
