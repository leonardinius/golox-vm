package vm

import (
	"fmt"
	"math"
	"os"
	"runtime"

	"github.com/leonardinius/goloxvm/internal/vm/bytecode"
	"github.com/leonardinius/goloxvm/internal/vm/hashtable"
	"github.com/leonardinius/goloxvm/internal/vm/vmchunk"
	"github.com/leonardinius/goloxvm/internal/vm/vmdebug"
	"github.com/leonardinius/goloxvm/internal/vm/vmmem"
	"github.com/leonardinius/goloxvm/internal/vm/vmstd"
	"github.com/leonardinius/goloxvm/internal/vm/vmvalue"
	"github.com/leonardinius/goloxvm/internal/vmcompiler"
)

const (
	MaxCallFrames = 64
	MaxStackCount = MaxCallFrames * (math.MaxUint8 + 1)
)

type CallFrame struct {
	Closure  *vmvalue.ObjClosure
	IP       int
	SlotsTop int
}

// VM is the virtual machine.
type VM struct {
	Frames       [MaxCallFrames]CallFrame
	FrameCount   int
	Stack        [MaxStackCount]vmvalue.Value
	StackTop     int
	OpenUpvalues *vmvalue.ObjUpvalue
}

var GlobalVM VM

type InterpretError int

const (
	_ InterpretError = iota
	InterpretCompileError
	InterpretRuntimeError
)

func (i InterpretError) Error() string {
	var err string
	switch i {
	case InterpretCompileError:
		err = "compile error"
	case InterpretRuntimeError:
		err = "runtime error"
	default:
		err = fmt.Sprintf("unknown error %d", i)
	}

	return err
}

func InitVM() {
	vmmem.SetGarbageCollector(GC)
	vmmem.SetGarbageCollectorRetain(func(v uint64) { Push(vmvalue.NanBoxedAsValue(v)) })
	vmmem.SetGarbageCollectorRelease(func() { _ = Pop() })
	hashtable.InitInternStrings()
	hashtable.InitGlobals()
	vmvalue.InitObjects()
	defineNative("clock", vmstd.StdClockNative, 0)
	resetStack()
}

func FreeVM() {
	hashtable.FreeGlobals()
	hashtable.FreeInternStrings()
	vmvalue.FreeObjects()
	resetStack()
}

func resetStack() {
	GlobalVM.StackTop = 0
	GlobalVM.FrameCount = 0
	GlobalVM.OpenUpvalues = nil
}

func Interpret(code []byte) (vmvalue.Value, error) {
	var fn *vmvalue.ObjFunction
	var ok bool

	if fn, ok = vmcompiler.Compile(code); !ok {
		return vmvalue.NilValue, InterpretCompileError
	}

	Push(vmvalue.ObjAsValue(fn))
	closure := vmvalue.NewClosure(fn)
	Pop()
	Push(vmvalue.ObjAsValue(closure))
	Call(closure, 0)

	return Run()
}

func traceInstruction(frame *CallFrame, chunk *vmchunk.Chunk) {
	if GlobalVM.StackTop > 0 {
		fmt.Print("        ")
		for i := range GlobalVM.StackTop {
			fmt.Print("[ ")
			vmdebug.PrintValue(StackAt(i))
			fmt.Print(" ]")
		}
		fmt.Println()
	}
	vmdebug.DisassembleInstruction(chunk, frame.IP)
}

func Push(value vmvalue.Value) {
	GlobalVM.Stack[GlobalVM.StackTop] = value
	GlobalVM.StackTop++
}

func Pop() vmvalue.Value {
	GlobalVM.StackTop--
	return GlobalVM.Stack[GlobalVM.StackTop]
}

func Peek(distance byte) vmvalue.Value {
	return GlobalVM.Stack[GlobalVM.StackTop-1-int(distance)]
}

func StackAt(at int) vmvalue.Value {
	return GlobalVM.Stack[at]
}

func SetStackAt(at int, v vmvalue.Value) {
	GlobalVM.Stack[at] = v
}

func CallValue(callee vmvalue.Value, argCount byte) (ok bool) {
	if vmvalue.IsObj(callee) {
		switch vmvalue.ObjTypeTag(callee) {
		case vmvalue.ObjTypeClosure:
			return Call(vmvalue.ValueAsClosure(callee), argCount)
		case vmvalue.ObjTypeNative:
			native := vmvalue.ValueAsNativeFn(callee)
			if argCount != native.Arity {
				return runtimeError("Expected %d arguments but got %d.",
					native.Arity, argCount)
			}

			iArgs := int(argCount)
			args := GlobalVM.Stack[GlobalVM.StackTop-iArgs : GlobalVM.StackTop]
			value := native.Fn(args...)
			GlobalVM.StackTop -= (iArgs + 1)
			Push(value)
			return true
		}
	}

	return runtimeError("Can only call functions and classes.")
}

func CaptureUpvalue(at int) *vmvalue.ObjUpvalue {
	value := &GlobalVM.Stack[at]
	valuePtr := vmvalue.UPtrFromValue(value)

	var prevUpvalue *vmvalue.ObjUpvalue
	upvalue := GlobalVM.OpenUpvalues
	for upvalue != nil && vmvalue.UPtrFromValue(upvalue.Location) > valuePtr {
		prevUpvalue = upvalue
		upvalue = upvalue.Next
	}
	if upvalue != nil && upvalue.Location == value {
		return upvalue
	}

	createdUpvalue := vmvalue.NewUpvalue(value)
	createdUpvalue.Next = upvalue
	if prevUpvalue == nil {
		GlobalVM.OpenUpvalues = createdUpvalue
	} else {
		prevUpvalue.Next = createdUpvalue
	}

	return createdUpvalue
}

func CloseUpvalues(at int) {
	last := &GlobalVM.Stack[at]
	lastPtr := vmvalue.UPtrFromValue(last)
	for GlobalVM.OpenUpvalues != nil &&
		vmvalue.UPtrFromValue(GlobalVM.OpenUpvalues.Location) >= lastPtr {
		upvalue := GlobalVM.OpenUpvalues
		upvalue.Closed = *upvalue.Location
		upvalue.Location = &upvalue.Closed
		GlobalVM.OpenUpvalues = upvalue.Next
	}
}

func Call(closure *vmvalue.ObjClosure, argCount byte) (ok bool) {
	iArgs := int(argCount)
	if iArgs != closure.Fn.Arity {
		return runtimeError("Expected %d arguments but got %d.",
			closure.Fn.Arity, argCount)
	}

	if GlobalVM.FrameCount == MaxCallFrames {
		return runtimeError("Stack overflow.")
	}

	frame := &GlobalVM.Frames[GlobalVM.FrameCount]
	GlobalVM.FrameCount++
	frame.Closure = closure
	frame.IP = 0
	frame.SlotsTop = GlobalVM.StackTop - iArgs - 1
	return true
}

func SetGlobal(name *vmvalue.ObjString, value vmvalue.Value) bool {
	return hashtable.SetGlobal(name, value)
}

func GetGlobal(name *vmvalue.ObjString) (vmvalue.Value, bool) {
	return hashtable.GetGlobal(name)
}

func DeleteGlobal(name *vmvalue.ObjString) bool {
	return hashtable.DeleteGlobal(name)
}

func Run() (vmvalue.Value, error) { //nolint:gocyclo,gocognit
	if vmdebug.DebugDisassembler {
		fmt.Println("== trace execution ==")
		defer fmt.Println()
	}

	ok := true
	frame, chunk := frameChunk()
	for {
		if !ok {
			return vmvalue.NilValue, InterpretRuntimeError
		}

		// Debug tracing.
		if vmdebug.DebugDisassembler {
			// Debug GC issues
			runtime.GC() // TODO: check for failures
			traceInstruction(frame, chunk)
		}

		instruction := bytecode.OpCode(readByte(frame, chunk))
		switch instruction {
		case bytecode.OpConstant:
			constant := readConstant(frame, chunk)
			Push(constant)
		case bytecode.OpNil:
			Push(vmvalue.NilValue)
		case bytecode.OpTrue:
			Push(vmvalue.TrueValue)
		case bytecode.OpFalse:
			Push(vmvalue.FalseValue)
		case bytecode.OpEqual:
			Push(vmvalue.BoolAsValue(vmvalue.IsValuesEqual(Pop(), Pop())))
		case bytecode.OpGreater:
			ok = binaryNumCompareOp(binOpGreater)
		case bytecode.OpLess:
			ok = binaryNumCompareOp(binOpLess)
		case bytecode.OpAdd:
			if vmvalue.IsString(Peek(0)) && vmvalue.IsString(Peek(1)) {
				ok = stringConcat()
			} else if vmvalue.IsNumber(Peek(0)) && vmvalue.IsNumber(Peek(1)) {
				ok = binaryNumMathOp(binOpAdd)
			} else {
				ok = runtimeError("Operands must be two numbers or two strings.")
			}
		case bytecode.OpSubtract:
			ok = binaryNumMathOp(binOpSubtract)
		case bytecode.OpMultiply:
			ok = binaryNumMathOp(binOpMultiply)
		case bytecode.OpDivide:
			ok = binaryNumMathOp(binOpDivide)
		case bytecode.OpNegate:
			ok = opNegate()
		case bytecode.OpNot:
			Push(vmvalue.BoolAsValue(!isTruey(Pop())))
		case bytecode.OpPop:
			Pop()
		case bytecode.OpPrint:
			PrintlnValue(Pop())
		case bytecode.OpGetLocal:
			slot := readByte(frame, chunk)
			local := StackAt(frame.SlotsTop + int(slot))
			Push(local)
		case bytecode.OpSetLocal:
			slot := readByte(frame, chunk)
			SetStackAt(frame.SlotsTop+int(slot), Peek(0))
		case bytecode.OpGetGlobal:
			name := readString(frame, chunk)
			if value, gok := GetGlobal(name); !gok {
				ok = runtimeError("Undefined variable '%s'.", string(name.Chars))
			} else {
				Push(value)
			}
		case bytecode.OpSetGlobal:
			name := readString(frame, chunk)
			if isNewKey := SetGlobal(name, Peek(0)); isNewKey {
				DeleteGlobal(name)
				ok = runtimeError("Undefined variable '%s'.", string(name.Chars))
			}
		case bytecode.OpDefineGlobal:
			name := readString(frame, chunk)
			SetGlobal(name, Peek(0))
			Pop()
		case bytecode.OpJump:
			offset := readShort(frame, chunk)
			frame.IP += int(offset)
		case bytecode.OpJumpIfFalse:
			offset := readShort(frame, chunk)
			if isFalsey(Peek(0)) {
				frame.IP += int(offset)
			}
		case bytecode.OpLoop:
			offset := readShort(frame, chunk)
			frame.IP -= int(offset)
		case bytecode.OpCall:
			argCount := readByte(frame, chunk)
			if ok = CallValue(Peek(argCount), argCount); ok {
				frame, chunk = frameChunk()
			}
		case bytecode.OpClosure:
			fn := vmvalue.ValueAsFunction(readConstant(frame, chunk))
			closure := vmvalue.NewClosure(fn)
			Push(vmvalue.ObjAsValue(closure))

			for i := range closure.Upvalues {
				islocal := readByte(frame, chunk)
				index := readByte(frame, chunk)
				if islocal != 0 {
					closure.Upvalues[i] = CaptureUpvalue(frame.SlotsTop + int(index))
				} else {
					closure.Upvalues[i] = frame.Closure.Upvalues[index]
				}
			}
		case bytecode.OpGetUpvalue:
			slot := readByte(frame, chunk)
			Push(*frame.Closure.Upvalues[slot].Location)
		case bytecode.OpSetUpvalue:
			slot := readByte(frame, chunk)
			*frame.Closure.Upvalues[slot].Location = Peek(0)
		case bytecode.OpCloseUpvalue:
			CloseUpvalues(GlobalVM.StackTop - 1)
			Pop()
		case bytecode.OpReturn:
			callReturnValue := Pop()
			CloseUpvalues(frame.SlotsTop)
			GlobalVM.FrameCount--
			if GlobalVM.FrameCount == 0 {
				Pop()
				return callReturnValue, nil
			}
			GlobalVM.StackTop = frame.SlotsTop
			Push(callReturnValue)
			frame, chunk = frameChunk()
		default:
			ok = runtimeError("Unexpected instruction")
		}
	}
}

func isTruey(value vmvalue.Value) bool {
	if vmvalue.IsBool(value) {
		return vmvalue.ValueAsBool(value)
	}
	return !vmvalue.IsNil(value)
}

func isFalsey(value vmvalue.Value) bool {
	return !isTruey(value)
}

func binaryNumOp(op func(vmvalue.Value, vmvalue.Value) vmvalue.Value) (ok bool) {
	if ok = (vmvalue.IsNumber(Peek(0)) && vmvalue.IsNumber(Peek(1))); !ok {
		runtimeError("Operands must be numbers.")
		return ok
	}

	b := Pop()
	a := Pop()
	Push(op(a, b))
	return ok
}

func binaryNumMathOp(op func(float64, float64) float64) (ok bool) {
	return binaryNumOp(func(a vmvalue.Value, b vmvalue.Value) vmvalue.Value {
		av := vmvalue.ValueAsNumber(a)
		bv := vmvalue.ValueAsNumber(b)
		return vmvalue.NumberAsValue(op(av, bv))
	})
}

func binaryNumCompareOp(op func(float64, float64) bool) (ok bool) {
	return binaryNumOp(func(a vmvalue.Value, b vmvalue.Value) vmvalue.Value {
		av := vmvalue.ValueAsNumber(a)
		bv := vmvalue.ValueAsNumber(b)
		return vmvalue.BoolAsValue(op(av, bv))
	})
}

func opNegate() (ok bool) {
	if ok = vmvalue.IsNumber(Peek(0)); !ok {
		runtimeError("Operand must be a number.")
		return ok
	}
	Push(vmvalue.NumberAsValue(-vmvalue.ValueAsNumber(Pop())))
	return ok
}

func stringConcat() (ok bool) {
	b := vmvalue.ValueAsStringChars(Peek(0))
	a := vmvalue.ValueAsStringChars(Peek(1))
	length := len(a) + len(b)
	chars := vmmem.AllocateSlice[byte](length)
	copy(chars, a)
	copy(chars[len(a):], b)
	str := hashtable.StringInternTake(chars)
	Pop()
	Pop()
	Push(vmvalue.ObjAsValue(str))
	return true
}

func binOpAdd(a, b float64) float64 {
	return a + b
}

func binOpSubtract(a, b float64) float64 {
	return a - b
}

func binOpMultiply(a, b float64) float64 {
	return a * b
}

func binOpDivide(a, b float64) float64 {
	return a / b
}

func binOpGreater(a, b float64) bool {
	return a > b
}

func binOpLess(a, b float64) bool {
	return a < b
}

func frameChunk() (*CallFrame, *vmchunk.Chunk) {
	frame := &GlobalVM.Frames[GlobalVM.FrameCount-1]
	ch := vmchunk.FromPtr(frame.Closure.Fn.Chunk)
	return frame, ch
}

func readByte(frame *CallFrame, chunk *vmchunk.Chunk) byte {
	frame.IP++
	return chunk.Code[frame.IP-1]
}

func readShort(frame *CallFrame, chunk *vmchunk.Chunk) uint16 {
	frame.IP += 2
	return (uint16(chunk.Code[frame.IP-2]) << 8) | uint16(chunk.Code[frame.IP-1])
}

func readConstant(frame *CallFrame, chunk *vmchunk.Chunk) vmvalue.Value {
	frame.IP++
	at := chunk.Code[frame.IP-1]
	return chunk.ConstantAt(int(at))
}

func readString(frame *CallFrame, chunk *vmchunk.Chunk) *vmvalue.ObjString {
	return vmvalue.ValueAsString(readConstant(frame, chunk))
}

func runtimeError(format string, messageAndArgs ...any) (ok bool) {
	fmt.Fprintf(os.Stderr, format, messageAndArgs...)
	fmt.Fprintln(os.Stderr)

	for i := range GlobalVM.FrameCount {
		frame := &GlobalVM.Frames[GlobalVM.FrameCount-1-i]
		fn := frame.Closure.Fn
		chunk := vmchunk.FromPtr(fn.Chunk)
		offset := frame.IP - 1
		line := chunk.Lines.GetLineByOffset(offset)
		fmt.Fprintf(os.Stderr, "[line %d] in ", line)
		if fn.Name == nil {
			fmt.Fprintln(os.Stderr, "script")
		} else {
			fmt.Fprintf(os.Stderr, "%s()\n", string(fn.Name.Chars))
		}
	}

	resetStack()
	return false
}

func PrintlnValue(v vmvalue.Value) {
	vmvalue.PrintlnValue(v)
}

func defineNative(name string, fn vmvalue.NativeFn, arity byte) {
	nameObj := hashtable.StringInternCopy([]byte(name))
	nameValue := vmvalue.ObjAsValue(nameObj)
	Push(nameValue)
	fnObj := vmvalue.NewNativeFunction(fn, arity)
	fnValue := vmvalue.ObjAsValue(fnObj)
	Push(fnValue)
	SetGlobal(nameObj, vmvalue.ObjAsValue(fnObj))
	Pop()
	Pop()
}
