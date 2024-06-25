package vmmem_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonardinius/goloxvm/internal/vmmem"
)

func TestGrowArrayShouldGrowCapacity(t *testing.T) {
	a := make([]int, 0)
	a = vmmem.GrowArray(a, 10)
	assert.Empty(t, a)
	assert.GreaterOrEqual(t, 10, cap(a))
}
