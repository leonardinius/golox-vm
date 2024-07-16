package vmvalue

import (
	"math"
)

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

func IsEqual(v1, v2 Value) bool {
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

// #define AS_CSTRING(value) (((ObjString *)AS_OBJ(value))->chars)
// #define AS_OBJ(value) ((Obj *)(uintptr_t)((value) & ~(SIGN_BIT | QNAN)))
// #define OBJ_VAL(obj) (Value)(SIGN_BIT | QNAN | (uint64_t)(uintptr_t)(obj))

// static inline bool isObjType(Value value, ObjType type) {
//     return IS_OBJ(value) && AS_OBJ(value)->type == type;
// }
// #define OBJ_TYPE(value) (AS_OBJ(value)->type)
// #define IS_BOUND_METHOD(value) isObjType(value, OBJ_BOUND_METHOD)
// #define IS_CLASS(value) isObjType(value, OBJ_CLASS)
// #define IS_CLOSURE(value) isObjType(value, OBJ_CLOSURE)
// #define IS_FUNCTION(value) isObjType(value, OBJ_FUNCTION)
// #define IS_INSTANCE(value) isObjType(value, OBJ_INSTANCE)
// #define IS_NATIVE(value) isObjType(value, OBJ_NATIVE)
// #define IS_STRING(value) isObjType(value, OBJ_STRING)

// #define AS_BOUND_METHOD(value) ((ObjBoundMethod *)AS_OBJ(value))
// #define AS_CLASS(value) ((ObjClass *)AS_OBJ(value))
// #define AS_CLOSURE(value) ((ObjClosure *)AS_OBJ(value))
// #define AS_FUNCTION(value) ((ObjFunction *)AS_OBJ(value))
// #define AS_INSTANCE(value) ((ObjInstance *)AS_OBJ(value))
// #define AS_NATIVE(value) (((ObjNative *)AS_OBJ(value))->function)
// #define AS_STRING(value) ((ObjString *)AS_OBJ(value))
// #define AS_CSTRING(value) (((ObjString *)AS_OBJ(value))->chars)
