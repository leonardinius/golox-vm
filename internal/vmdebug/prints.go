package vmdebug

import (
	"fmt"

	"github.com/leonardinius/goloxvm/internal/vmvalue"
)

func PrintValue(v vmvalue.Value) {
	fmt.Printf("%g", v)
}
