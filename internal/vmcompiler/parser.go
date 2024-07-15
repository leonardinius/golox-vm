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
	case vmscanner.TokenMinus:
		emitCode1(vmchunk.OpSubtract)
	case vmscanner.TokenStar:
		emitCode1(vmchunk.OpMultiply)
	case vmscanner.TokenSlash:
		emitCode1(vmchunk.OpDivide)
	default:
		panic(fmt.Sprintf("Unreachable operator: %s (%d)", operatorType, operatorType))
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
		vmscanner.TokenLeftParen:    {grouping, nil, PrecedenceNone},
		vmscanner.TokenRightParen:   {nil, nil, PrecedenceNone},
		vmscanner.TokenLeftBrace:    {nil, nil, PrecedenceNone},
		vmscanner.TokenRightBrace:   {nil, nil, PrecedenceNone},
		vmscanner.TokenComma:        {nil, nil, PrecedenceNone},
		vmscanner.TokenDot:          {nil, nil, PrecedenceNone},
		vmscanner.TokenMinus:        {unary, binary, PrecedenceTerm},
		vmscanner.TokenPlus:         {nil, binary, PrecedenceTerm},
		vmscanner.TokenSemicolon:    {nil, nil, PrecedenceNone},
		vmscanner.TokenSlash:        {nil, binary, PrecedenceFactor},
		vmscanner.TokenStar:         {nil, binary, PrecedenceFactor},
		vmscanner.TokenBang:         {nil, nil, PrecedenceNone},
		vmscanner.TokenBangEqual:    {nil, nil, PrecedenceNone},
		vmscanner.TokenEqual:        {nil, nil, PrecedenceNone},
		vmscanner.TokenEqualEqual:   {nil, nil, PrecedenceNone},
		vmscanner.TokenGreater:      {nil, nil, PrecedenceNone},
		vmscanner.TokenGreaterEqual: {nil, nil, PrecedenceNone},
		vmscanner.TokenLess:         {nil, nil, PrecedenceNone},
		vmscanner.TokenLessEqual:    {nil, nil, PrecedenceNone},
		vmscanner.TokenIdentifier:   {nil, nil, PrecedenceNone},
		vmscanner.TokenString:       {nil, nil, PrecedenceNone},
		vmscanner.TokenNumber:       {number, nil, PrecedenceNone},
		vmscanner.TokenAnd:          {nil, nil, PrecedenceNone},
		vmscanner.TokenClass:        {nil, nil, PrecedenceNone},
		vmscanner.TokenElse:         {nil, nil, PrecedenceNone},
		vmscanner.TokenFalse:        {nil, nil, PrecedenceNone},
		vmscanner.TokenFor:          {nil, nil, PrecedenceNone},
		vmscanner.TokenFun:          {nil, nil, PrecedenceNone},
		vmscanner.TokenIf:           {nil, nil, PrecedenceNone},
		vmscanner.TokenNil:          {nil, nil, PrecedenceNone},
		vmscanner.TokenOr:           {nil, nil, PrecedenceNone},
		vmscanner.TokenPrint:        {nil, nil, PrecedenceNone},
		vmscanner.TokenReturn:       {nil, nil, PrecedenceNone},
		vmscanner.TokenSuper:        {nil, nil, PrecedenceNone},
		vmscanner.TokenThis:         {nil, nil, PrecedenceNone},
		vmscanner.TokenTrue:         {nil, nil, PrecedenceNone},
		vmscanner.TokenVar:          {nil, nil, PrecedenceNone},
		vmscanner.TokenWhile:        {nil, nil, PrecedenceNone},
		vmscanner.TokenError:        {nil, nil, PrecedenceNone},
		vmscanner.TokenEOF:          {nil, nil, PrecedenceNone},
	}
}
