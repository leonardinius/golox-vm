# (Go) Lox Bytecode VM (Crafting Interpreters)

This is an implementation of [The Lox Programming Language](https://www.craftinginterpreters.com/the-lox-language.html) in Go.

This builds upon my previous work:

- [leonardinius/golox](https://github.com/leonardinius/golox) - Golang Lox interpreter.
- [leonardinius/clox](https://github.com/leonardinius/clox) - C language bytecode VM, following the "Crafting Interpreters" book closely.

This repository is a re-implementation of what I learned from the clox VM bytecode using Go.

`golox-vm` is not a direct copy of `clox`, but rather a reimplementation using Go syntax. The main effort went into dealing with Go's limitations like cyclic dependencies, garbage collection, memory management, NaN boxing, and unsafe pointers.

I've also learned a thing or two about (removed experiment) CGO (Calling Go from C) with C.alloc/C.free memory management, etc.

## What's Included?

* `golox-vm` bytecode
* Mark & sweep garbage collector with NaN boxed values
* LOX features: functions, OOP, etc.
* Acceptance tests: Copied from [munificent/craftinginterpreters:test/](https://github.com/munificent/craftinginterpreters/tree/master/test)
* Benchmarks
* pprof profiler support: `GLOX_PPROF`=0/1,`GLOX_PPROF_CPU`=0/1,`GLOX_PPROF_MEM`=0/1

## Completeness & Speed

* `make test_e2e` passes all original test suites.
* `make bench` executes basic tests. Currently, the speed is better than [leonardinius/golox](https://github.com/leonardinius/golox) but not as fast as [leonardinius/clox](https://github.com/leonardinius/clox).
* More detailed analysis is needed.

## How to Run Locally

```terminal
λ make 
Usage: make <target>
 Default
        help                  Display this help
 Build/Run
        all                   ALL, builds the world
        clean                 Clean-up build artifacts
        test                  Runs all tests
        test_e2e              Runs all e2e tests
        bench                 Runs all benchmarks
        lint                  Runs all linters
        release               Build RELEASE (debug off)
        debug                 Build DEBUG (debug on)
        run                   Runs golox-vm. Use 'make ARGS="script.lox" run' to pass arguments
        rund                  Runs goloxd-vm. Use 'make ARGS="script.lox" run' to pass arguments
```
