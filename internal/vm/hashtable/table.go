package hashtable

import (
	"slices"

	"github.com/leonardinius/goloxvm/internal/vm/vmdebug"
	"github.com/leonardinius/goloxvm/internal/vm/vmmem"
	"github.com/leonardinius/goloxvm/internal/vm/vmvalue"
)

const TableMaxLoad float64 = 0.75

type Table struct {
	entries []entry
	count   int
}

type entry struct {
	key   *vmvalue.ObjString
	value vmvalue.Value
}

func NewHashtable() Table {
	h := Table{}
	h.reset()
	return h
}

func (h *Table) reset() {
	h.entries = nil
	h.count = 0
}

func (h *Table) Free() {
	h.entries = vmmem.FreeArray(h.entries)
	h.reset()
}

func (h *Table) Count() int {
	return h.count
}

func (h *Table) Set(
	key *vmvalue.ObjString,
	value vmvalue.Value,
) bool {
	loadLimit := int(float64(len(h.entries)) * TableMaxLoad)
	if h.count+1 > loadLimit {
		capacity := vmmem.GrowCapacity(len(h.entries))
		h.adjustCapacity(capacity)
	}

	entry := h.findEntry(h.entries, key)
	var isNewKey bool
	if entry.key == nil {
		isNewKey = true
	}

	if isNewKey && vmvalue.IsNil(entry.value) {
		h.count++
	}
	entry.key = key
	entry.value = value
	return isNewKey
}

func (h *Table) findEntry(entries []entry, key *vmvalue.ObjString) *entry {
	capacity := uint64(len(entries))
	vmdebug.Assertf(capacity%2 == 0, "capacity must be greater than 0 (%d)", capacity)
	mask := capacity - 1
	index := key.Hash & mask
	var tombstone *entry = nil
	for {
		el := &entries[index]
		if el.key == nil {
			if vmvalue.IsNil(el.value) {
				if tombstone != nil {
					return tombstone
				}
				return el
			} else if tombstone == nil {
				tombstone = el
			}
		} else if el.key == key {
			return el
		}
		index = (index + 1) & mask
	}
}

func (h *Table) adjustCapacity(capacity int) {
	entries := vmmem.GrowArray(h.entries, capacity)
	for i := range entries {
		el := &entries[i]
		el.key = nil
		el.value = vmvalue.NilValue
	}

	h.count = 0
	for i := range h.entries {
		el := &h.entries[i]
		if el.key == nil {
			continue
		}

		dest := h.findEntry(entries, el.key)
		dest.key = el.key
		dest.value = el.value
		h.count++
	}

	h.entries = vmmem.FreeArray(h.entries)
	h.entries = entries
}

func (h *Table) PutAll(from *Table) {
	for i := range from.entries {
		el := &from.entries[i]
		if el.key != nil {
			h.Set(el.key, el.value)
		}
	}
}

func (h *Table) Get(key *vmvalue.ObjString) (vmvalue.Value, bool) {
	if h.count == 0 {
		return vmvalue.NilValue, false
	}

	if el := h.findEntry(h.entries, key); el.key != nil {
		return el.value, true
	}

	return vmvalue.NilValue, false
}

func (h *Table) Delete(key *vmvalue.ObjString) bool {
	if h.count == 0 {
		return false
	}

	el := h.findEntry(h.entries, key)
	if el.key == nil {
		return false
	}

	el.key = nil
	el.value = vmvalue.BoolAsValue(true)
	h.count--
	return true
}

func (h *Table) findString(chars []byte, hash uint64) *vmvalue.ObjString {
	if h.count == 0 {
		return nil
	}

	capacity := uint64(len(h.entries))
	vmdebug.Assertf(capacity%2 == 0, "capacity must be greater than 0 (%d)", capacity)
	mask := capacity - 1
	index := hash & mask
	for {
		el := &h.entries[index]
		if el.key == nil && vmvalue.IsNil(el.value) {
			return nil
		} else if hash == el.key.Hash && slices.Equal(chars, el.key.Chars) {
			return el.key
		}
		index = (index + 1) & mask
	}
}
