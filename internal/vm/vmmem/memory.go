package vmmem

// #include <stdio.h>
// #include <stdlib.h>
// #include <errno.h>
import "C"      //nolint:gocritic // dupImport
import "unsafe" //nolint:gocritic // dupImport

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

func CMalloc(length int) unsafe.Pointer {
	// force C compiler to allocate memory
	return C.malloc(C.ulong(length))
}

func CFree[T any](ptr *T) {
	C.free(unsafe.Pointer(ptr))
}

func CMallocBytes(length int) []byte {
	return unsafe.Slice((*byte)(CMalloc(length)), length)
}

func CFreeBytes(bytes []byte) {
	// force C compiler to free memory
	CFree(unsafe.SliceData(bytes))
}
