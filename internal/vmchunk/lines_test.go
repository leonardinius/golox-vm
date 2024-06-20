package vmchunk_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonardinius/goloxvm/internal/vmchunk"
)

func TestSetOffsetShouldValidateInput(t *testing.T) {
	lines := vmchunk.Lines{}
	lines.Init()

	lines.MustWriteOffset(1, 1)
	lines.MustWriteOffset(2, 1)
	assert.PanicsWithValue(t, "offset must be non-decreasing", func() { lines.MustWriteOffset(2, 1) })
}

func TestGetLine404ShouldFailGracefully(t *testing.T) {
	lines := vmchunk.Lines{}
	lines.Init()
	assert.Equal(t, -1, lines.GetLineByOffset(-1))
	assert.Equal(t, -1, lines.GetLineByOffset(0))
	assert.Equal(t, -1, lines.GetLineByOffset(1))

	lines.MustWriteOffset(0, 1)
	assert.Equal(t, -1, lines.GetLineByOffset(-1))
	assert.Equal(t, 1, lines.GetLineByOffset(0))
	assert.Equal(t, -1, lines.GetLineByOffset(1))

	lines.MustWriteOffset(1, 1)
	assert.Equal(t, 1, lines.GetLineByOffset(0))
	assert.Equal(t, 1, lines.GetLineByOffset(1))
	assert.Equal(t, -1, lines.GetLineByOffset(2))
}

func TestEncodeDecodeLinesInformation(t *testing.T) {
	lines := vmchunk.Lines{}
	lines.Init()

	lines.MustWriteOffset(0, 1)
	lines.MustWriteOffset(1, 2)
	lines.MustWriteOffset(2, 3)
	assert.Equal(t, 1, lines.GetLineByOffset(0))
	assert.Equal(t, 2, lines.GetLineByOffset(1))
	assert.Equal(t, 3, lines.GetLineByOffset(2))

	lines.MustWriteOffset(3, 3)
	lines.MustWriteOffset(4, 3)
	lines.MustWriteOffset(5, 3)
	lines.MustWriteOffset(6, 3)
	lines.MustWriteOffset(7, 3)
	lines.MustWriteOffset(8, 3)
	lines.MustWriteOffset(9, 3)
	assert.Equal(t, 3, lines.GetLineByOffset(3))
	assert.Equal(t, 3, lines.GetLineByOffset(9))

	lines.MustWriteOffset(265, 5)
	lines.MustWriteOffset(266, 6)
	lines.MustWriteOffset(267, 7)
	assert.Equal(t, 5, lines.GetLineByOffset(265))
	assert.Equal(t, 6, lines.GetLineByOffset(266))
	assert.Equal(t, 7, lines.GetLineByOffset(267))

	lines.MustWriteOffset(1024, 8)
	lines.MustWriteOffset(1025, 8)
	lines.MustWriteOffset(1026, 8)
	lines.MustWriteOffset(1027, 9)
	lines.MustWriteOffset(1028, 10)
	assert.Equal(t, 8, lines.GetLineByOffset(1024))
	assert.Equal(t, 8, lines.GetLineByOffset(1025))
	assert.Equal(t, 8, lines.GetLineByOffset(1026))
	assert.Equal(t, 9, lines.GetLineByOffset(1027))
	assert.Equal(t, 10, lines.GetLineByOffset(1028))

	lines.Free()
}
