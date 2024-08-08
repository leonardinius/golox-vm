package vmvalue

import (
	"fmt"
	"hash/maphash"
	"unsafe"

	"github.com/leonardinius/goloxvm/internal/vm/vmmem"
)

type ObjType byte

const (
	_ ObjType = iota
	ObjTypeString
	ObjTypeFunction
	ObjTypeNative
	ObjTypeClosure
)

type VMObjectable interface {
	Obj | ObjString | ObjFunction | ObjNative | ObjClosure
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

type NativeFn func(args ...Value) Value

type ObjNative struct {
	Obj
	Fn    NativeFn
	Arity byte
}

type ObjClosure struct {
	Obj
	Fn *ObjFunction
}

type ObjString struct {
	Obj
	Chars []byte
	Hash  uint64
}

var (
	GRoots = (*Obj)(nil)
	gSeed  = maphash.MakeSeed()

	gObjStringSize   = int(unsafe.Sizeof(ObjString{}))
	gObjFunctionSize = int(unsafe.Sizeof(ObjFunction{}))
	gObjNativeFnSize = int(unsafe.Sizeof(ObjNative{}))
	gObjClosureSize  = int(unsafe.Sizeof(ObjClosure{}))
)

func castObject[T VMObjectable](o *Obj) *T {
	return (*T)(unsafe.Pointer(o)) //nolint:gosec
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
		v := castObject[ObjString](obj)
		vmmem.FreeArray(v.Chars)
		vmmem.FreeUnsafePtr[byte](v, gObjStringSize)
	case ObjTypeFunction:
		v := castObject[ObjFunction](obj)
		v.FreeChunkFn()
		vmmem.FreeUnsafePtr[byte](v, gObjFunctionSize)
	case ObjTypeNative:
		v := castObject[ObjNative](obj)
		vmmem.FreeUnsafePtr[byte](v, gObjNativeFnSize)
	case ObjTypeClosure:
		v := castObject[ObjClosure](obj)
		vmmem.FreeUnsafePtr[byte](v, gObjClosureSize)
	default:
		panic(fmt.Sprintf("unable to free object of type %d", obj.Type))
	}
}

func NewTakeString(chars []byte, hash uint64) *ObjString {
	value := allocateObject[ObjString](ObjTypeString, gObjStringSize)
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
	value := allocateObject[ObjFunction](ObjTypeFunction, gObjFunctionSize)
	value.Arity = 0
	value.Name = nil
	value.ChunkPtr = chunkPtr
	value.FreeChunkFn = nil
	return value
}

func NewNativeFunction(fn NativeFn, arity byte) *ObjNative {
	value := allocateObject[ObjNative](ObjTypeNative, gObjNativeFnSize)
	value.Fn = fn
	value.Arity = arity
	return value
}

func NewClosure(fn *ObjFunction) *ObjClosure {
	value := allocateObject[ObjClosure](ObjTypeClosure, gObjClosureSize)
	value.Fn = fn
	return value
}

func PrintObject(obj *Obj) {
	switch obj.Type {
	case ObjTypeString:
		v := string(castObject[ObjString](obj).Chars)
		fmt.Print(v)
	case ObjTypeFunction:
		v := castObject[ObjFunction](obj)
		printFunction(v)
	case ObjTypeNative:
		fmt.Print("<native fn>")
	case ObjTypeClosure:
		v := castObject[ObjClosure](obj)
		printFunction(v.Fn)
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
