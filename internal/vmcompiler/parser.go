package vmcompiler

import (
	"fmt"
	"os"
	"strconv"

	"github.com/leonardinius/goloxvm/internal/bytecode"
	"github.com/leonardinius/goloxvm/internal/scanner"
	"github.com/leonardinius/goloxvm/internal/tokens"
	"github.com/leonardinius/goloxvm/internal/vmvalue"
)

type Parser struct {
	current   scanner.Token
	previous  scanner.Token
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

func (p ParsePrecedence) Next() ParsePrecedence {
	return p + 1
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
	gParser.previous = gParser.current

	for {
		gParser.current = gScanner.ScanToken()
		if gParser.current.Type != tokens.TokenError {
			break
		}
		// use TokenError lexeme as error message
		errorAtCurrent(gParser.current.Lexeme())
	}
}

func parsePrecedence(precedence ParsePrecedence) {
	advance()

	prefixRule := mustGetRule(gParser.previous.Type).prefixRule
	if prefixRule == nil {
		errorAtPrev("Expect expression.")
		return
	}
	prefixRule()

	for precedence <= mustGetRule(gParser.current.Type).precedence {
		advance()
		infixRule := mustGetRule(gParser.previous.Type).infixRule
		infixRule()
	}
}

func expression() {
	parsePrecedence(PrecedenceAssignment)
}

func number() {
	v, err := strconv.ParseFloat(gParser.previous.Lexeme(), 64)
	if err != nil {
		errorAtPrev(err.Error())
	}
	emitConstant(vmvalue.NumberValue(v))
}

func grouping() {
	expression()
	consume(tokens.TokenRightParen, "Expect ')' after expression.")
}

func binary() {
	// the 1st (left) operand has been already parsed and consumed by this point

	// operator type
	operatorType := gParser.previous.Type
	// rule for the operator
	rule := mustGetRule(operatorType)
	// parse the second (right) operand
	parsePrecedence(rule.precedence.Next())

	switch operatorType {
	case tokens.TokenPlus:
		emitByte(bytecode.OpAdd)
	case tokens.TokenMinus:
		emitByte(bytecode.OpSubtract)
	case tokens.TokenStar:
		emitByte(bytecode.OpMultiply)
	case tokens.TokenSlash:
		emitByte(bytecode.OpDivide)
	default:
		panic(fmt.Sprintf("Unreachable operator: %s (%d)", operatorType, operatorType))
	}
}

func unary() {
	operatorType := gParser.previous.Type
	parsePrecedence(PrecedenceUnary)

	// emit the operator instruction
	switch operatorType {
	case tokens.TokenMinus:
		emitByte(bytecode.OpNegate)
	default:
		panic("Unreachable unary: " + gParser.previous.Lexeme())
	}
}

func mustGetRule(t tokens.TokenType) *ParseRule {
	if r, ok := rules[t]; ok {
		return r
	} else {
		panic(fmt.Sprintf("get rule %s (%d)", t, t))
	}
}

func consume(stype tokens.TokenType, message string) {
	if gParser.current.Type == stype {
		advance()
		return
	}

	errorAtCurrent(message)
}

func errorAtCurrent(message string) {
	errorAt(&gParser.current, message)
}

func errorAtPrev(message string) {
	errorAt(&gParser.previous, message)
}

func errorAt(token *scanner.Token, message string) {
	if gParser.panicMode {
		return
	}
	gParser.panicMode = true
	fmt.Fprintf(os.Stderr, "[line %d] Error", token.Line)

	if token.Type == tokens.TokenEOF {
		fmt.Fprintf(os.Stderr, " at end")
	} else if token.Type == tokens.TokenError {
		// Nothing.
	} else {
		fmt.Fprintf(os.Stderr, " at '%s'", token.Lexeme())
	}

	fmt.Fprintf(os.Stderr, ": %s\n", message)
	gParser.hadError = true
}

var rules map[tokens.TokenType]*ParseRule

func init() {
	rules = map[tokens.TokenType]*ParseRule{
		tokens.TokenLeftParen:    {grouping, nil, PrecedenceNone},
		tokens.TokenRightParen:   {nil, nil, PrecedenceNone},
		tokens.TokenLeftBrace:    {nil, nil, PrecedenceNone},
		tokens.TokenRightBrace:   {nil, nil, PrecedenceNone},
		tokens.TokenComma:        {nil, nil, PrecedenceNone},
		tokens.TokenDot:          {nil, nil, PrecedenceNone},
		tokens.TokenMinus:        {unary, binary, PrecedenceTerm},
		tokens.TokenPlus:         {nil, binary, PrecedenceTerm},
		tokens.TokenSemicolon:    {nil, nil, PrecedenceNone},
		tokens.TokenSlash:        {nil, binary, PrecedenceFactor},
		tokens.TokenStar:         {nil, binary, PrecedenceFactor},
		tokens.TokenBang:         {nil, nil, PrecedenceNone},
		tokens.TokenBangEqual:    {nil, nil, PrecedenceNone},
		tokens.TokenEqual:        {nil, nil, PrecedenceNone},
		tokens.TokenEqualEqual:   {nil, nil, PrecedenceNone},
		tokens.TokenGreater:      {nil, nil, PrecedenceNone},
		tokens.TokenGreaterEqual: {nil, nil, PrecedenceNone},
		tokens.TokenLess:         {nil, nil, PrecedenceNone},
		tokens.TokenLessEqual:    {nil, nil, PrecedenceNone},
		tokens.TokenIdentifier:   {nil, nil, PrecedenceNone},
		tokens.TokenString:       {nil, nil, PrecedenceNone},
		tokens.TokenNumber:       {number, nil, PrecedenceNone},
		tokens.TokenAnd:          {nil, nil, PrecedenceNone},
		tokens.TokenClass:        {nil, nil, PrecedenceNone},
		tokens.TokenElse:         {nil, nil, PrecedenceNone},
		tokens.TokenFalse:        {nil, nil, PrecedenceNone},
		tokens.TokenFor:          {nil, nil, PrecedenceNone},
		tokens.TokenFun:          {nil, nil, PrecedenceNone},
		tokens.TokenIf:           {nil, nil, PrecedenceNone},
		tokens.TokenNil:          {nil, nil, PrecedenceNone},
		tokens.TokenOr:           {nil, nil, PrecedenceNone},
		tokens.TokenPrint:        {nil, nil, PrecedenceNone},
		tokens.TokenReturn:       {nil, nil, PrecedenceNone},
		tokens.TokenSuper:        {nil, nil, PrecedenceNone},
		tokens.TokenThis:         {nil, nil, PrecedenceNone},
		tokens.TokenTrue:         {nil, nil, PrecedenceNone},
		tokens.TokenVar:          {nil, nil, PrecedenceNone},
		tokens.TokenWhile:        {nil, nil, PrecedenceNone},
		tokens.TokenError:        {nil, nil, PrecedenceNone},
		tokens.TokenEOF:          {nil, nil, PrecedenceNone},
	}
}
