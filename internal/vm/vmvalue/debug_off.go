//go:build !debug

package vmvalue

func debugFreeObject(header *Obj, size int, message string) {}

func debugAllocateObject(header *Obj, size int, message string) {}

func debugMarkObject(obj *Obj) {}

func debugBlackenObject(obj *Obj) {}

func debugStressGC() {}
