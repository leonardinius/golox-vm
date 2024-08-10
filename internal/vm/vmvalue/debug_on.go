//go:build debug

package vmvalue

import "fmt"

func DebugFreeObject(obj *Obj, message string) {
	fmt.Printf("%p %-7s [%-12s] ", obj, message, obj.Type)
	PrintObject(obj)
	fmt.Println()
}
