package vmcompiler

type Scanner struct {
	source  []byte
	start   int
	current int
	line    int
}

func NewScanner(source []byte) Scanner {
	return Scanner{
		source:  source,
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) Free() {
	s.source = nil
}

func (s *Scanner) ScanToken() Token {
	s.skipWhitespace()

	s.start = s.current
	if s.isAtEnd() {
		return s.makeToken(TokenEOF)
	}

	c := s.advance()
	if s.isAlpha(c) {
		return s.identifier()
	}
	if s.isDigit(c) {
		return s.number()
	}

	switch c {
	case '(':
		return s.makeToken(TokenLeftParen)
	case ')':
		return s.makeToken(TokenRightParen)
	case '{':
		return s.makeToken(TokenLeftBrace)
	case '}':
		return s.makeToken(TokenRightBrace)
	case ';':
		return s.makeToken(TokenSemicolon)
	case ',':
		return s.makeToken(TokenComma)
	case '.':
		return s.makeToken(TokenDot)
	case '-':
		return s.makeToken(TokenMinus)
	case '+':
		return s.makeToken(TokenPlus)
	case '/':
		return s.makeToken(TokenSlash)
	case '*':
		return s.makeToken(TokenStar)
	case '!':
		if s.match('=') {
			return s.makeToken(TokenBangEqual)
		}
		return s.makeToken(TokenBang)
	case '=':
		if s.match('=') {
			return s.makeToken(TokenEqualEqual)
		}
		return s.makeToken(TokenEqual)
	case '<':
		if s.match('=') {
			return s.makeToken(TokenLessEqual)
		}
		return s.makeToken(TokenLess)
	case '>':
		if s.match('=') {
			return s.makeToken(TokenGreaterEqual)
		}
		return s.makeToken(TokenGreater)
	case '"':
		return s.string()
	}

	return s.errorToken("Unexpected character.")
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() byte {
	s.current++
	return s.source[s.current-1]
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current+1]
}

func (s *Scanner) match(c byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != c {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) skipWhitespace() {
	for !s.isAtEnd() {
		switch c := s.peek(); c {
		case ' ', '\r', '\t':
			s.advance()
		case '\n':
			s.line++
			s.advance()
		case '/':
			if s.peekNext() == '/' {
				for !s.isAtEnd() && s.peek() != '\n' {
					s.advance()
				}
				break
			}
			return
		default:
			return
		}
	}
}

func (s *Scanner) string() Token {
	startLine := s.line
	for !s.isAtEnd() && s.peek() != '"' {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		return s.errorToken("Unterminated string.")
	}

	// Consume the closing quote.
	s.advance()
	token := s.makeToken(TokenString)
	// otherwise it will report the token to start ending "
	// if multiline string - it is wrong
	token.Line = startLine
	return token
}

func (s *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) number() Token {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	// Look for a fractional part.
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		// Consume the "."
		s.advance()
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	return s.makeToken(TokenNumber)
}

func (s *Scanner) isAlpha(c byte) bool {
	return c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func (s *Scanner) checkKeyword(start, length int, rest string, token TokenType) TokenType {
	if s.current-s.start != start+length {
		return TokenIdentifier
	}

	if rest != string(s.source[s.start+start:s.start+start+length]) {
		return TokenIdentifier
	}

	return token
}

func (s *Scanner) identifierType() TokenType {
	switch s.source[s.start] {
	case 'a':
		return s.checkKeyword(1, 2, "nd", TokenAnd)
	case 'c':
		return s.checkKeyword(1, 4, "lass", TokenClass)
	case 'e':
		return s.checkKeyword(1, 3, "lse", TokenElse)
	case 'f': // for, fun
		if s.current-s.start > 1 {
			switch s.source[s.start+1] {
			case 'a':
				return s.checkKeyword(2, 3, "lse", TokenFalse)
			case 'o':
				return s.checkKeyword(2, 1, "r", TokenFor)
			case 'u':
				return s.checkKeyword(2, 1, "n", TokenFun)
			}
		}
	case 'i':
		return s.checkKeyword(1, 1, "f", TokenIf)
	case 'n':
		return s.checkKeyword(1, 2, "il", TokenNil)
	case 'o':
		return s.checkKeyword(1, 1, "r", TokenOr)
	case 'p':
		return s.checkKeyword(1, 4, "rint", TokenPrint)
	case 'r':
		return s.checkKeyword(1, 5, "eturn", TokenReturn)
	case 's':
		return s.checkKeyword(1, 4, "uper", TokenSuper)
	case 't': // this, true
		if s.current-s.start > 1 {
			switch s.source[s.start+1] {
			case 'h':
				return s.checkKeyword(2, 2, "is", TokenThis)
			case 'r':
				return s.checkKeyword(2, 2, "ue", TokenTrue)
			}
		}
	case 'v':
		return s.checkKeyword(1, 2, "ar", TokenVar)
	case 'w':
		return s.checkKeyword(1, 4, "hile", TokenWhile)
	}

	return TokenIdentifier
}

func (s *Scanner) identifier() Token {
	for s.isAlpha(s.peek()) || s.isDigit(s.peek()) {
		s.advance()
	}

	return s.makeToken(s.identifierType())
}

func (s *Scanner) makeToken(token TokenType) Token {
	return MakeToken(s, token)
}

func (s *Scanner) errorToken(message string) Token {
	return MakeErrorToken(s, message)
}
