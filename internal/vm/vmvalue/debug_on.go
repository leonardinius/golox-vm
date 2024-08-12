//go:build debug

package vmvalue

import (
	"fmt"
)

func debugAssertf(condition bool, message string, args ...any) {
	if !condition {
		panic(fmt.Errorf(message, args...))
	}
}

func debugPrintf(message string, args ...any) {
	fmt.Printf(message, args...)
}

func debugPrintFreeObject(obj *Obj, size int) {
	debugAssertf(obj != nil, "debugPrintFreeObject: o is nil %v", obj)
	fmt.Printf("%p %-7s %04d for [%-12s] '", obj, "free", size, obj.Type)
	PrintObject(obj)
	fmt.Println("'")
}

func debugPrintAllocateObject(obj *Obj, size int) {
	debugAssertf(obj != nil, "debugPrintAllocateObject: o is nil %v", obj)
	fmt.Printf("%p %-7s %04d for [%-12s]\n", obj, "allocate", size, obj.Type)
}

func debugPrintMarkObject(obj *Obj) {
	debugAssertf(obj != nil, "debugPrintMarkObject: o is nil %v", obj)
	fmt.Printf("%p %-7s [%-12s] '", obj, "mark", obj.Type)
	PrintObject(obj)
	fmt.Println("'")
}

func debugPrintSkipMarkedObject(obj *Obj) {
	debugAssertf(obj != nil, "debugPrintMarkObject: o is nil %v", obj)
	fmt.Printf("%p %-7s [%-12s] '", obj, "skip", obj.Type)
	PrintObject(obj)
	fmt.Println("'")
}

func debugPrintBlackenObject(obj *Obj) {
	debugAssertf(obj != nil, "debugPrintBlackenObject: o is nil %v", obj)
	fmt.Printf("%p %-7s [%-12s] '", obj, "blacken", obj.Type)
	PrintObject(obj)
	fmt.Println("'")
}
