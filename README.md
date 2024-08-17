# (Go) Lox Bytecode VM (Crafting interpreters)

This is an implementation of [The Lox Programming Language](https://www.craftinginterpreters.com/the-lox-language.html) implemented in Go.

This is based on my previous work:

- [leonardinius/golox](https://github.com/leonardinius/golox)
  Golang Lox interpreter.
- [leonardinius/clox](https://github.com/leonardinius/clox)
  C language bytecode VM, basically just going slowly page by page in book and ^C^V the code from there.

This repository is re-implementation of what I have learned in clox VM bytecode with golang.

`golox-vm` is `clox` copy with Go syntax. I ended up copying the design mostly, the main effort went into how to make it work with Go cyclic dependency, GC and memory management, NaN boxing and unsafe pointers.

I feel I've learned a thing or two about (removed experiment) CGO with C.alloc/C.free memory management etc.

## What's included?

- golox-vm bytecode.
- mark & sweep garbage collector, NaN boxed values.
- LOX: functions, OOP etc..
- acceptance tests: ^C^V from [munificent/craftinginterpreters:test/](https://github.com/munificent/craftinginterpreters/tree/master/test).
- benchmarks.
- pprof profiler support: `GLOX_PPROF`=0/1,`GLOX_PPROF_CPU`=0/1,`GLOX_PPROF_MEM`=0/1

## Completeess & Speed

- `make test_e2e` pass all original test suite.
- `make bench` executes basic tests. At the moment the speed is better than [leonardinius/golox](https://github.com/leonardinius/golox) but not as fast as [leonardinius/clox](https://github.com/leonardinius/clox).
  More detailed analysis is due.

## How to XY locally?

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
