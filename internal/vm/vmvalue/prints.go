package vmvalue

import (
	"fmt"
)

func PrintValue(v Value) {
	switch {
	case IsNumber(v):
		fv := ValueAsNumber(v)
		fmt.Printf("%v", fv)
	case IsNil(v):
		fmt.Print("nil")
	case IsBool(v):
		if ValueAsBool(v) {
			fmt.Print("true")
		} else {
			fmt.Print("false")
		}
	case IsObj(v):
		PrintObject(ValueAsObj(v))
	default:
		panic(fmt.Sprintf("unexpected value type: %#v", v))
	}
}

func PrintlnValue(v Value) {
	PrintValue(v)
	fmt.Println()
}
