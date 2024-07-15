package vmcompiler

import (
	"fmt"
	"os"
	"strconv"

	"github.com/leonardinius/goloxvm/internal/vmchunk"
	"github.com/leonardinius/goloxvm/internal/vmscanner"
	"github.com/leonardinius/goloxvm/internal/vmvalue"
)

var rules map[vmscanner.TokenType]*ParseRule

type Parser struct {
	current   vmscanner.Token
	previous  vmscanner.Token
	hadError  bool
	panicMode bool
}

func NewParser() Parser {
	return Parser{
		hadError:  false,
		panicMode: false,
	}
}

type ParsePrecedence int

const (
	_ ParsePrecedence = iota
	PrecedenceNone
	PrecedenceAssignment // =
	PrecedenceOr         // or
	PrecedenceAnd        // and
	PrecedenceEquality   // == !=
	PrecedenceComparison // < > <= >=
	PrecedenceTerm       // + -
	PrecedenceFactor     // * /
	PrecedenceUnary      // ! -
	PrecedenceCall       // . ()
	PrecedencePrimary
)

func (p ParsePrecedence) Inc() ParsePrecedence {
	return p + 1
}

func (p ParsePrecedence) Dec() ParsePrecedence {
	return p - 1
}

type (
	ParseFn      func()
	InfixParseFn func()
)

type ParseRule struct {
	prefixRule ParseFn
	infixRule  InfixParseFn
	precedence ParsePrecedence
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

func parsePrecedence(precedence ParsePrecedence) {
	advance()

	prefixRule := mustGetRule(parser.previous.Type).prefixRule
	if prefixRule == nil {
		errorAtPrev("Expect expression.")
		return
	}
	prefixRule()

	for precedence <= mustGetRule(parser.current.Type).precedence {
		advance()
		infixRule := mustGetRule(parser.previous.Type).infixRule
		infixRule()
	}
}

func expression() {
	parsePrecedence(PrecedenceAssignment)
}

func number() {
	v, err := strconv.ParseFloat(parser.previous.Lexeme(), 64)
	if err != nil {
		errorAtPrev(err.Error())
	}
	emitConstant(vmvalue.Value(v))
}

func grouping() {
	expression()
	consume(vmscanner.TokenRightParen, "Expect ')' after expression.")
}

func binary() {
	operatorType := parser.previous.Type
	rule := mustGetRule(operatorType)
	parsePrecedence(rule.precedence.Inc())

	switch operatorType {
	case vmscanner.TokenPlus:
		emitCode1(vmchunk.OpAdd)
		break
	case vmscanner.TokenMinus:
		emitCode1(vmchunk.OpSubtract)
		break
	case vmscanner.TokenStar:
		emitCode1(vmchunk.OpMultiply)
		break
	case vmscanner.TokenSlash:
		emitCode1(vmchunk.OpDivide)
		break
	default:
		return // Unreachable.
	}
}

func unary() {
	parsePrecedence(PrecedenceUnary)

	// compile the operand
	expression()

	// emit the operator instruction
	switch parser.previous.Type {
	case vmscanner.TokenMinus:
		emitCode1(vmchunk.OpNegate)
	default:
		panic("Unreachable unary: " + parser.previous.Lexeme())
	}
}

func mustGetRule(t vmscanner.TokenType) *ParseRule {
	if r, ok := rules[t]; ok {
		return r
	} else {
		panic(fmt.Sprintf("get rule %s (%d)", t, t))
	}
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

func init() {
	rules = map[vmscanner.TokenType]*ParseRule{
		vmscanner.TokenLeftParen:  {grouping, nil, PrecedenceNone},
		vmscanner.TokenRightParen: {nil, nil, PrecedenceNone},
		vmscanner.TokenLeftBrace:  {nil, nil, PrecedenceNone},
		vmscanner.TokenRightBrace: {nil, nil, PrecedenceNone},
		vmscanner.TokenComma:      {nil, nil, PrecedenceNone},
		vmscanner.TokenDot:        {nil, nil, PrecedenceNone},
		vmscanner.TokenMinus:      {unary, binary, PrecedenceTerm},
		vmscanner.TokenPlus:       {nil, binary, PrecedenceTerm},
		//	[TOKEN_SEMICOLON]     = {NULL,     NULL,   PREC_NONE},
		vmscanner.TokenSemicolon: {nil, nil, PrecedenceNone},
		//	[TOKEN_SLASH]         = {NULL,     binary, PREC_FACTOR},
		vmscanner.TokenSlash: {nil, binary, PrecedenceFactor},
		//	[TOKEN_STAR]          = {NULL,     binary, PREC_FACTOR},

		//	[TOKEN_BANG]          = {NULL,     NULL,   PREC_NONE},
		//	[TOKEN_BANG_EQUAL]    = {NULL,     NULL,   PREC_NONE},
		//	[TOKEN_EQUAL]         = {NULL,     NULL,   PREC_NONE},
		//	[TOKEN_EQUAL_EQUAL]   = {NULL,     NULL,   PREC_NONE},
		//	[TOKEN_GREATER]       = {NULL,     NULL,   PREC_NONE},
		//	[TOKEN_GREATER_EQUAL] = {NULL,     NULL,   PREC_NONE},
		//	[TOKEN_LESS]          = {NULL,     NULL,   PREC_NONE},
		//	[TOKEN_LESS_EQUAL]    = {NULL,     NULL,   PREC_NONE},
		//	[TOKEN_IDENTIFIER]    = {NULL,     NULL,   PREC_NONE},
		//	[TOKEN_STRING]        = {NULL,     NULL,   PREC_NONE},
		//	[TOKEN_NUMBER]        = {number,   NULL,   PREC_NONE},
		//	[TOKEN_AND]           = {NULL,     NULL,   PREC_NONE},
		//	[TOKEN_CLASS]         = {NULL,     NULL,   PREC_NONE},
		//	[TOKEN_ELSE]          = {NULL,     NULL,   PREC_NONE},
		//	[TOKEN_FALSE]         = {NULL,     NULL,   PREC_NONE},
		//
		// [TOKEN_FOR]           = {NULL,     NULL,   PREC_NONE},
		// [TOKEN_FUN]           = {NULL,     NULL,   PREC_NONE},
		// [TOKEN_IF]            = {NULL,     NULL,   PREC_NONE},
		// [TOKEN_NIL]           = {NULL,     NULL,   PREC_NONE},
		// [TOKEN_OR]            = {NULL,     NULL,   PREC_NONE},
		// [TOKEN_PRINT]         = {NULL,     NULL,   PREC_NONE},
		// [TOKEN_RETURN]        = {NULL,     NULL,   PREC_NONE},
		// [TOKEN_SUPER]         = {NULL,     NULL,   PREC_NONE},
		// [TOKEN_THIS]          = {NULL,     NULL,   PREC_NONE},
		// [TOKEN_TRUE]          = {NULL,     NULL,   PREC_NONE},
		// [TOKEN_VAR]           = {NULL,     NULL,   PREC_NONE},
		// [TOKEN_WHILE]         = {NULL,     NULL,   PREC_NONE},
		// [TOKEN_ERROR]         = {NULL,     NULL,   PREC_NONE},
		// [TOKEN_EOF]           = {NULL,     NULL,   PREC_NONE},
	}
}
