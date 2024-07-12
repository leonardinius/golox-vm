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
	TokenError
	TokenEOF
)

var tokenTypeStrings = map[TokenType]string{
	TokenLeftParen:    "TokenLeftParen",
	TokenRightParen:   "TokenRightParen",
	TokenLeftBrace:    "TokenLeftBrace",
	TokenRightBrace:   "TokenRightBrace",
	TokenComma:        "TokenComma",
	TokenDot:          "TokenDot",
	TokenMinus:        "TokenMinus",
	TokenPlus:         "TokenPlus",
	TokenSemicolon:    "TokenSemicolon",
	TokenSlash:        "TokenSlash",
	TokenStar:         "TokenStar",
	TokenBang:         "TokenBang",
	TokenBangEqual:    "TokenBangEqual",
	TokenEqual:        "TokenEqual",
	TokenEqualEqual:   "TokenEqualEqual",
	TokenGreater:      "TokenGreater",
	TokenGreaterEqual: "TokenGreaterEqual",
	TokenLess:         "TokenLess",
	TokenLessEqual:    "TokenLessEqual",
	TokenIdentifier:   "TokenIdentifier",
	TokenString:       "TokenString",
	TokenNumber:       "TokenNumber",
	TokenAnd:          "TokenAnd",
	TokenClass:        "TokenClass",
	TokenElse:         "TokenElse",
	TokenFalse:        "TokenFalse",
	TokenFor:          "TokenFor",
	TokenFun:          "TokenFun",
	TokenIf:           "TokenIf",
	TokenNil:          "TokenNil",
	TokenOr:           "TokenOr",
	TokenPrint:        "TokenPrint",
	TokenReturn:       "TokenReturn",
	TokenSuper:        "TokenSuper",
	TokenThis:         "TokenThis",
	TokenTrue:         "TokenTrue",
	TokenVar:          "TokenVar",
	TokenWhile:        "TokenWhile",
	TokenError:        "TokenError",
	TokenEOF:          "TokenEOF",
}

func (t TokenType) String() string {
	if str, ok := tokenTypeStrings[t]; ok {
		return str
	}

	return "Unknown"
}
