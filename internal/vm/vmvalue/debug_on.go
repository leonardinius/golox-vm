//go:build debug

package vmvalue

import (
	"fmt"

	"github.com/leonardinius/goloxvm/internal/vm/vmmem"
)

func debugFreeObject(header *Obj, size int, message string) {
	fmt.Printf("%p %-7s %-4d [%-12s] ", header, message, size, header.Type)
	PrintObject(header)
	fmt.Println()
}

func debugAllocateObject(header *Obj, size int, message string) {
	fmt.Printf("%p %-7s %-4d [%-12s]\n", header, message, size, header.Type)
}

func debugMarkObject(obj *Obj) {
	fmt.Printf("%p mark '", obj)
	PrintObject(obj)
	fmt.Println("'")
}

func debugBlackenObject(obj *Obj) {
	fmt.Printf("%p blacken '", obj)
	PrintObject(obj)
	fmt.Println("'")
}

func debugStressGC() {
	vmmem.DebugStressGC()
}
