package bytecode

import "fmt"

type OpCode byte

const (
	_ OpCode = iota
	OpConstant
	OpNil
	OpTrue
	OpFalse
	OpEqual
	OpGreater
	OpLess
	OpPop
	OpAdd
	OpSubtract
	OpMultiply
	OpDivide
	OpNot
	OpNegate
	OpReturn
)

var gOpCodeStrings = map[OpCode]string{
	OpConstant: "OP_CONSTANT",
	OpNil:      "OP_NIL",
	OpTrue:     "OP_TRUE",
	OpFalse:    "OP_FALSE",
	OpEqual:    "OP_EQUAL",
	OpGreater:  "OP_GREATER",
	OpLess:     "OP_LESS",
	OpPop:      "OP_POP",
	OpAdd:      "OP_ADD",
	OpSubtract: "OP_SUBTRACT",
	OpMultiply: "OP_MULTIPLY",
	OpDivide:   "OP_DIVIDE",
	OpNot:      "OP_NOT",
	OpNegate:   "OP_NEGATE",
	OpReturn:   "OP_RETURN",
}

func (op OpCode) String() string {
	if str, ok := gOpCodeStrings[op]; ok {
		return str
	}

	panic(fmt.Sprintf("unknown opcode: %d", op))
}
