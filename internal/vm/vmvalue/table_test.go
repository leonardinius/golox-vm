package vmvalue_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonardinius/goloxvm/internal/vm/vmvalue"

	_ "github.com/leonardinius/goloxvm/internal/tests"
)

func TestBasicOps(t *testing.T) {
	h := vmvalue.NewHashtable()
	t.Cleanup(h.Free)

	chars1 := []byte("s1")
	s1 := vmvalue.NewTakeString(chars1, vmvalue.HashString(chars1))
	h.Set(s1, vmvalue.NumberAsValue(10))

	chars2 := []byte("s2")
	s2 := vmvalue.NewTakeString(chars2, vmvalue.HashString(chars2))
	h.Set(s2, vmvalue.NumberAsValue(20))

	// Get
	v, ok := h.Get(s1)
	assert.True(t, ok)
	assert.Equal(t, int(10), int(vmvalue.ValueAsNumber(v)))
	v, ok = h.Get(s2)
	assert.True(t, ok)
	assert.Equal(t, int(20), int(vmvalue.ValueAsNumber(v)))

	// Upset
	h.Set(s1, vmvalue.NumberAsValue(11))
	v, ok = h.Get(s1)
	assert.True(t, ok)
	assert.Equal(t, int(11), int(vmvalue.ValueAsNumber(v)))

	chars3 := []byte("s3")
	s3 := vmvalue.NewTakeString(chars3, vmvalue.HashString(chars3))
	v, ok = h.Get(s3)
	assert.False(t, ok)
	assert.True(t, vmvalue.IsNil(v))
	h.Set(s3, vmvalue.NumberAsValue(30))
	v, ok = h.Get(s3)
	assert.True(t, ok)
	assert.Equal(t, int(30), int(vmvalue.ValueAsNumber(v)))
}

func TestAdjustSize(t *testing.T) {
	h := vmvalue.NewHashtable()
	t.Cleanup(h.Free)

	m := make(map[int]*vmvalue.ObjString)

	for i := range 255 {
		chars := []byte("string" + strconv.Itoa(i))
		s := vmvalue.NewTakeString(chars, vmvalue.HashString(chars))
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

	for e := range 255 {
		i := 255 + e
		chars := []byte("string" + strconv.Itoa(i))
		s := vmvalue.NewTakeString(chars, vmvalue.HashString(chars))
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
