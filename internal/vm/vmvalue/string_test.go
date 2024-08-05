package vmvalue_test

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonardinius/goloxvm/internal/vm/vmvalue"
)

func TestStringHash(t *testing.T) {
	chars := []byte("Hello")

	hash1 := vmvalue.HashString(slices.Clone(chars))
	hash2 := vmvalue.HashString(slices.Clone(chars))

	assert.Equal(t, hash1, hash2)
}
