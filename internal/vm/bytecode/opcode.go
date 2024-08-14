package bytecode

import "fmt"

type OpCode byte

const (
	_ OpCode = iota
	OpConstant
	OpNil
	OpTrue
	OpFalse
	OpPop
	OpGetLocal
	OpSetLocal
	OpGetUpvalue
	OpSetUpvalue
	OpDefineGlobal
	OpGetGlobal
	OpSetGlobal
	OpEqual
	OpGreater
	OpLess
	OpAdd
	OpSubtract
	OpMultiply
	OpDivide
	OpNot
	OpNegate
	OpPrint
	OpJump
	OpJumpIfFalse
	OpLoop
	OpCall
	OpClosure
	OpCloseUpvalue
	OpClass
	OpReturn
)

var gOpCodeStrings = map[OpCode]string{
	OpConstant:     "OP_CONSTANT",
	OpNil:          "OP_NIL",
	OpTrue:         "OP_TRUE",
	OpFalse:        "OP_FALSE",
	OpEqual:        "OP_EQUAL",
	OpGreater:      "OP_GREATER",
	OpLess:         "OP_LESS",
	OpPop:          "OP_POP",
	OpGetLocal:     "OP_GET_LOCAL",
	OpSetLocal:     "OP_SET_LOCAL",
	OpGetUpvalue:   "OP_GET_UPVALUE",
	OpSetUpvalue:   "OP_SET_UPVALUE",
	OpGetGlobal:    "OP_GET_GLOBAL",
	OpSetGlobal:    "OP_SET_GLOBAL",
	OpDefineGlobal: "OP_DEFINE_GLOBAL",
	OpAdd:          "OP_ADD",
	OpSubtract:     "OP_SUBTRACT",
	OpMultiply:     "OP_MULTIPLY",
	OpDivide:       "OP_DIVIDE",
	OpNot:          "OP_NOT",
	OpNegate:       "OP_NEGATE",
	OpPrint:        "OP_PRINT",
	OpJump:         "OP_JUMP",
	OpJumpIfFalse:  "OP_JUMP_IF_FALSE",
	OpLoop:         "OP_LOOP",
	OpCall:         "OP_CALL",
	OpClosure:      "OP_CLOSURE",
	OpCloseUpvalue: "OP_CLOSE_UPVALUE",
	OpClass:        "OP_CLASS",
	OpReturn:       "OP_RETURN",
}

func (op OpCode) String() string {
	if str, ok := gOpCodeStrings[op]; ok {
		return str
	}

	panic(fmt.Sprintf("unknown opcode: %d", op))
}
