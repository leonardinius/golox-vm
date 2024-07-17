package vmobject

import (
	"fmt"
	"slices"
	"unsafe"

	"github.com/leonardinius/goloxvm/internal/vmmem"
)

type ObjType byte

const (
	_ ObjType = iota
	// ObjTypeBoundMethod
	// ObjTypeClass
	// ObjTypeClosure
	// ObjTypeFunction
	// ObjTypeInstance
	// ObjTypeNative.
	ObjTypeString
	// ObjTypeUpvalue.
)

type VMObjectable interface {
	Obj | ObjString
}

type Obj struct {
	Type ObjType
}

type ObjString struct {
	Obj
	Chars []byte
}

var (
	GObjSize       = int(unsafe.Sizeof(Obj{}))
	GObjStringSize = int(unsafe.Sizeof(ObjString{}))
)

func castObject[T VMObjectable](o *Obj) *T {
	return (*T)(unsafe.Pointer(o)) //nolint:gosec // unsafe.Pointer is used here
}

func AllocateObject[T VMObjectable](objType ObjType, sizeBytes int) *T {
	ptr := vmmem.AllocateUnsafePtr[byte](sizeBytes)
	((*Obj)(ptr)).Type = objType
	return (*T)(ptr)
}

func NewObjString(chars []byte) *ObjString {
	value := AllocateObject[ObjString](ObjTypeString, GObjStringSize)
	value.Chars = chars
	return value
}

func CopyString(chars []byte) *ObjString {
	cloned := vmmem.AllocateSlice[byte](len(chars))
	copy(cloned, chars)
	return NewObjString(cloned)
}

func IsObjectsEqual(a, b *Obj) bool {
	if a.Type != b.Type {
		return false
	}
	switch a.Type {
	case ObjTypeString:
		aval := castObject[ObjString](a)
		bval := castObject[ObjString](b)
		return slices.Equal(aval.Chars, bval.Chars)
	default:
		panic(fmt.Sprintf("unable to compare object of type %d", a.Type))
	}
}

func PrintObject(obj *Obj) {
	switch obj.Type {
	case ObjTypeString:
		svalue := string(castObject[ObjString](obj).Chars)
		fmt.Print("\"" + svalue + "\"")
	default:
		panic(fmt.Sprintf("unable to print object of type %d", obj.Type))
	}
}
