package vmvalue

type ObjType byte

const (
	_ ObjType = iota
	ObjBoundMethod
	ObjClass
	ObjClosure
	ObjFunction
	ObjInstance
	ObjNative
	ObjString
	ObjUpvalue
)

type Obj struct {
	Type     ObjType
	IsMarked bool
	Next     *Obj
}
