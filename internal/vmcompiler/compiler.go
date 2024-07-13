package vmcompiler

import (
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/leonardinius/goloxvm/internal/vmchunk"
	"github.com/leonardinius/goloxvm/internal/vmscanner"
	"github.com/leonardinius/goloxvm/internal/vmvalue"
)

var (
	scanner        vmscanner.Scanner
	parser         Parser
	compilingChunk *vmchunk.Chunk
)

func Compile(source []byte, chunk *vmchunk.Chunk) bool {
	scanner = vmscanner.NewScanner(source)
	defer scanner.Free()
	compilingChunk = chunk
	defer endCompiler()

	parser = NewParser()

	advance()
	expression()
	consume(vmscanner.TokenEOF, "Expect end of expression.")

	return !parser.hadError
}

func advance() {
	parser.previous = parser.current

	for {
		parser.current = scanner.ScanToken()
		if parser.current.Type != vmscanner.TokenError {
			break
		}
		// use TokenError lexeme as error message
		errorAtCurrent(parser.current.Lexeme())
	}
}

func expression() {
	advance()
	number()
}

func number() {
	v, err := strconv.ParseFloat(parser.previous.Lexeme(), 64)
	if err != nil {
		errorAtPrev(err.Error())
	}
	emitConstant(vmvalue.Value(v))
}

func consume(stype vmscanner.TokenType, message string) {
	if parser.current.Type == stype {
		advance()
		return
	}

	errorAtCurrent(message)
}

func errorAtCurrent(message string) {
	errorAt(&parser.current, message)
}

func errorAtPrev(message string) {
	errorAt(&parser.previous, message)
}

func errorAt(token *vmscanner.Token, message string) {
	if parser.panicMode {
		return
	}
	parser.panicMode = true
	fmt.Fprintf(os.Stderr, "[line %d] Error", token.Line)

	if token.Type == vmscanner.TokenEOF {
		fmt.Fprintf(os.Stderr, " at end")
	} else if token.Type == vmscanner.TokenError {
		// Nothing.
	} else {
		fmt.Fprintf(os.Stderr, " at '%s'", token.Lexeme())
	}

	fmt.Fprintf(os.Stderr, ": %s\n", message)
	parser.hadError = true
}

func currentChunk() *vmchunk.Chunk {
	return compilingChunk
}

func emitCode1(op vmchunk.OpCode) {
	currentChunk().WriteOpcode(op, parser.previous.Line)
}

func emitCode2(op vmchunk.OpCode, b byte) {
	currentChunk().WriteOpcode(op, parser.previous.Line)
	currentChunk().Write(b, parser.previous.Line)
}

func emitConstant(v vmvalue.Value) {
	emitCode2(vmchunk.OpConstant, makeConstant(v))
}

func makeConstant(v vmvalue.Value) byte {
	constant := currentChunk().AddConstant(v)
	if constant > math.MaxUint8 {
		errorAtPrev("Too many constants in one chunk.")
		return 0
	}
	return byte(constant)
}

func emitReturn() {
	emitCode1(vmchunk.OpReturn)
}

func endCompiler() {
	emitReturn()
}
