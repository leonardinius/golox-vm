package vmvalue

import "github.com/leonardinius/goloxvm/internal/vm/vmmem"

type ValueArray []Value

func (va *ValueArray) At(i int) Value {
	return (*va)[i]
}

func (va *ValueArray) Write(v Value) int {
	if cap(*va) < len(*va)+1 {
		capacity := vmmem.GrowCapacity(cap(*va))
		*va = vmmem.GrowArray(*va, capacity)
	}
	*va = append(*va, v)
	return len(*va) - 1
}

func InitValueArray(va *ValueArray) {
	*va = nil
}

func FreeValueArray(va *ValueArray) {
	*va = vmmem.FreeArray(*va)
	InitValueArray(va)
}
