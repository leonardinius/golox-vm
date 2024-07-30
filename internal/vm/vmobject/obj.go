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
	ObjTypeFunction
	ObjTypeString
)

type VMObjectable interface {
	Obj | ObjString | ObjFunction
}

type Obj struct {
	Type ObjType
	Next *Obj
}

type ObjFunction struct {
	Obj
	Arity       int
	ChunkPtr    uintptr
	FreeChunkFn func()
	Name        *ObjString
}

type ObjString struct {
	Obj
	Chars []byte
	Hash  uint64
}

var (
	GRoots = (*Obj)(nil)
	gSeed  = maphash.MakeSeed()

	GObjSize         = int(unsafe.Sizeof(Obj{}))
	GObjStringSize   = int(unsafe.Sizeof(ObjString{}))
	GObjFunctionSize = int(unsafe.Sizeof(ObjFunction{}))
)

func castObject[T VMObjectable](o *Obj) *T {
	return (*T)(unsafe.Pointer(o)) //nolint:gosec // unsafe.Pointer is used here
}

func allocateObject[T VMObjectable](objType ObjType, sizeBytes int) *T {
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
	case ObjTypeFunction:
		f := castObject[ObjFunction](obj)
		if f.FreeChunkFn != nil {
			f.FreeChunkFn()
		}
		vmmem.FreeUnsafePtr[byte](f, GObjFunctionSize)
	default:
		panic(fmt.Sprintf("unable to free object of type %d", obj.Type))
	}
}

func NewTakeString(chars []byte, hash uint64) *ObjString {
	value := allocateObject[ObjString](ObjTypeString, GObjStringSize)
	value.Chars = chars
	value.Hash = hash
	return value
}

func NewCopyString(chars []byte, hash uint64) *ObjString {
	cloned := vmmem.AllocateSlice[byte](len(chars))
	copy(cloned, chars)
	return NewTakeString(cloned, hash)
}

func NewFunction(chunkPtr uintptr) *ObjFunction {
	value := allocateObject[ObjFunction](ObjTypeFunction, GObjFunctionSize)
	value.Arity = 0
	value.Name = nil
	value.ChunkPtr = chunkPtr
	value.FreeChunkFn = nil
	return value
}

func PrintObject(obj *Obj) {
	switch obj.Type {
	case ObjTypeString:
		value := string(castObject[ObjString](obj).Chars)
		fmt.Print(value)
	case ObjTypeFunction:
		f := castObject[ObjFunction](obj)
		printFunction(f)
	default:
		panic(fmt.Sprintf("unable to print object of type %d", obj.Type))
	}
}

func printFunction(f *ObjFunction) {
	if f.Name == nil {
		fmt.Print("<script>")
		return
	}

	name := string(f.Name.Chars)
	fmt.Print("<fn " + name + ">")
}

func HashString(chars []byte) uint64 {
	return maphash.Bytes(gSeed, chars)
}
