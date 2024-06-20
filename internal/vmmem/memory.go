package vmmem

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
		s = append(s[:cap(s)], make([]E, newSize-oldSize)...)
	} else if newSize < oldSize {
		s = s[:newSize]
	}
	return s
}
