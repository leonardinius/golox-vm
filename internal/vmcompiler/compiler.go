package vmcompiler

import (
	"github.com/leonardinius/goloxvm/internal/vmchunk"
)

func Compile(source []byte, chunk *vmchunk.Chunk) error {
	scanner := NewScanner(source)
	defer scanner.Free()

	advance()
	expression()
	consume(TokenEOF, "Expect end of expression.")

	return nil
}

func advance() {
	panic("unimplemented")
}

func expression() {
	panic("unimplemented")
}

func consume(tk TokenType, message string) {
	panic("unimplemented")
}
