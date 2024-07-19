package hashtable_test

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonardinius/goloxvm/internal/vm/hashtable"
)

func TestStringInternTake(t *testing.T) {
	hashtable.InitInternStrings()
	t.Cleanup(hashtable.FreeInternStrings)

	chars := []byte("Hello")
	s1 := hashtable.StringInternTake(slices.Clone(chars))
	s2 := hashtable.StringInternTake(slices.Clone(chars))
	s3 := hashtable.StringInternTake(slices.Clone(chars))
	assert.Same(t, s1, s2)
	assert.Same(t, s2, s3)
}

func TestStringInternCopy(t *testing.T) {
	hashtable.InitInternStrings()
	t.Cleanup(hashtable.FreeInternStrings)

	chars := []byte("Hello")
	s1 := hashtable.StringInternCopy(chars)
	s2 := hashtable.StringInternCopy(chars)
	s3 := hashtable.StringInternCopy(chars)
	assert.Same(t, s1, s2)
	assert.Same(t, s2, s3)
}
