package vmmem_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonardinius/goloxvm/internal/vm/vmmem"
)

func TestGrowArrayShouldGrowCapacity(t *testing.T) {
	a := make([]int, 0)
	a = vmmem.GrowSlice(a, 10)
	assert.Len(t, a, 10)
	assert.GreaterOrEqual(t, 10, cap(a))
}
