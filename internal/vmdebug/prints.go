package vmdebug

import (
	"fmt"

	"github.com/leonardinius/goloxvm/internal/vmvalue"
)

func PrintValue(v vmvalue.Value) {
	fmt.Printf("%g", vmvalue.ValueAsNumber(v))
}

func PrintlnValue(v vmvalue.Value) {
	PrintValue(v)
	fmt.Println()
}
