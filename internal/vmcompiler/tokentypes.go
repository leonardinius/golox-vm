package vmcompiler

type TokenType int

const (
	_ TokenType = iota
	// Single-character tokens.
	TokenLeftParen
	TokenRightParen
	TokenLeftBrace
	TokenRightBrace
	TokenComma
	TokenDot
	TokenMinus
	TokenPlus
	TokenSemicolon
	TokenSlash
	TokenStar

	// One or two character tokens.
	TokenBang
	TokenBangEqual
	TokenEqual
	TokenEqualEqual
	TokenGreater
	TokenGreaterEqual
	TokenLess
	TokenLessEqual

	// Literals.
	TokenIdentifier
	TokenString
	TokenNumber

	// Keywords.
	TokenAnd
	TokenClass
	TokenElse
	TokenFalse
	TokenFor
	TokenFun
	TokenIf
	TokenNil
	TokenOr
	TokenPrint
	TokenReturn
	TokenSuper
	TokenThis
	TokenTrue
	TokenVar
	TokenWhile

	// Special control tokens.
	TokenError
	TokenEOF
)

var tokenTypeStrings = map[TokenType]string{
	TokenLeftParen:    "TOKEN_LEFT_PAREN",
	TokenRightParen:   "TOKEN_RIGHT_PAREN",
	TokenLeftBrace:    "TOKEN_LEFT_BRACE",
	TokenRightBrace:   "TOKEN_RIGHT_BRACE",
	TokenComma:        "TOKEN_COMMA",
	TokenDot:          "TOKEN_DOT",
	TokenMinus:        "TOKEN_MINUS",
	TokenPlus:         "TOKEN_PLUS",
	TokenSemicolon:    "TOKEN_SEMICOLON",
	TokenSlash:        "TOKEN_SLASH",
	TokenStar:         "TOKEN_STAR",
	TokenBang:         "TOKEN_BANG",
	TokenBangEqual:    "TOKEN_BANG_EQUAL",
	TokenEqual:        "TOKEN_EQUAL",
	TokenEqualEqual:   "TOKEN_EQUAL_EQUAL",
	TokenGreater:      "TOKEN_GREATER",
	TokenGreaterEqual: "TOKEN_GREATER_EQUAL",
	TokenLess:         "TOKEN_LESS",
	TokenLessEqual:    "TOKEN_LESS_EQUAL",
	TokenIdentifier:   "TOKEN_IDENTIFIER",
	TokenString:       "TOKEN_STRING",
	TokenNumber:       "TOKEN_NUMBER",
	TokenAnd:          "TOKEN_AND",
	TokenClass:        "TOKEN_CLASS",
	TokenElse:         "TOKEN_ELSE",
	TokenFalse:        "TOKEN_FALSE",
	TokenFor:          "TOKEN_FOR",
	TokenFun:          "TOKEN_FUN",
	TokenIf:           "TOKEN_IF",
	TokenNil:          "TOKEN_NIL",
	TokenOr:           "TOKEN_OR",
	TokenPrint:        "TOKEN_PRINT",
	TokenReturn:       "TOKEN_RETURN",
	TokenSuper:        "TOKEN_SUPER",
	TokenThis:         "TOKEN_THIS",
	TokenTrue:         "TOKEN_TRUE",
	TokenVar:          "TOKEN_VAR",
	TokenWhile:        "TOKEN_WHILE",
	TokenError:        "TOKEN_ERROR",
	TokenEOF:          "TOKEN_EOF",
}

func (t TokenType) String() string {
	if str, ok := tokenTypeStrings[t]; ok {
		return str
	}

	return "Unknown"
}
