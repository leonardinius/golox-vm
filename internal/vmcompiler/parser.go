package vmcompiler

import (
	"fmt"
	"os"
	"strconv"

	"github.com/leonardinius/goloxvm/internal/vm/bytecode"
	"github.com/leonardinius/goloxvm/internal/vm/vmobject"
	"github.com/leonardinius/goloxvm/internal/vm/vmvalue"
	"github.com/leonardinius/goloxvm/internal/vmcompiler/scanner"
	"github.com/leonardinius/goloxvm/internal/vmcompiler/tokens"
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
	emitConstant(vmvalue.NumberAsValue(v))
}

func string_() {
	t := gParser.previous
	bytes := t.Source[t.Start+1 : t.Start+t.Length-1]
	emitConstant(vmvalue.ObjAsValue(vmobject.NewCopyString(bytes)))
}

func grouping() {
	expression()
	consume(tokens.TokenRightParen, "Expect ')' after expression.")
}

func literal() {
	switch literalType := gParser.previous.Type; literalType {
	case tokens.TokenFalse:
		emitOpcode(bytecode.OpFalse)
	case tokens.TokenNil:
		emitOpcode(bytecode.OpNil)
	case tokens.TokenTrue:
		emitOpcode(bytecode.OpTrue)
	default:
		panic(fmt.Sprintf("unexpected literal type: %s (%d)", literalType, literalType))
	}
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
	case tokens.TokenBangEqual:
		emitOpcodes(bytecode.OpEqual, bytecode.OpNot)
	case tokens.TokenEqualEqual:
		emitOpcode(bytecode.OpEqual)
	case tokens.TokenGreater:
		emitOpcode(bytecode.OpGreater)
	case tokens.TokenGreaterEqual:
		emitOpcodes(bytecode.OpLess, bytecode.OpNot)
	case tokens.TokenLess:
		emitOpcode(bytecode.OpLess)
	case tokens.TokenLessEqual:
		emitOpcodes(bytecode.OpGreater, bytecode.OpNot)
	case tokens.TokenPlus:
		emitOpcode(bytecode.OpAdd)
	case tokens.TokenMinus:
		emitOpcode(bytecode.OpSubtract)
	case tokens.TokenStar:
		emitOpcode(bytecode.OpMultiply)
	case tokens.TokenSlash:
		emitOpcode(bytecode.OpDivide)
	default:
		panic(fmt.Sprintf("unreachable operator: %s (%d)", operatorType, operatorType))
	}
}

func unary() {
	operatorType := gParser.previous.Type
	parsePrecedence(PrecedenceUnary)

	// emit the operator instruction
	switch operatorType {
	case tokens.TokenBang:
		emitOpcode(bytecode.OpNot)
	case tokens.TokenMinus:
		emitOpcode(bytecode.OpNegate)
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
		tokens.TokenBang:         {unary, nil, PrecedenceNone},
		tokens.TokenBangEqual:    {nil, binary, PrecedenceEquality},
		tokens.TokenEqual:        {nil, nil, PrecedenceNone},
		tokens.TokenEqualEqual:   {nil, binary, PrecedenceEquality},
		tokens.TokenGreater:      {nil, binary, PrecedenceComparison},
		tokens.TokenGreaterEqual: {nil, binary, PrecedenceComparison},
		tokens.TokenLess:         {nil, binary, PrecedenceComparison},
		tokens.TokenLessEqual:    {nil, binary, PrecedenceComparison},
		tokens.TokenIdentifier:   {nil, nil, PrecedenceNone},
		tokens.TokenString:       {string_, nil, PrecedenceNone},
		tokens.TokenNumber:       {number, nil, PrecedenceNone},
		tokens.TokenAnd:          {nil, nil, PrecedenceNone},
		tokens.TokenClass:        {nil, nil, PrecedenceNone},
		tokens.TokenElse:         {nil, nil, PrecedenceNone},
		tokens.TokenFalse:        {literal, nil, PrecedenceNone},
		tokens.TokenFor:          {nil, nil, PrecedenceNone},
		tokens.TokenFun:          {nil, nil, PrecedenceNone},
		tokens.TokenIf:           {nil, nil, PrecedenceNone},
		tokens.TokenNil:          {literal, nil, PrecedenceNone},
		tokens.TokenOr:           {nil, nil, PrecedenceNone},
		tokens.TokenPrint:        {nil, nil, PrecedenceNone},
		tokens.TokenReturn:       {nil, nil, PrecedenceNone},
		tokens.TokenSuper:        {nil, nil, PrecedenceNone},
		tokens.TokenThis:         {nil, nil, PrecedenceNone},
		tokens.TokenTrue:         {literal, nil, PrecedenceNone},
		tokens.TokenVar:          {nil, nil, PrecedenceNone},
		tokens.TokenWhile:        {nil, nil, PrecedenceNone},
		tokens.TokenError:        {nil, nil, PrecedenceNone},
		tokens.TokenEOF:          {nil, nil, PrecedenceNone},
	}
}
