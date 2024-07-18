package vmmem

import "unsafe"

func GrowCapacity(n int) int {
	if n < 8 {
		return 8
	}
	return n * 2
}

func GrowArray[S ~[]E, E any](s S, n int) S {
	return Reallocate(s, cap(s), n)
}

func FreeArray[S ~[]E, E any](s S) S {
	return Reallocate(s, cap(s), 0)
}

func Reallocate[S ~[]E, E any](s S, oldSize, newSize int) S {
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

func AllocateSlice[E any](length int) []E {
	var empty []E = nil
	return Reallocate(empty, 0, length)
}

func AllocateUnsafePtr[E any](length int) unsafe.Pointer {
	slice := AllocateSlice[E](length)
	return unsafe.Pointer(unsafe.SliceData(slice)) //nolint:gosec // return unsafe Pointer here
}

func FreeUnsafePtr[E any, T any](ptr *T, length int) *T {
	data := (*E)(unsafe.Pointer(ptr))   //nolint:gosec //  unsafe Pointer here
	slice := unsafe.Slice(data, length) //nolint:gosec //  unsafe Pointer here
	FreeArray(slice)
	return nil
}
