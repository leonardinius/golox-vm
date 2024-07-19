package vmobject_test

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"

	"github.com/leonardinius/goloxvm/internal/vm/vmobject"
	"github.com/leonardinius/goloxvm/internal/vm/vmvalue"
)

func TestObjValueNanBoxing(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		t.Run("NewObjString", func(t *testing.T) {
			chars1 := []byte("Hello")
			objString := vmobject.NewTakeString(chars1, vmobject.HashString(chars1))
			value := vmvalue.ObjAsValue(objString)
			assert.True(t, vmvalue.IsString(value))
			chars2 := vmvalue.ValueAsStringChars(value)
			assert.Equal(t, "Hello", string(chars2))
			assert.Same(t, unsafe.SliceData(chars1), unsafe.SliceData(chars2))
		})

		t.Run("CopyString", func(t *testing.T) {
			chars1 := []byte("Hello")
			objString := vmobject.NewCopyString(chars1, vmobject.HashString(chars1))
			value := vmvalue.ObjAsValue(objString)
			assert.True(t, vmvalue.IsString(value))
			chars2 := vmvalue.ValueAsStringChars(value)
			assert.Equal(t, "Hello", string(chars2))
			assert.NotSame(t, unsafe.SliceData(chars1), unsafe.SliceData(chars2))
		})
	})
}
