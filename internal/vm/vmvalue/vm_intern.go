package vmvalue

import (
	"github.com/leonardinius/goloxvm/internal/vm/vmmem"
)

var gInternStrings Table

const internMarkerValue = NilValue

func InitInternStrings() {
	gInternStrings = NewHashtable()
}

func FreeInternStrings() {
	gInternStrings.Free()
}

func StringInternTake(chars []byte) *ObjString {
	hash := HashString(chars)

	if str := findString(chars, hash); str != nil {
		return str
	}

	str := NewTakeString(chars, hash)
	vmmem.PushRetainGC(uint64(ObjAsValue(str)))
	defer vmmem.PopReleaseGC()
	gInternStrings.Set(str, internMarkerValue)
	return str
}

func StringInternCopy(chars []byte) *ObjString {
	hash := HashString(chars)

	if str := findString(chars, hash); str != nil {
		return str
	}

	str := NewCopyString(chars, hash)
	vmmem.PushRetainGC(ValueAsNanBoxed(ObjAsValue(str)))
	defer vmmem.PopReleaseGC()
	gInternStrings.Set(str, internMarkerValue)
	return str
}

func findString(chars []byte, hash uint64) *ObjString {
	return gInternStrings.findString(chars, hash)
}

func RemoveWhiteInternStrings() {
	gInternStrings.removeWhiteKeys()
}
