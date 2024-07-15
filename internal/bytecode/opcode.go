package bytecode

type OpCode byte

const (
	_ OpCode = iota
	OpConstant
	OpPop
	OpAdd
	OpSubtract
	OpMultiply
	OpDivide
	OpNegate
	OpReturn
)

var gOpCodeStrings = map[OpCode]string{
	OpConstant: "OP_CONSTANT",
	OpPop:      "OP_POP",
	OpAdd:      "OP_ADD",
	OpSubtract: "OP_SUBTRACT",
	OpMultiply: "OP_MULTIPLY",
	OpDivide:   "OP_DIVIDE",
	OpNegate:   "OP_NEGATE",
	OpReturn:   "OP_RETURN",
}

func (op OpCode) String() string {
	if str, ok := gOpCodeStrings[op]; ok {
		return str
	}

	return "Unknown"
}
