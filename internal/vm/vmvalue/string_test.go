package vmvalue_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonardinius/goloxvm/internal/vm/vmvalue"
)

func TestStringHash(t *testing.T) {
	chars := []byte("Hello")

	hash1 := vmvalue.HashString(bytes.Clone(chars))
	hash2 := vmvalue.HashString(bytes.Clone(chars))

	assert.Equal(t, hash1, hash2)
}
