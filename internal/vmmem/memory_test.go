package vmmem_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonardinius/goloxvm/internal/vmmem"
)

func TestGrowArrayShouldAllocateZeroMemory(t *testing.T) {
	a := make([]int, 0)
	a = vmmem.GrowArray(a, 10)
	a[9] = 1
	assert.Equal(t, 0, a[8])
	assert.Equal(t, 1, a[9])

	a = vmmem.GrowArray(a, 255)
	a[254] = 2
	assert.Equal(t, 0, a[253])
	assert.Equal(t, 2, a[254])
	assert.Len(t, a, 255)
	assert.Less(t, 255, cap(a))

	a = vmmem.GrowArray(a, 513)
	a[512] = 3
	assert.Equal(t, 0, a[511])
	assert.Equal(t, 3, a[512])
	assert.Len(t, a, 513)
	assert.Less(t, 513, cap(a))
}
