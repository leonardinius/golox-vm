package vmcompiler

import (
	"github.com/leonardinius/goloxvm/internal/vmchunk"
	"github.com/leonardinius/goloxvm/internal/vmscanner"
)

func Compile(source []byte, chunk *vmchunk.Chunk) error {
	scanner := vmscanner.NewScanner(source)
	defer scanner.Free()

	advance()
	expression()
	consume(vmscanner.TokenEOF, "Expect end of expression.")

	return nil
}

func advance() {
	panic("unimplemented")
}

func expression() {
	panic("unimplemented")
}

func consume(tk vmscanner.TokenType, message string) {
	panic("unimplemented")
}
