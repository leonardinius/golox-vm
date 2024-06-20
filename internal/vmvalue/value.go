package vmvalue

import "github.com/leonardinius/goloxvm/internal/vmmem"

type Value float64

type ValueArray []Value

func (va *ValueArray) Count() int {
	return len(*va)
}

func (va *ValueArray) Cap() int {
	return cap(*va)
}

func (va *ValueArray) At(i int) Value {
	return (*va)[i]
}

func (va *ValueArray) Write(v Value) int {
	if va.Cap() < va.Count()+1 {
		capacity := vmmem.GrowCapacity(va.Cap())
		*va = vmmem.GrowArray(*va, capacity)
	}
	*va = append(*va, v)
	return va.Count() - 1
}

func InitValueArray(va *ValueArray) {
	*va = nil
}

func FreeValueArray(va *ValueArray) {
	*va = vmmem.FreeArray(*va)
	InitValueArray(va)
}
