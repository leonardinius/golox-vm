package vmvalue

import (
	"slices"
	"strconv"

	"github.com/leonardinius/goloxvm/internal/vm/vmmem"
)

type ValueArray []Value

func NewValueArray() ValueArray {
	return nil
}

func (va *ValueArray) Init() {
	*va = nil
}

func (va *ValueArray) Free() {
	*va = vmmem.FreeSlice(*va)
}

func (va *ValueArray) Mark() {
	for i := range *va {
		MarkValue((*va)[i])
	}
	debugPrintf("-- start of constants\n")
	for i := range *va {
		v := (*va)[i]
		debugPrintf("   >> constant [%d] = '", i)
		PrintValue(v)
		debugPrintf("'")
		if IsObj(v) {
			obj := ValueAsObj(v)
			debugPrintf(" %s", strconv.FormatBool(obj.Marked))
		}
		debugPrintf("\n")
	}
	debugPrintf("-- end of constants\n")
}

func (va *ValueArray) At(i int) Value {
	return (*va)[i]
}

func (va *ValueArray) Write(v Value) int {
	length := len(*va)

	if cap(*va) < len(*va)+1 {
		capacity := vmmem.GrowCapacity(cap(*va))
		*va = slices.Grow(*va, capacity)
		vaarray := *va
		*va = vaarray[0:length:capacity]
	}

	*va = append(*va, v)
	return len(*va) - 1
}
