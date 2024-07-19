package vmobject_test

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonardinius/goloxvm/internal/vm/vmobject"
)

func TestStringHash(t *testing.T) {
	chars := []byte("Hello")

	hash1 := vmobject.HashString(slices.Clone(chars))
	hash2 := vmobject.HashString(slices.Clone(chars))

	assert.Equal(t, hash1, hash2)
}
