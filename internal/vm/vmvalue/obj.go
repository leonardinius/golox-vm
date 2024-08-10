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

var gObjTypeStrings = map[ObjType]string{
	ObjTypeString:   "OBJ_STRING",
	ObjTypeFunction: "OBJ_FUNCTION",
	ObjTypeNative:   "OBJ_NATIVE",
	ObjTypeClosure:  "OBJ_CLOSURE",
}

// String implements fmt.Stringer.
func (op ObjType) String() string {
	if str, ok := gObjTypeStrings[op]; ok {
		return str
	}

	panic(fmt.Sprintf("unknown object type: %d", op))
}

type VMObjectable interface {
	Obj | ObjString | ObjFunction | ObjNative | ObjClosure
}

type vmGc interface {
	gc()
}

var (
	_ vmGc = (*Obj)(nil)
	_ vmGc = (*ObjString)(nil)
	_ vmGc = (*ObjFunction)(nil)
	_ vmGc = (*ObjNative)(nil)
	_ vmGc = (*ObjClosure)(nil)
)

type Obj struct {
	Type ObjType
	Next *Obj
}

// gc implements vmGc.
func (o *Obj) gc() {
	panic(fmt.Sprintf("unable to free object of type %d", o.Type))
}

type ObjFunction struct {
	Obj
	Arity       int
	Chunk       any
	FreeChunkFn func()
	Name        *ObjString
}

// gc implements vmGc.
// Subtle: this method shadows the method (Obj).gc of ObjFunction.Obj.
func (o *ObjFunction) gc() {
	o.FreeChunkFn()
}

type NativeFn func(args ...Value) Value

type ObjNative struct {
	Obj
	Fn    NativeFn
	Arity byte
}

// gc implements vmGc.
// Subtle: this method shadows the method (Obj).gc of ObjNative.Obj.
func (o *ObjNative) gc() {
}

type ObjClosure struct {
	Obj
	Fn *ObjFunction
}

// gc implements vmGc.
// Subtle: this method shadows the method (Obj).gc of ObjClosure.Obj.
func (o *ObjClosure) gc() {
}

type ObjString struct {
	Obj
	Chars []byte
	Hash  uint64
}

// gc implements vmGc.
// Subtle: this method shadows the method (Obj).gc of ObjString.Obj.
func (o *ObjString) gc() {
	o.Chars = vmmem.FreeSlice(o.Chars)
}

var (
	GRoots = (*Obj)(nil)
	gSeed  = maphash.MakeSeed()
)

func castObject[T VMObjectable](o *Obj) *T {
	return (*T)(unsafe.Pointer(o)) //nolint:gosec
}

func allocateObject[T VMObjectable](objType ObjType) *T {
	o := new(T)
	object := (*Obj)(unsafe.Pointer(o)) //nolint:gosec
	object.Type = objType
	object.Next = GRoots
	GRoots = object
	return o
}

func FreeObjects() {
	for GRoots != nil {
		var obj *Obj = GRoots.Next
		FreeObject(GRoots)
		GRoots = obj
	}
}

func FreeObject(o *Obj) {
	DebugFreeObject(o, "free")
	switch o.Type {
	case ObjTypeString:
		v := castObject[ObjString](o)
		v.gc()
	case ObjTypeFunction:
		v := castObject[ObjFunction](o)
		v.gc()
	case ObjTypeNative:
		v := castObject[ObjNative](o)
		v.gc()
	case ObjTypeClosure:
		v := castObject[ObjClosure](o)
		v.gc()
	default:
		o.gc()
	}
}

func NewTakeString(chars []byte, hash uint64) *ObjString {
	value := allocateObject[ObjString](ObjTypeString)
	value.Chars = chars
	value.Hash = hash
	return value
}

func NewCopyString(chars []byte, hash uint64) *ObjString {
	cloned := vmmem.AllocateSlice[byte](len(chars))
	copy(cloned, chars)
	return NewTakeString(cloned, hash)
}

func NewFunction(chunk any, freeChunkFn func()) *ObjFunction {
	value := allocateObject[ObjFunction](ObjTypeFunction)
	value.Arity = 0
	value.Name = nil
	value.Chunk = chunk
	value.FreeChunkFn = freeChunkFn
	return value
}

func NewNativeFunction(fn NativeFn, arity byte) *ObjNative {
	value := allocateObject[ObjNative](ObjTypeNative)
	value.Fn = fn
	value.Arity = arity
	return value
}

func NewClosure(fn *ObjFunction) *ObjClosure {
	value := allocateObject[ObjClosure](ObjTypeClosure)
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
