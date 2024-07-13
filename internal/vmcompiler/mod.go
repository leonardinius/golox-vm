package vmcompiler

import (
	"fmt"
	"os"

	"github.com/leonardinius/goloxvm/internal/vmchunk"
	"github.com/leonardinius/goloxvm/internal/vmscanner"
)

var (
	scanner vmscanner.Scanner
	parser  Parser
)

func Compile(source []byte, chunk *vmchunk.Chunk) bool {
	scanner = vmscanner.NewScanner(source)
	defer scanner.Free()

	parser = NewParser()

	advance()
	expression()
	consume(vmscanner.TokenEOF, "Expect end of expression.")

	constant1 := chunk.AddConstant(1.1)
	chunk.WriteOpcode(vmchunk.OpConstant, 1)
	chunk.Write(byte(constant1), 1)
	chunk.WriteOpcode(vmchunk.OpReturn, 1)

	return !parser.hadError
}

func advance() {
	parser.previous = parser.current

	for {
		parser.current = scanner.ScanToken()
		if parser.current.Type != vmscanner.TokenError {
			break
		}

		errorAtCurrent(parser.current.Lexeme())
	}
}

func expression() {
}

func consume(tk vmscanner.TokenType, message string) {
}

func errorAtCurrent(message string) {
	errorAt(&parser.current, message)
}

// func errorWith(message string) {
// 	errorAt(&parser.previous, message)
// }

func errorAt(token *vmscanner.Token, message string) {
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
