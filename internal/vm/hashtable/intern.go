package hashtable

import (
	"github.com/leonardinius/goloxvm/internal/vm/vmobject"
	"github.com/leonardinius/goloxvm/internal/vm/vmvalue"
)

var gInternStrings Table

const marker = vmvalue.NilValue

func InitInternStrings() {
	gInternStrings = NewHashtable()
}

func FreeInternStrings() {
	gInternStrings.Free()
}

func StringInternTake(chars []byte) *vmobject.ObjString {
	hash := vmobject.HashString(chars)

	if str := findString(chars, hash); str != nil {
		return str
	}

	str := vmobject.NewTakeString(chars, hash)
	gInternStrings.Set(str, marker)
	return str
}

func StringInternCopy(chars []byte) *vmobject.ObjString {
	hash := vmobject.HashString(chars)

	if str := findString(chars, hash); str != nil {
		return str
	}

	str := vmobject.NewCopyString(chars, hash)
	gInternStrings.Set(str, marker)
	return str
}

func findString(chars []byte, hash uint64) *vmobject.ObjString {
	return gInternStrings.findString(chars, hash)
}
