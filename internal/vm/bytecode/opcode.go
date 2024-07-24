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
	OpGetGlobal
	OpDefineGlobal
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
	OpGetGlobal:    "OP_GET_GLOBAL",
	OpDefineGlobal: "OP_DEFINE_GLOBAL",
	OpSetGlobal:    "OP_SET_GLOBAL",
	OpAdd:          "OP_ADD",
	OpSubtract:     "OP_SUBTRACT",
	OpMultiply:     "OP_MULTIPLY",
	OpDivide:       "OP_DIVIDE",
	OpNot:          "OP_NOT",
	OpNegate:       "OP_NEGATE",
	OpPrint:        "OP_PRINT",
	OpReturn:       "OP_RETURN",
}

func (op OpCode) String() string {
	if str, ok := gOpCodeStrings[op]; ok {
		return str
	}

	panic(fmt.Sprintf("unknown opcode: %d", op))
}
