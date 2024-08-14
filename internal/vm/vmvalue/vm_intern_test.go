package vmvalue_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	_ "github.com/leonardinius/goloxvm/internal/tests"
	"github.com/leonardinius/goloxvm/internal/vm/vmvalue"
)

func TestStringInternTake(t *testing.T) {
	vmvalue.InitInternStrings()
	t.Cleanup(vmvalue.FreeInternStrings)

	chars := []byte("Hello")
	s1 := vmvalue.StringInternTake(bytes.Clone(chars))
	s2 := vmvalue.StringInternTake(bytes.Clone(chars))
	s3 := vmvalue.StringInternTake(bytes.Clone(chars))
	assert.Same(t, s1, s2)
	assert.Same(t, s2, s3)
}

func TestStringInternCopy(t *testing.T) {
	vmvalue.InitInternStrings()
	t.Cleanup(vmvalue.FreeInternStrings)

	chars := []byte("Hello")
	s1 := vmvalue.StringInternCopy(chars)
	s2 := vmvalue.StringInternCopy(chars)
	s3 := vmvalue.StringInternCopy(chars)
	assert.Same(t, s1, s2)
	assert.Same(t, s2, s3)
}
