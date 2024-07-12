package vmcompiler

type Token struct {
	Type   TokenType
	Source []byte
	Start  int
	Length int
	Line   int
}

func MakeToken(scanner *Scanner, token TokenType) Token {
	return Token{
		Type:   token,
		Source: scanner.source,
		Start:  scanner.start,
		Length: scanner.current - scanner.start,
		Line:   scanner.line,
	}
}

func MakeErrorToken(scanner *Scanner, message string) Token {
	bytes := []byte(message)
	return Token{
		Type:   TokenError,
		Source: bytes,
		Start:  0,
		Length: len(bytes),
		Line:   scanner.line,
	}
}

func (t *Token) Lexeme() string {
	return string(t.Source[t.Start : t.Start+t.Length])
}
