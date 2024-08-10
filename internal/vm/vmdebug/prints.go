package vmdebug

import (
	"fmt"

	"github.com/leonardinius/goloxvm/internal/vm/vmvalue"
)

func PrintValue(v vmvalue.Value) {
	switch {
	case vmvalue.IsNumber(v):
		fv := vmvalue.ValueAsNumber(v)
		if float64(int64(fv)) == fv {
			fmt.Printf("%d", int(fv))
		} else {
			fmt.Printf("%f", fv)
		}
	case vmvalue.IsNil(v):
		fmt.Print("nil")
	case vmvalue.IsBool(v):
		if vmvalue.ValueAsBool(v) {
			fmt.Print("true")
		} else {
			fmt.Print("false")
		}
	case vmvalue.IsObj(v):
		vmvalue.PrintObject(vmvalue.ValueAsObj(v))
	default:
		panic(fmt.Sprintf("unexpected value type: %#v", v))
	}
}

func PrintlnValue(v vmvalue.Value) {
	PrintValue(v)
	fmt.Println()
}
