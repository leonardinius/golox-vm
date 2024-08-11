package vmvalue

import (
	"fmt"
	"hash/maphash"
	"unsafe"

	"github.com/leonardinius/goloxvm/internal/vm/vmmem"
)

var (
	GRoots *Obj         = nil
	gSeed  maphash.Seed = maphash.MakeSeed()
)

type ObjType byte

const (
	_ ObjType = iota
	ObjTypeString
	ObjTypeFunction
	ObjTypeNative
	ObjTypeClosure
	ObjTypeUpvalue
)

var gObjTypeStrings = map[ObjType]string{
	ObjTypeString:   "OBJ_STRING",
	ObjTypeFunction: "OBJ_FUNCTION",
	ObjTypeNative:   "OBJ_NATIVE",
	ObjTypeClosure:  "OBJ_CLOSURE",
	ObjTypeUpvalue:  "OBJ_UPVALUE",
}

// String implements fmt.Stringer.
func (op ObjType) String() string {
	if str, ok := gObjTypeStrings[op]; ok {
		return str
	}

	panic(fmt.Sprintf("unknown object type: %d", op))
}

type VMObjectable interface {
	Obj | ObjString | ObjFunction | ObjNative | ObjClosure | ObjUpvalue
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
	_ vmGc = (*ObjUpvalue)(nil)
)

type Obj struct {
	Type ObjType
	Next *Obj
}

// gc implements vmGc.
func (o *Obj) gc() {
	// TODO: vmmeory GC tracking
}

type ObjString struct {
	Obj
	Chars []byte
	Hash  uint64
}

func NewTakeString(chars []byte, hash uint64) *ObjString {
	obj := allocateObject[ObjString](ObjTypeString)
	obj.Chars = chars
	obj.Hash = hash
	return obj
}

func NewCopyString(chars []byte, hash uint64) *ObjString {
	cloned := vmmem.AllocateSlice[byte](len(chars))
	copy(cloned, chars)
	return NewTakeString(cloned, hash)
}

func HashString(chars []byte) uint64 {
	return maphash.Bytes(gSeed, chars)
}

// gc implements vmGc.
func (o *ObjString) gc() {
	o.Chars = vmmem.FreeSlice(o.Chars)
}

type ObjFunction struct {
	Obj
	Arity        int
	Chunk        any
	FreeChunkFn  func()
	UpvalueCount int
	Name         *ObjString
}

func NewFunction(chunk any, freeChunkFn func()) *ObjFunction {
	obj := allocateObject[ObjFunction](ObjTypeFunction)
	obj.Chunk = chunk
	obj.FreeChunkFn = freeChunkFn
	obj.Arity = 0
	obj.UpvalueCount = 0
	obj.Name = nil
	return obj
}

// gc implements vmGc.
func (o *ObjFunction) gc() {
	o.FreeChunkFn()
}

type NativeFn func(args ...Value) Value

type ObjNative struct {
	Obj
	Fn    NativeFn
	Arity byte
}

func NewNativeFunction(fn NativeFn, arity byte) *ObjNative {
	obj := allocateObject[ObjNative](ObjTypeNative)
	obj.Fn = fn
	obj.Arity = arity
	return obj
}

// gc implements vmGc.
func (o *ObjNative) gc() {
}

type ObjClosure struct {
	Obj
	Fn       *ObjFunction
	Upvalues []*ObjUpvalue
}

func NewClosure(fn *ObjFunction) *ObjClosure {
	obj := allocateObject[ObjClosure](ObjTypeClosure)
	obj.Fn = fn
	obj.Upvalues = vmmem.AllocateSlice[*ObjUpvalue](fn.UpvalueCount)
	return obj
}

// gc implements vmGc.
func (o *ObjClosure) gc() {
	o.Upvalues = vmmem.FreeSlice(o.Upvalues)
}

type ObjUpvalue struct {
	Obj
	Location *Value
	Closed   Value
	Next     *ObjUpvalue
}

func NewUpvalue(slot *Value) *ObjUpvalue {
	obj := allocateObject[ObjUpvalue](ObjTypeUpvalue)
	obj.Location = slot
	obj.Closed = NilValue
	obj.Next = nil
	return obj
}

// gc implements vmGc.
func (o *ObjUpvalue) gc() {
}

func InitObjects() {
	GRoots = nil
	gSeed = maphash.MakeSeed()
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
	case ObjTypeUpvalue:
		v := castObject[ObjUpvalue](o)
		v.gc()
	}

	// call shared GC part
	// TODO: vmmeory GC tracking
	o.gc()
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
	case ObjTypeUpvalue:
		fmt.Print("upvalue")
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

func castObject[T VMObjectable](o *Obj) *T {
	return (*T)(unsafe.Pointer(o)) //nolint:gosec
}

func allocateObject[T VMObjectable](objType ObjType) *T {
	o := new(T)                         // TODO: vmmeory GC tracking
	object := (*Obj)(unsafe.Pointer(o)) //nolint:gosec
	object.Type = objType
	object.Next = GRoots
	GRoots = object
	return o
}
