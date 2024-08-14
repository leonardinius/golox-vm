package vmvalue

import (
	"bytes"

	"github.com/leonardinius/goloxvm/internal/vm/vmmem"
)

const TableMaxLoad float64 = 0.75

type Table struct {
	entries []entry
	count   int
}

type entry struct {
	key   *ObjString
	value Value
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
	h.entries = vmmem.FreeSlice(h.entries)
	h.reset()
}

func (h *Table) Set(
	key *ObjString,
	value Value,
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

	if isNewKey && IsNil(entry.value) {
		h.count++
	}
	entry.key = key
	entry.value = value
	return isNewKey
}

func (h *Table) findEntry(entries []entry, key *ObjString) *entry {
	capacity := uint64(len(entries))
	debugAssertIsPowerOfTwo(capacity)
	mask := capacity - 1
	index := key.Hash & mask
	var tombstone *entry = nil

	for {
		el := &entries[index]
		if el.key == nil {
			if IsNil(el.value) {
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
	entries := vmmem.GrowSlice(h.entries, capacity)
	for i := range entries {
		el := &entries[i]
		el.key = nil
		el.value = NilValue
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

	h.entries = vmmem.FreeSlice(h.entries)
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

func (h *Table) Get(key *ObjString) (Value, bool) {
	if h.count == 0 {
		return NilValue, false
	}

	if el := h.findEntry(h.entries, key); el.key != nil {
		return el.value, true
	}

	return NilValue, false
}

func (h *Table) Delete(key *ObjString) bool {
	if h.count == 0 {
		return false
	}

	el := h.findEntry(h.entries, key)
	if el.key == nil {
		return false
	}

	el.key = nil
	el.value = BoolAsValue(true)
	// h.count should stay same
	return true
}

func (h *Table) findString(chars []byte, hash uint64) *ObjString {
	if h.count == 0 {
		return nil
	}

	capacity := uint64(len(h.entries))
	debugAssertIsPowerOfTwo(capacity)
	mask := capacity - 1
	index := hash & mask
	for {
		el := &h.entries[index]
		if el.key == nil {
			// Stop if we find an empty non-tombstone entry.
			if IsNil(el.value) {
				return nil
			}
		} else if hash == el.key.Hash &&
			bytes.Equal(chars, el.key.Chars) {
			// we found it
			return el.key
		}
		index = (index + 1) & mask
	}
}

func (h *Table) markTable() {
	for i := range h.entries {
		el := &h.entries[i]
		MarkObject(el.key)
		MarkValue(el.value)
	}
}

func (h *Table) removeWhiteKeys() {
	for i := range h.entries {
		el := &h.entries[i]
		if el.key != nil && !el.key.Marked {
			h.Delete(el.key)
		}
	}
}

func debugAssertIsPowerOfTwo(n uint64) {
	debugAssertf(n != 0 && (n&(n-1)) == 0, "capacity must be power of 2 (%d)", n)
}
