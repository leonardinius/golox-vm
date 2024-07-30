package vmvalue

import (
	"math"
	"unsafe"

	"github.com/leonardinius/goloxvm/internal/vm/vmobject"
)

// Value we use NaN boxing, NaN tagging here.
// See https://craftinginterpreters.com/optimization.html.
type Value uint64

type ValueType byte

const (
	_ ValueType = iota
	ValBool
	ValNil
	ValNumber
	ValObj
)

const (
	QNAN     uint64 = 0x7ffc000000000000
	SignBit  uint64 = 1 << 63
	TagNil   uint64 = 1 // 01
	TagFalse uint64 = 2 // 10
	TagTrue  uint64 = 3 // 11
	//
	NilValue   Value = Value(QNAN | TagNil)
	TrueValue  Value = Value(QNAN | TagTrue)
	FalseValue Value = Value(QNAN | TagFalse)
)

func IsBool(v Value) bool {
	return (v | 1) == TrueValue
}

func IsTrue(v Value) bool {
	return (v) == TrueValue
}

func IsFalse(v Value) bool {
	return (v) == FalseValue
}

func IsNil(v Value) bool {
	return (v) == NilValue
}

func IsNumber(v Value) bool {
	return (uint64(v) & QNAN) != QNAN
}

func IsObj(v Value) bool {
	return (uint64(v) & (QNAN | SignBit)) == (QNAN | SignBit)
}

func IsValuesEqual(v1, v2 Value) bool {
	if IsNumber(v1) && IsNumber(v2) {
		return ValueAsNumber(v1) == ValueAsNumber(v2)
	}
	return v1 == v2
}

func ValueAsBool(v Value) bool {
	return v == TrueValue
}

func ValueAsNumber(v Value) float64 {
	return math.Float64frombits(uint64(v))
}

func NumberAsValue(num float64) Value {
	return Value(math.Float64bits(num))
}

func BoolAsValue(b bool) Value {
	if b {
		return TrueValue
	}
	return FalseValue
}

func ObjAsValue[T vmobject.VMObjectable](v *T) Value {
	ptr := uintptr(unsafe.Pointer(v)) //nolint:gosec // unsafe.Pointer is used here
	return Value(SignBit | QNAN | uint64(ptr))
}

func ValueAsObjType[T vmobject.VMObjectable](v Value) *T {
	addr := (uint64(v) & ^(SignBit | QNAN))
	ptr := uintptr(addr)
	uptr := (**T)(unsafe.Pointer(&ptr)) //nolint:gosec // unsafe.Pointer is used here
	return *uptr
}

func ValueAsObj(v Value) *vmobject.Obj {
	return ValueAsObjType[vmobject.Obj](v)
}

func ObjTypeTag(v Value) vmobject.ObjType {
	return ValueAsObj(v).Type
}

func isObjType(v Value, objType vmobject.ObjType) bool {
	return IsObj(v) && ObjTypeTag(v) == objType
}

func IsString(v Value) bool {
	return isObjType(v, vmobject.ObjTypeString)
}

func IsFunction(v Value) bool {
	return isObjType(v, vmobject.ObjTypeFunction)
}

func ValueAsString(v Value) *vmobject.ObjString {
	return ValueAsObjType[vmobject.ObjString](v)
}

func ValueAsStringChars(v Value) []byte {
	return ValueAsObjType[vmobject.ObjString](v).Chars
}

func ValueAsFunction(v Value) *vmobject.ObjFunction {
	return ValueAsObjType[vmobject.ObjFunction](v)
}
