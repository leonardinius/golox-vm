package vmvalue

import (
	"fmt"
	"hash/maphash"
	"slices"
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

var (
	_                = int(unsafe.Sizeof(Obj{}))
	gObjStringSize   = int(unsafe.Sizeof(ObjString{}))
	gObjFunctionSize = int(unsafe.Sizeof(ObjFunction{}))
	gObjNativeSize   = int(unsafe.Sizeof(ObjNative{}))
	gObjClosureSize  = int(unsafe.Sizeof(ObjClosure{}))
	gObjUpvalueSize  = int(unsafe.Sizeof(ObjUpvalue{}))
)

type Obj struct {
	Type   ObjType
	Marked bool
	Next   *Obj
}

type ObjString struct {
	Obj
	Chars []byte
	Hash  uint64
}

func NewTakeString(chars []byte, hash uint64) *ObjString {
	obj := allocateObject[ObjString](ObjTypeString, gObjStringSize)
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

type ObjFunction struct {
	Obj
	Arity                int
	Chunk                any
	ChunkFreeFn          func()
	ChunkMarkConstantsFn func()
	UpvalueCount         int
	Name                 *ObjString
}

func NewFunction(chunk any, chunkFreeFn, chunkMarkFn func()) *ObjFunction {
	obj := allocateObject[ObjFunction](ObjTypeFunction, gObjFunctionSize)
	obj.Chunk = chunk
	obj.ChunkFreeFn = chunkFreeFn
	obj.ChunkMarkConstantsFn = chunkMarkFn
	obj.Arity = 0
	obj.UpvalueCount = 0
	obj.Name = nil
	return obj
}

type NativeFn func(args ...Value) (Value, error)

type ObjNative struct {
	Obj
	Fn    NativeFn
	Arity byte
}

func NewNativeFunction(fn NativeFn, arity byte) *ObjNative {
	obj := allocateObject[ObjNative](ObjTypeNative, gObjNativeSize)
	obj.Fn = fn
	obj.Arity = arity
	return obj
}

type ObjClosure struct {
	Obj
	Fn       *ObjFunction
	Upvalues []*ObjUpvalue
}

func NewClosure(fn *ObjFunction) *ObjClosure {
	obj := allocateObject[ObjClosure](ObjTypeClosure, gObjClosureSize)
	obj.Fn = fn
	obj.Upvalues = vmmem.AllocateSlice[*ObjUpvalue](fn.UpvalueCount)
	return obj
}

type ObjUpvalue struct {
	Obj
	Location *Value
	Closed   Value
	Next     *ObjUpvalue
}

func NewUpvalue(slot *Value) *ObjUpvalue {
	obj := allocateObject[ObjUpvalue](ObjTypeUpvalue, gObjUpvalueSize)
	obj.Location = slot
	obj.Closed = NilValue
	obj.Next = nil
	return obj
}

func InitObjects() {
	GRoots = nil
	gcTrace = gcTraceStack{}
	gSeed = maphash.MakeSeed()
}

func FreeObjects() {
	for GRoots != nil {
		var obj *Obj = GRoots.Next
		FreeObject(GRoots)
		GRoots = obj
	}
	gcTrace.grayStack = vmmem.FreeSlice(gcTrace.grayStack)
}

func FreeObject(obj *Obj) {
	switch obj.Type {
	case ObjTypeString:
		v := castObject[ObjString](obj)
		debugPrintFreeObject(obj, gObjStringSize)
		v.Chars = vmmem.FreeSlice(v.Chars)
		vmmem.TriggerGC(gObjStringSize, 1, 0)
	case ObjTypeFunction:
		v := castObject[ObjFunction](obj)
		debugPrintFreeObject(obj, gObjFunctionSize)
		v.ChunkFreeFn()
		vmmem.TriggerGC(gObjFunctionSize, 1, 0)
	case ObjTypeNative:
		debugPrintFreeObject(obj, gObjNativeSize)
		vmmem.TriggerGC(gObjNativeSize, 1, 0)
	case ObjTypeClosure:
		v := castObject[ObjClosure](obj)
		debugPrintFreeObject(obj, gObjClosureSize)
		v.Upvalues = vmmem.FreeSlice(v.Upvalues)
		vmmem.TriggerGC(gObjClosureSize, 1, 0)
	case ObjTypeUpvalue:
		debugPrintFreeObject(obj, gObjUpvalueSize)
		vmmem.TriggerGC(gObjUpvalueSize, 1, 0)
	default:
		panic(fmt.Sprintf("unable to free object of type %d", obj.Type))
	}
}

func PrintAnyObject[T VMObjectable](o *T) {
	PrintObject(castObjectable(o))
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

//go:nosplit
//go:nocheckptr
func castObject[T VMObjectable](o *Obj) *T {
	return (*T)(unsafe.Pointer(o)) //nolint:gosec
}

//go:nosplit
//go:nocheckptr
func castObjectable[T VMObjectable](o *T) *Obj {
	return (*Obj)(unsafe.Pointer(o)) //nolint:gosec
}

func allocateObject[T VMObjectable](objType ObjType, size int) *T {
	vmmem.TriggerGC(size, 0, 1)
	o := new(T)
	object := castObjectable(o)
	object.Type = objType
	object.Marked = false
	object.Next = GRoots
	GRoots = object
	debugPrintAllocateObject(object, size)
	return o
}

type gcTraceStack struct {
	grayStack []*Obj
}

var gcTrace gcTraceStack = gcTraceStack{}

func MarkObject[T VMObjectable](o *T) {
	if o == nil {
		return
	}
	obj := castObjectable(o)
	if obj.Marked {
		return
	}

	if len(gcTrace.grayStack)+1 > cap(gcTrace.grayStack) {
		newCapacity := vmmem.GrowCapacity(cap(gcTrace.grayStack))
		gcTrace.grayStack = slices.Grow(gcTrace.grayStack, newCapacity)
	}

	debugPrintMarkObject(obj)
	obj.Marked = true
	gcTrace.grayStack = append(gcTrace.grayStack, obj)
}

func GCTraceReferences() {
	//nolint:intrange // grayStack gwors during blacken
	for i := 0; i < len(gcTrace.grayStack); i++ {
		obj := gcTrace.grayStack[i]
		blackenObject(obj)
	}
	gcTrace.grayStack = gcTrace.grayStack[:0]
}

func blackenObject(obj *Obj) {
	debugPrintBlackenObject(obj)

	switch obj.Type {
	case ObjTypeString, ObjTypeNative:
		// do nothing
	case ObjTypeUpvalue:
		v := castObject[ObjUpvalue](obj)
		MarkValue(v.Closed)
	case ObjTypeFunction:
		v := castObject[ObjFunction](obj)
		MarkObject(v.Name)
		v.ChunkMarkConstantsFn()
	case ObjTypeClosure:
		v := castObject[ObjClosure](obj)
		MarkObject(v.Fn)
		for i := range v.Upvalues {
			MarkObject(v.Upvalues[i])
		}
	default:
		panic(fmt.Sprintf("unable to print object of type %d", obj.Type))
	}
}

func GCSweep() {
	var previous *Obj
	obj := GRoots

	for obj != nil {
		if obj.Marked {
			obj.Marked = false
			previous = obj
			obj = obj.Next
		} else {
			unreached := obj
			obj = obj.Next
			if previous != nil {
				previous.Next = obj
			} else {
				GRoots = obj
			}
			FreeObject(unreached)
		}
	}
}
