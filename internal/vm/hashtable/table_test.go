package hashtable_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonardinius/goloxvm/internal/vm/hashtable"
	"github.com/leonardinius/goloxvm/internal/vm/vmobject"
	"github.com/leonardinius/goloxvm/internal/vm/vmvalue"
)

func TestBasicOps(t *testing.T) {
	h := hashtable.NewHashtable()
	t.Cleanup(h.Free)
	assert.Equal(t, 0, h.Count())

	chars1 := []byte("s1")
	s1 := vmobject.NewTakeString(chars1, vmobject.HashString(chars1))
	h.Set(s1, vmvalue.NumberAsValue(10))
	assert.Equal(t, 1, h.Count())

	chars2 := []byte("s2")
	s2 := vmobject.NewTakeString(chars2, vmobject.HashString(chars2))
	h.Set(s2, vmvalue.NumberAsValue(20))
	assert.Equal(t, 2, h.Count())

	// Get
	v, ok := h.Get(s1)
	assert.True(t, ok)
	assert.Equal(t, int(10), int(vmvalue.ValueAsNumber(v)))
	v, ok = h.Get(s2)
	assert.True(t, ok)
	assert.Equal(t, int(20), int(vmvalue.ValueAsNumber(v)))

	// Upset
	h.Set(s1, vmvalue.NumberAsValue(11))
	assert.Equal(t, 2, h.Count())
	v, ok = h.Get(s1)
	assert.True(t, ok)
	assert.Equal(t, int(11), int(vmvalue.ValueAsNumber(v)))

	chars3 := []byte("s3")
	s3 := vmobject.NewTakeString(chars3, vmobject.HashString(chars3))
	v, ok = h.Get(s3)
	assert.False(t, ok)
	assert.True(t, vmvalue.IsNil(v))
	h.Set(s3, vmvalue.NumberAsValue(30))
	assert.Equal(t, 3, h.Count())
	v, ok = h.Get(s3)
	assert.True(t, ok)
	assert.Equal(t, int(30), int(vmvalue.ValueAsNumber(v)))
}

func TestAdjustSize(t *testing.T) {
	h := hashtable.NewHashtable()
	t.Cleanup(h.Free)

	m := make(map[int]*vmobject.ObjString)

	for i := range 255 {
		chars := []byte("string" + strconv.Itoa(i))
		s := vmobject.NewTakeString(chars, vmobject.HashString(chars))
		h.Set(s, vmvalue.NumberAsValue(float64(i)))
		m[i] = s
	}

	for i := range 255 {
		if i%5 != 0 {
			continue
		}
		s := m[i]
		ok := h.Delete(s)
		assert.True(t, ok)
	}
	assert.Equal(t, 255-51, h.Count())

	for e := range 255 {
		i := 255 + e
		chars := []byte("string" + strconv.Itoa(i))
		s := vmobject.NewTakeString(chars, vmobject.HashString(chars))
		h.Set(s, vmvalue.NumberAsValue(float64(i)))
		m[i] = s
	}

	for i := range 255 {
		if i%5 != 0 {
			continue
		}
		s := m[255+i]
		ok := h.Delete(s)
		assert.True(t, ok)
	}
	assert.Equal(t, 255-51+255-51, h.Count())

	for i := range 255 + 255 {
		if i%5 != 0 {
			_, ok := h.Get(m[i])
			assert.True(t, ok)
		} else {
			_, ok := h.Get(m[i])
			assert.False(t, ok)
		}
	}
}
