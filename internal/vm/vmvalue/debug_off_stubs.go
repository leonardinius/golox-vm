//go:build !debug

package vmvalue

func debugPrintFreeObject(header *Obj, size int) {}

func debugPrintAllocateObject(header *Obj, size int) {}

func debugPrintMarkObject(obj *Obj) {}

func debugPrintSkipMarkedObject(obj *Obj) {}

func debugPrintBlackenObject(obj *Obj) {}
