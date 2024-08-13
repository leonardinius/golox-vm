package vmcompiler

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	"github.com/leonardinius/goloxvm/internal/vm/bytecode"
	"github.com/leonardinius/goloxvm/internal/vm/hashtable"
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

func (p ParsePrecedence) CanAssign() bool {
	return p <= PrecedenceAssignment
}

type (
	ParseFn func(precedence ParsePrecedence)
)

type ParseRule struct {
	prefixRule ParseFn
	infixRule  ParseFn
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
		errorAtCurrent(gParser.current.LexemeAsString())
	}
}

func parsePrecedence(precedence ParsePrecedence) {
	advance()

	prefixRule := mustGetRule(gParser.previous.Type).prefixRule
	if prefixRule == nil {
		errorAtPrev("Expect expression.")
		return
	}
	prefixRule(precedence)

	for precedence <= mustGetRule(gParser.current.Type).precedence {
		advance()
		infixRule := mustGetRule(gParser.previous.Type).infixRule
		infixRule(precedence)
	}

	if precedence.CanAssign() && match(tokens.TokenEqual) {
		errorAtPrev("Invalid assignment target.")
	}
}

func identifierConstant(token *scanner.Token) int {
	identifier := hashtable.StringInternCopy(token.Lexeme())
	value := vmvalue.ObjAsValue(identifier)
	return makeConstant(value)
}

func resolveLocal(compiler *Compiler, name *scanner.Token) (slot int, ok bool) {
	for i := compiler.LocalCount - 1; i >= 0; i-- {
		local := &compiler.Locals[i]
		if bytes.Equal(name.Lexeme(), local.Name.Lexeme()) {
			if local.Depth == -1 {
				errorAtPrev("Can't read local variable in its own initializer.")
			}
			return i, true
		}
	}

	return 0, false
}

func addLocal(name scanner.Token) {
	if gCurrent.LocalCount == len(gCurrent.Locals) {
		errorAtPrev("Too many local variables in function.")
		return
	}

	local := &gCurrent.Locals[gCurrent.LocalCount]
	gCurrent.LocalCount++
	local.Name = name
	local.Depth = -1
	local.IsCaptured = false
}

func resolveUpvalue(compiler *Compiler, name *scanner.Token) (slot int, ok bool) {
	if compiler.Enclosing == nil {
		return 0, false
	}

	if local, ok := resolveLocal(compiler.Enclosing, name); ok {
		compiler.Enclosing.Locals[local].IsCaptured = true
		return addUpvalue(compiler, local, 1), true
	}

	if upvalue, ok := resolveUpvalue(compiler.Enclosing, name); ok {
		return addUpvalue(compiler, upvalue, 0), true
	}

	return 0, false
}

func addUpvalue(compiler *Compiler, index int, islocal byte) int {
	upvalueCount := compiler.Function.UpvalueCount

	for i := range upvalueCount {
		upvalue := &compiler.Upvalues[i]
		if upvalue.Index == index && upvalue.Local == islocal {
			return i
		}
	}

	if upvalueCount == MaxUpvalueCount {
		errorAtPrev("Too many closure variables in function.")
		return 0
	}

	upvalue := &compiler.Upvalues[upvalueCount]
	upvalue.Local = islocal
	upvalue.Index = index
	compiler.Function.UpvalueCount++
	return upvalueCount
}

func declareVariable() {
	if gCurrent.ScoreDepth == 0 {
		return
	}

	name := &gParser.previous
	// search for local variable
	for i := gCurrent.LocalCount - 1; i >= 0; i-- {
		local := &gCurrent.Locals[i]
		if local.Depth != -1 && local.Depth < gCurrent.ScoreDepth {
			break
		}

		if bytes.Equal(name.Lexeme(), local.Name.Lexeme()) {
			errorAtPrev("Already a variable with this name in this scope.")
		}
	}

	addLocal(*name)
}

func parseVariable(errorMessage string) int {
	consume(tokens.TokenIdentifier, errorMessage)

	declareVariable()
	if gCurrent.ScoreDepth > 0 {
		return 0
	}

	return identifierConstant(&gParser.previous)
}

func markInitialized() {
	if gCurrent.ScoreDepth == 0 {
		return
	}
	gCurrent.Locals[gCurrent.LocalCount-1].Depth = gCurrent.ScoreDepth
}

func defineVariable(global int) {
	if gCurrent.ScoreDepth > 0 {
		markInitialized()
		return
	}

	emitOpByte(bytecode.OpDefineGlobal, byte(global))
}

func and_(ParsePrecedence) {
	endJump := emitJump(bytecode.OpJumpIfFalse)

	emitOpcode(bytecode.OpPop)
	parsePrecedence(PrecedenceAnd)

	patchJump(endJump)
}

func or_(ParsePrecedence) {
	elseJump := emitJump(bytecode.OpJumpIfFalse)
	endJump := emitJump(bytecode.OpJump)

	patchJump(elseJump)
	emitOpcode(bytecode.OpPop)

	parsePrecedence(PrecedenceOr)
	patchJump(endJump)
}

func expression() {
	parsePrecedence(PrecedenceAssignment)
}

func block() {
	for !check(tokens.TokenRightBrace) && !check(tokens.TokenEOF) {
		declaration()
	}

	consume(tokens.TokenRightBrace, "Expect '}' after block.")
}

func function(fnType FunctionType, fnName *vmvalue.ObjString) {
	compiler := NewCompiler(fnType, fnName)
	beginScope()

	consume(tokens.TokenLeftParen, "Expect '(' after function name.")
	if !check(tokens.TokenRightParen) {
		for {
			gCurrent.Function.Arity++
			if gCurrent.Function.Arity > MaxArity {
				errorAtCurrent("Can't have more than 255 parameters.")
			}

			paramConstant := parseVariable("Expect parameter name.")
			defineVariable(paramConstant)

			if !match(tokens.TokenComma) {
				break
			}
		}
	}
	consume(tokens.TokenRightParen, "Expect ')' after parameters.")

	consume(tokens.TokenLeftBrace, "Expect '{' before function body.")
	block()

	// end of function
	fn := endCompiler()
	emitOpByte(bytecode.OpClosure, byte(makeConstant(vmvalue.ObjAsValue(fn))))
	for i := range fn.UpvalueCount {
		upvalue := &compiler.Upvalues[i]
		emitByte(upvalue.Local)
		emitByte(byte(upvalue.Index))
	}
}

func funDeclaration() {
	global := parseVariable("Expect function name.")
	markInitialized()
	function(FunctionTypeFunction, hashtable.StringInternTake(gParser.previous.Lexeme()))
	defineVariable(global)
}

func varDeclaration() {
	global := parseVariable("Expect variable name.")

	if match(tokens.TokenEqual) {
		expression()
	} else {
		emitOpcode(bytecode.OpNil)
	}
	consume(tokens.TokenSemicolon, "Expect ';' after variable declaration.")

	defineVariable(global)
}

func printStatement() {
	expression()
	consume(tokens.TokenSemicolon, "Expect ';' after value.")
	emitOpcode(bytecode.OpPrint)
}

func returnStatement() {
	if gCurrent.FnType == FunctionTypeScript {
		errorAtPrev("Can't return from top-level code.")
	}

	if match(tokens.TokenSemicolon) {
		emitReturn()
	} else {
		expression()
		consume(tokens.TokenSemicolon, "Expect ';' after return value.")
		emitOpcode(bytecode.OpReturn)
	}
}

func synchronize() {
	gParser.panicMode = false

	for gParser.current.Type != tokens.TokenEOF {
		if gParser.previous.Type == tokens.TokenSemicolon {
			return
		}

		switch gParser.current.Type {
		case tokens.TokenClass:
		case tokens.TokenFun:
		case tokens.TokenVar:
		case tokens.TokenFor:
		case tokens.TokenIf:
		case tokens.TokenWhile:
		case tokens.TokenPrint:
		case tokens.TokenReturn:
			return
		default: // Do nothing.
		}

		advance()
	}
}

func declaration() {
	switch {
	case match(tokens.TokenFun):
		funDeclaration()
	case match(tokens.TokenVar):
		varDeclaration()
	default:
		statement()
	}

	if gParser.panicMode {
		synchronize()
	}
}

func statement() {
	switch {
	case match(tokens.TokenPrint):
		printStatement()
	case match(tokens.TokenFor):
		forStatement()
	case match(tokens.TokenIf):
		ifStatement()
	case match(tokens.TokenWhile):
		whileStatement()
	case match(tokens.TokenReturn):
		returnStatement()
	case match(tokens.TokenLeftBrace):
		func() {
			beginScope()
			defer endScope()
			block()
		}()
	default:
		expressionStatement()
	}
}

func expressionStatement() {
	expression()
	consume(tokens.TokenSemicolon, "Expect ';' after expression.")
	emitOpcode(bytecode.OpPop)
}

func ifStatement() {
	consume(tokens.TokenLeftParen, "Expect '(' after 'if'.")
	expression()
	consume(tokens.TokenRightParen, "Expect ')' after condition.")

	// start of if execution
	// (1.) eval the condition
	// if condition is false, jump to else (3.)
	// pop condition and continue otherwise
	thenJump := emitJump(bytecode.OpJumpIfFalse)
	emitOpcode(bytecode.OpPop)
	statement()
	// (2.) iftrue statement execution ended
	// jump to the end of else (5.)
	elseJump := emitJump(bytecode.OpJump)

	// (3.) end of iftrue, (1.) will jump here if condition is false
	// pop condition and continue.
	patchJump(thenJump)
	emitOpcode(bytecode.OpPop)

	// (4.) else statement execution
	// if there is no else, jump to the end of if
	// otherwise, continue
	if match(tokens.TokenElse) {
		statement()
	}
	// (5.) end of else (end of if).
	patchJump(elseJump)
}

func whileStatement() {
	loopStart := currentChunk().Count
	consume(tokens.TokenLeftParen, "Expect '(' after 'while'.")
	expression()
	consume(tokens.TokenRightParen, "Expect ')' after condition.")

	exitJump := emitJump(bytecode.OpJumpIfFalse)
	emitOpcode(bytecode.OpPop)
	statement()
	emitLoop(loopStart)

	patchJump(exitJump)
	emitOpcode(bytecode.OpPop)
}

func forStatement() {
	beginScope()
	defer endScope()

	consume(tokens.TokenLeftParen, "Expect '(' after 'for'.")

	if match(tokens.TokenSemicolon) {
		// No initializer.
	} else if match(tokens.TokenVar) {
		varDeclaration()
	} else {
		expressionStatement()
	}

	loopStart := currentChunk().Count
	exitJump := -1

	if !match(tokens.TokenSemicolon) {
		expression()
		consume(tokens.TokenSemicolon, "Expect ';' after loop condition.")

		exitJump = emitJump(bytecode.OpJumpIfFalse)
		emitOpcode(bytecode.OpPop) // Condition.
	}

	if !match(tokens.TokenRightParen) {
		bodyJump := emitJump(bytecode.OpJump)
		incrementStart := currentChunk().Count
		expression()
		emitOpcode(bytecode.OpPop) // discard expression result
		consume(tokens.TokenRightParen, "Expect ')' after for clauses.")

		emitLoop(loopStart)
		loopStart = incrementStart
		patchJump(bodyJump)
	}

	statement()
	emitLoop(loopStart)

	if exitJump != -1 {
		patchJump(exitJump)
		emitOpcode(bytecode.OpPop) // Condition.
	}
}

func number(ParsePrecedence) {
	v, err := strconv.ParseFloat(gParser.previous.LexemeAsString(), 64)
	if err != nil {
		errorAtPrev(err.Error())
	}
	emitConstant(vmvalue.NumberAsValue(v))
}

func string_(ParsePrecedence) {
	t := gParser.previous
	chars := t.Source[t.Start+1 : t.Start+t.Length-1]
	str := hashtable.StringInternCopy(chars)
	emitConstant(vmvalue.ObjAsValue(str))
}

func namedVariable(name scanner.Token, precedence ParsePrecedence) {
	canAssign := precedence.CanAssign()
	var getOp, setOp bytecode.OpCode

	arg, ok := resolveLocal(gCurrent, &name)
	if ok {
		getOp = bytecode.OpGetLocal
		setOp = bytecode.OpSetLocal
	} else if arg, ok = resolveUpvalue(gCurrent, &name); ok {
		getOp = bytecode.OpGetUpvalue
		setOp = bytecode.OpSetUpvalue
	} else {
		arg = identifierConstant(&name)
		getOp = bytecode.OpGetGlobal
		setOp = bytecode.OpSetGlobal
	}

	if canAssign && match(tokens.TokenEqual) {
		expression()
		emitOpByte(setOp, byte(arg))
	} else {
		emitOpByte(getOp, byte(arg))
	}
}

func variable(precedence ParsePrecedence) {
	namedVariable(gParser.previous, precedence)
}

func grouping(ParsePrecedence) {
	expression()
	consume(tokens.TokenRightParen, "Expect ')' after expression.")
}

func literal(ParsePrecedence) {
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

func binary(ParsePrecedence) {
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

func call(ParsePrecedence) {
	argCount := argumentList()
	emitOpByte(bytecode.OpCall, argCount)
}

func argumentList() byte {
	argCount := 0
	if !check(tokens.TokenRightParen) {
		for {
			expression()
			argCount++
			if argCount > MaxArity {
				errorAtPrev("Can't have more than 255 arguments.")
			}
			if !match(tokens.TokenComma) {
				break
			}
		}
	}
	consume(tokens.TokenRightParen, "Expect ')' after arguments.")
	return byte(argCount)
}

func unary(ParsePrecedence) {
	operatorType := gParser.previous.Type
	parsePrecedence(PrecedenceUnary)

	// emit the operator instruction
	switch operatorType {
	case tokens.TokenBang:
		emitOpcode(bytecode.OpNot)
	case tokens.TokenMinus:
		emitOpcode(bytecode.OpNegate)
	default:
		panic("Unreachable unary: " + gParser.previous.LexemeAsString())
	}
}

func mustGetRule(t tokens.TokenType) *ParseRule {
	if r, ok := rules[t]; ok {
		return r
	} else {
		panic(fmt.Sprintf("get rule %s (%d)", t, t))
	}
}

func consume(t tokens.TokenType, message string) {
	if gParser.current.Type == t {
		advance()
		return
	}

	errorAtCurrent(message)
}

func match(t tokens.TokenType) bool {
	if !check(t) {
		return false
	}
	advance()
	return true
}

func check(t tokens.TokenType) bool {
	return gParser.current.Type == t
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
		fmt.Fprintf(os.Stderr, " at '%s'", token.LexemeAsString())
	}

	fmt.Fprintf(os.Stderr, ": %s\n", message)
	gParser.hadError = true
}

var rules map[tokens.TokenType]*ParseRule

func init() {
	rules = map[tokens.TokenType]*ParseRule{
		tokens.TokenLeftParen:    {grouping, call, PrecedenceCall},
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
		tokens.TokenIdentifier:   {variable, nil, PrecedenceNone},
		tokens.TokenString:       {string_, nil, PrecedenceNone},
		tokens.TokenNumber:       {number, nil, PrecedenceNone},
		tokens.TokenAnd:          {nil, and_, PrecedenceAnd},
		tokens.TokenClass:        {nil, nil, PrecedenceNone},
		tokens.TokenElse:         {nil, nil, PrecedenceNone},
		tokens.TokenFalse:        {literal, nil, PrecedenceNone},
		tokens.TokenFor:          {nil, nil, PrecedenceNone},
		tokens.TokenFun:          {nil, nil, PrecedenceNone},
		tokens.TokenIf:           {nil, nil, PrecedenceNone},
		tokens.TokenNil:          {literal, nil, PrecedenceNone},
		tokens.TokenOr:           {nil, or_, PrecedenceOr},
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
