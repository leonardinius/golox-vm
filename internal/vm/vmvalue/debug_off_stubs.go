//go:build !debug

package vmvalue

func debugAssertf(condition bool, message string, args ...any) {}

func debugPrintFreeObject(header *Obj, size int) {}

func debugPrintAllocateObject(header *Obj, size int) {}

func debugPrintMarkObject(obj *Obj) {}

func debugPrintBlackenObject(obj *Obj) {}
