package vmchunk

import (
	"fmt"

	"github.com/leonardinius/goloxvm/internal/vm/vmmem"
)

type Lines struct {
	raw    []byte
	index  int
	offset int
}

func (l *Lines) Init() {
	l.raw = nil
	l.index = -1
	l.offset = -1
}

func (l *Lines) Free() {
	l.Init()
}

func (l *Lines) GetLineByOffset(offset int) int {
	if offset > l.offset || offset < 0 {
		return -1
	}

	running := 0
	for i := range l.count() {
		pos := l.wrapLineFromBytes(i)
		count := pos.count()

		running += count
		if running >= offset {
			return pos.line()
		}
	}

	return -1
}

func (l *Lines) MustWriteOffset(offset, line int) {
	if offset <= l.offset {
		panic("offset must be non-decreasing")
	}

	if l.isEmpty() {
		l.index = 0
		l.offset = offset
		l.ensureCapacity(l.index)
		l.putLineToBytes(wrapLineWithCount(l.offset, line), l.index)
		return
	}

	pos := l.wrapLineFromBytes(l.index)
	countDiff := offset - l.offset
	if pos.line() == line && pos.count()+countDiff < 255 {
		l.setCount(l.index, pos.count()+countDiff)
		l.offset = offset
		return
	}

	for countDiff > 0 {
		l.ensureCapacity(l.index + 1)
		count := max(1, min(255, countDiff))
		l.putLineToBytes(wrapLineWithCount(count, line), l.index+1)
		countDiff -= count
		l.index++
		l.offset += count
	}

	if l.offset != offset {
		panic(fmt.Sprintf("offset invariant does not match %d != %d", l.offset, offset))
	}
}

func (l *Lines) ensureCapacity(lineIndex int) {
	if cap(l.raw) < (lineIndex+1)*3 {
		capacity := vmmem.GrowCapacity(len(l.raw))
		l.raw = vmmem.GrowSlice(l.raw, capacity)
		l.raw = l.raw[:cap(l.raw)]
	}
}

func (l *Lines) isEmpty() bool {
	return len(l.raw) == 0
}

func (l *Lines) count() int {
	return len(l.raw) / 3
}

func (l *Lines) wrapLineFromBytes(index int) linepos {
	i := index * 3
	return linepos((uint32(l.raw[i]) << 16) | (uint32(l.raw[i+1]) << 8) | uint32(l.raw[i+2]))
}

func (l *Lines) setCount(index, count int) {
	l.raw[index*3] = byte(count)
}

func (l *Lines) putLineToBytes(line linepos, index int) {
	// 3 bytes per line
	u := uint32(line)
	i := index * 3
	l.raw[i] = byte((u >> 16) & 0xFF)
	l.raw[i+1] = byte((u >> 8) & 0xFF)
	l.raw[i+2] = byte((u) & 0xFF)
}

// line[i, i+1, i+2] is a run encoded line information
// [3 bytes] { [count] [2 bytes for line number] }
//
// the sum of all counts is the current offset.
type linepos uint32

func wrapLineWithCount(count, line int) linepos {
	return linepos((count << 16) | (line & 0xFFFF))
}

func (p linepos) count() int {
	return int((uint32(p) >> 16) & 0xFF)
}

func (p linepos) line() int {
	return int(uint32(p) & 0xFFFF)
}
