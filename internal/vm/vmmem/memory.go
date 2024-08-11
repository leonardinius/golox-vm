package vmmem

import "unsafe"

func GrowCapacity(n int) int {
	if n < 8 {
		return 8
	}
	return n * 2
}

func GrowSlice[S ~[]E, E any](s S, n int) S {
	return ReallocateSlice(s, cap(s), n)
}

func FreeSlice[S ~[]E, E any](s S) S {
	return ReallocateSlice(s, cap(s), 0)
}

func ReallocateSlice[S ~[]E, E any](s S, oldSize, newSize int) S {
	var v E
	TriggerGC(int(unsafe.Sizeof(v)), oldSize, newSize)

	if newSize == 0 {
		s = nil
	} else if newSize > oldSize {
		// modification of slices.Grow
		s = append(s[:cap(s)], make([]E, newSize-oldSize)...)[:newSize]
	} else if newSize < oldSize {
		s = s[:newSize]
	}
	return s
}

func AllocateSlice[E any](size int) []E {
	var slice []E
	return ReallocateSlice(slice, 0, size)
}

type memgc struct {
	collect        func()
	bytesAllocated int
	nextGC         int
}

var gc = memgc{}

const gcHeapGrowFactor = 2

func SetGarbageCollector(f func()) {
	// Proper way is to go with atomic.Value
	// but! to not overkill in toy project
	// I'll just use this hacky way
	gc.collect = f
	gc.bytesAllocated = 0
	gc.nextGC = 1024 * 1024
}

func TriggerGC(elemSize, oldSize, newSize int) {
	debugStressGC()
	gc.bytesAllocated += ((elemSize * newSize) - (elemSize * oldSize))
	if newSize > oldSize && gc.bytesAllocated > gc.nextGC {
		CollectGarbage()
	}
}

func CollectGarbage() {
	debugPrintln("-- gc begin")
	before := gc.bytesAllocated
	gc.collect()
	if before > gc.nextGC {
		gc.nextGC = gc.bytesAllocated * gcHeapGrowFactor
	}
	after := gc.bytesAllocated
	debugPrintln("-- gc end")
	debugPrintlf("   collected %d bytes (from %d to %d) next at %d", before-after, before, after, gc.nextGC)
}
