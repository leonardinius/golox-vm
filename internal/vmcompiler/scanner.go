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

func (s *Scanner) ScanToken() Token {
	s.skipWhitespace()

	s.start = s.current
	if s.isAtEnd() {
		return s.makeToken(TokenEOF)
	}

	switch c := s.advance(); c {
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
	return s.source[s.current]
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
		c := s.peek()
		if c != ' ' && c != '\r' && c != '\t' {
			break
		}
		s.advance()
	}
}

func (s *Scanner) makeToken(token TokenType) Token {
	return MakeToken(s, token)
}

func (s *Scanner) errorToken(message string) Token {
	return MakeErrorToken(s, message)
}
