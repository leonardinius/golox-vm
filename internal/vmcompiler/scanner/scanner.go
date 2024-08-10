package scanner

import "github.com/leonardinius/goloxvm/internal/vmcompiler/tokens"

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
		return s.makeToken(tokens.TokenEOF)
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
		return s.makeToken(tokens.TokenLeftParen)
	case ')':
		return s.makeToken(tokens.TokenRightParen)
	case '{':
		return s.makeToken(tokens.TokenLeftBrace)
	case '}':
		return s.makeToken(tokens.TokenRightBrace)
	case ';':
		return s.makeToken(tokens.TokenSemicolon)
	case ',':
		return s.makeToken(tokens.TokenComma)
	case '.':
		return s.makeToken(tokens.TokenDot)
	case '-':
		return s.makeToken(tokens.TokenMinus)
	case '+':
		return s.makeToken(tokens.TokenPlus)
	case '/':
		return s.makeToken(tokens.TokenSlash)
	case '*':
		return s.makeToken(tokens.TokenStar)
	case '!':
		if s.match('=') {
			return s.makeToken(tokens.TokenBangEqual)
		}
		return s.makeToken(tokens.TokenBang)
	case '=':
		if s.match('=') {
			return s.makeToken(tokens.TokenEqualEqual)
		}
		return s.makeToken(tokens.TokenEqual)
	case '<':
		if s.match('=') {
			return s.makeToken(tokens.TokenLessEqual)
		}
		return s.makeToken(tokens.TokenLess)
	case '>':
		if s.match('=') {
			return s.makeToken(tokens.TokenGreaterEqual)
		}
		return s.makeToken(tokens.TokenGreater)
	case '"':
		return s.string()
	}

	return s.errorToken("Unexpected character.")
}

func (s *Scanner) isAtEndPeek(peek int) bool {
	return s.current+peek >= len(s.source)
}

func (s *Scanner) isAtEnd() bool {
	return s.isAtEndPeek(0)
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
	if s.isAtEndPeek(1) {
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
	token := s.makeToken(tokens.TokenString)
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

	return s.makeToken(tokens.TokenNumber)
}

func (s *Scanner) isAlpha(c byte) bool {
	return c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func (s *Scanner) checkKeyword(start, length int, rest string, token tokens.TokenType) tokens.TokenType {
	if s.current-s.start != start+length {
		return tokens.TokenIdentifier
	}

	if rest != string(s.source[s.start+start:s.start+start+length]) {
		return tokens.TokenIdentifier
	}

	return token
}

func (s *Scanner) identifierType() tokens.TokenType {
	switch s.source[s.start] {
	case 'a':
		return s.checkKeyword(1, 2, "nd", tokens.TokenAnd)
	case 'c':
		return s.checkKeyword(1, 4, "lass", tokens.TokenClass)
	case 'e':
		return s.checkKeyword(1, 3, "lse", tokens.TokenElse)
	case 'f': // for, fun
		if s.current-s.start > 1 {
			switch s.source[s.start+1] {
			case 'a':
				return s.checkKeyword(2, 3, "lse", tokens.TokenFalse)
			case 'o':
				return s.checkKeyword(2, 1, "r", tokens.TokenFor)
			case 'u':
				return s.checkKeyword(2, 1, "n", tokens.TokenFun)
			}
		}
	case 'i':
		return s.checkKeyword(1, 1, "f", tokens.TokenIf)
	case 'n':
		return s.checkKeyword(1, 2, "il", tokens.TokenNil)
	case 'o':
		return s.checkKeyword(1, 1, "r", tokens.TokenOr)
	case 'p':
		return s.checkKeyword(1, 4, "rint", tokens.TokenPrint)
	case 'r':
		return s.checkKeyword(1, 5, "eturn", tokens.TokenReturn)
	case 's':
		return s.checkKeyword(1, 4, "uper", tokens.TokenSuper)
	case 't': // this, true
		if s.current-s.start > 1 {
			switch s.source[s.start+1] {
			case 'h':
				return s.checkKeyword(2, 2, "is", tokens.TokenThis)
			case 'r':
				return s.checkKeyword(2, 2, "ue", tokens.TokenTrue)
			}
		}
	case 'v':
		return s.checkKeyword(1, 2, "ar", tokens.TokenVar)
	case 'w':
		return s.checkKeyword(1, 4, "hile", tokens.TokenWhile)
	}

	return tokens.TokenIdentifier
}

func (s *Scanner) identifier() Token {
	for s.isAlpha(s.peek()) || s.isDigit(s.peek()) {
		s.advance()
	}

	return s.makeToken(s.identifierType())
}

func (s *Scanner) makeToken(tokenType tokens.TokenType) Token {
	return MakeToken(s, tokenType)
}

func (s *Scanner) errorToken(message string) Token {
	return MakeErrorToken(s, message)
}
