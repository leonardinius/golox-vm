package main

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/leonardinius/goloxvm/internal/cmd"
	"github.com/leonardinius/goloxvm/internal/vm/vmdebug"
)

func main() {
	os.Exit(mainCli(os.Args[1:]...))
}

func mainCli(args ...string) int {
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

	code := cmd.Main(args...)

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

func ioClose(c io.Closer) {
	if err := c.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "[WARN ] close: %s\n", err)
	}
}
