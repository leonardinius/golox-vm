package hashtable

import (
	"github.com/leonardinius/goloxvm/internal/vm/vmvalue"
)

var gGlobalEnv Table

func InitGlobals() {
	gGlobalEnv = NewHashtable()
}

func FreeGlobals() {
	gGlobalEnv.Free()
}

func SetGlobal(name *vmvalue.ObjString, value vmvalue.Value) bool {
	return gGlobalEnv.Set(name, value)
}

func GetGlobal(name *vmvalue.ObjString) (vmvalue.Value, bool) {
	return gGlobalEnv.Get(name)
}

func DeleteGlobal(name *vmvalue.ObjString) bool {
	return gGlobalEnv.Delete(name)
}
