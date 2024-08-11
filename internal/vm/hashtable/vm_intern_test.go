package hashtable_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonardinius/goloxvm/internal/vm/hashtable"

	_ "github.com/leonardinius/goloxvm/internal/tests"
)

func TestStringInternTake(t *testing.T) {
	hashtable.InitInternStrings()
	t.Cleanup(hashtable.FreeInternStrings)

	chars := []byte("Hello")
	s1 := hashtable.StringInternTake(bytes.Clone(chars))
	s2 := hashtable.StringInternTake(bytes.Clone(chars))
	s3 := hashtable.StringInternTake(bytes.Clone(chars))
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
