package vmvalue_test

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"

	"github.com/leonardinius/goloxvm/internal/vm/vmvalue"
)

func TestObjValueNanBoxing(t *testing.T) {
	t.Parallel()
	t.Run("String", func(t *testing.T) {
		t.Run("NewObjString", func(t *testing.T) {
			chars1 := []byte("Hello")
			objString := vmvalue.NewTakeString(chars1, vmvalue.HashString(chars1))
			gc()
			value := vmvalue.ObjAsValue(objString)
			gc()
			assert.True(t, vmvalue.IsString(value))
			chars2 := vmvalue.ValueAsStringChars(value)
			gc()
			assert.Equal(t, "Hello", string(chars2))
			assert.Same(t, unsafe.SliceData(chars1), unsafe.SliceData(chars2))
		})

		t.Run("CopyString", func(t *testing.T) {
			chars1 := []byte("Hello")
			objString := vmvalue.NewCopyString(chars1, vmvalue.HashString(chars1))
			gc()
			value := vmvalue.ObjAsValue(objString)
			gc()
			assert.True(t, vmvalue.IsString(value))
			chars2 := vmvalue.ValueAsStringChars(value)
			gc()
			assert.Equal(t, "Hello", string(chars2))
			assert.NotSame(t, unsafe.SliceData(chars1), unsafe.SliceData(chars2))
		})
	})
}

func gc() {
	// runtime.GC()
}
