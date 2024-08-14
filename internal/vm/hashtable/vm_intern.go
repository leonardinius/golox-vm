package hashtable

import (
	"github.com/leonardinius/goloxvm/internal/vm/vmmem"
	"github.com/leonardinius/goloxvm/internal/vm/vmvalue"
)

var gInternStrings Table

const internMarkerValue = vmvalue.NilValue

func InitInternStrings() {
	gInternStrings = NewHashtable()
}

func FreeInternStrings() {
	gInternStrings.Free()
}

func StringInternTake(chars []byte) *vmvalue.ObjString {
	hash := vmvalue.HashString(chars)

	if str := findString(chars, hash); str != nil {
		return str
	}

	str := vmvalue.NewTakeString(chars, hash)
	vmmem.PushRetainGC(uint64(vmvalue.ObjAsValue(str)))
	defer vmmem.PopReleaseGC()
	gInternStrings.Set(str, internMarkerValue)
	return str
}

func StringInternCopy(chars []byte) *vmvalue.ObjString {
	hash := vmvalue.HashString(chars)

	if str := findString(chars, hash); str != nil {
		return str
	}

	str := vmvalue.NewCopyString(chars, hash)
	vmmem.PushRetainGC(vmvalue.ValueAsNanBoxed(vmvalue.ObjAsValue(str)))
	defer vmmem.PopReleaseGC()
	gInternStrings.Set(str, internMarkerValue)
	return str
}

func findString(chars []byte, hash uint64) *vmvalue.ObjString {
	return gInternStrings.findString(chars, hash)
}

func RemoveWhiteInternStrings() {
	gInternStrings.removeWhiteKeys()
}
