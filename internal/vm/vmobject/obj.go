package vmobject

import (
	"fmt"
	"hash/maphash"
	"unsafe"

	"github.com/leonardinius/goloxvm/internal/vm/vmmem"
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
	Next *Obj
}

type ObjString struct {
	Obj
	Chars []byte
	Hash  uint64
}

var (
	GObjSize       = int(unsafe.Sizeof(Obj{}))
	GObjStringSize = int(unsafe.Sizeof(ObjString{}))
	GRoots         = (*Obj)(nil)
	gSeed          = maphash.MakeSeed()
)

func castObject[T VMObjectable](o *Obj) *T {
	return (*T)(unsafe.Pointer(o)) //nolint:gosec // unsafe.Pointer is used here
}

func AllocateObject[T VMObjectable](objType ObjType, sizeBytes int) *T {
	ptr := vmmem.AllocateUnsafePtr[byte](sizeBytes)
	object := ((*Obj)(ptr))
	object.Type = objType
	object.Next = GRoots
	GRoots = object
	return (*T)(ptr)
}

func FreeObjects() {
	obj := GRoots
	for obj != nil {
		next := obj.Next
		FreeObject(obj)
		obj = next
	}
}

func FreeObject(obj *Obj) {
	switch obj.Type {
	case ObjTypeString:
		s := castObject[ObjString](obj)
		vmmem.FreeArray(s.Chars)
		vmmem.FreeUnsafePtr[byte](s, GObjStringSize)
	default:
		panic(fmt.Sprintf("unable to free object of type %d", obj.Type))
	}
}

func NewTakeString(chars []byte, hash uint64) *ObjString {
	value := AllocateObject[ObjString](ObjTypeString, GObjStringSize)
	value.Chars = chars
	value.Hash = hash
	return value
}

func NewCopyString(chars []byte, hash uint64) *ObjString {
	cloned := vmmem.AllocateSlice[byte](len(chars))
	copy(cloned, chars)
	return NewTakeString(cloned, hash)
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

func HashString(chars []byte) uint64 {
	return maphash.Bytes(gSeed, chars)
}
