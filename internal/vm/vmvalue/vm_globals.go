package vmvalue

var gGlobalEnv Table

func InitGlobals() {
	gGlobalEnv = NewHashtable()
}

func FreeGlobals() {
	gGlobalEnv.Free()
}

func MarkGlobals() {
	gGlobalEnv.Mark()
}

func SetGlobal(name *ObjString, value Value) bool {
	return gGlobalEnv.Set(name, value)
}

func GetGlobal(name *ObjString) (Value, bool) {
	return gGlobalEnv.Get(name)
}

func DeleteGlobal(name *ObjString) bool {
	return gGlobalEnv.Delete(name)
}
