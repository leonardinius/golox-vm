package vmobject_test

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"

	"github.com/leonardinius/goloxvm/internal/vmobject"
	"github.com/leonardinius/goloxvm/internal/vmvalue"
)

func TestObjValueNanBoxing(t *testing.T) {
	t.Parallel()

	t.Run("String", func(t *testing.T) {
		t.Run("NewObjString", func(t *testing.T) {
			bytes := []byte("Hello")
			objString := vmobject.NewObjString(bytes)
			value := vmvalue.ObjAsValue(objString)
			assert.True(t, vmvalue.IsString(value))
			chars := vmvalue.ValueAsStringChars(value)
			assert.Equal(t, "Hello", string(chars))
			assert.Same(t, unsafe.SliceData(bytes), unsafe.SliceData(chars))
		})

		t.Run("CopyString", func(t *testing.T) {
			bytes := []byte("Hello")
			objString := vmobject.CopyString(bytes)
			value := vmvalue.ObjAsValue(objString)
			assert.True(t, vmvalue.IsString(value))
			chars := vmvalue.ValueAsStringChars(value)
			assert.Equal(t, "Hello", string(chars))
			assert.NotSame(t, unsafe.SliceData(bytes), unsafe.SliceData(chars))
		})
	})
}
