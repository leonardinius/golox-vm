package vmdebug

import (
	"fmt"

	"github.com/leonardinius/goloxvm/internal/vm/vmobject"
	"github.com/leonardinius/goloxvm/internal/vm/vmvalue"
)

func DebugValue(v vmvalue.Value) {
	switch {
	case vmvalue.IsNumber(v):
		fmt.Printf("%g", vmvalue.ValueAsNumber(v))
	case vmvalue.IsNil(v):
		fmt.Print("nil")
	case vmvalue.IsBool(v):
		if vmvalue.ValueAsBool(v) {
			fmt.Print("true")
		} else {
			fmt.Print("false")
		}
	case vmvalue.IsObj(v):
		vmobject.DebugObject(vmvalue.ValueAsObj(v))
	default:
		panic(fmt.Sprintf("unexpected value type: %#v", v))
	}
}

func DebuglnValue(v vmvalue.Value) {
	DebugValue(v)
	fmt.Println()
}
