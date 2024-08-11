package vmvalue_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonardinius/goloxvm/internal/vm/vmvalue"

	_ "github.com/leonardinius/goloxvm/internal/tests"
)

func TestWriteIncrementsByOne(t *testing.T) {
	va := vmvalue.NewValueArray()
	for constant := range 1024 {
		at := va.Write(vmvalue.NumberAsValue(float64(constant)))
		assert.Equal(t, constant, at)
	}
}
