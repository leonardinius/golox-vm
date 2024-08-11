package tests

import (
	"github.com/leonardinius/goloxvm/internal/vm/vmmem"
)

func TestInitGarabageCollector() {
	vmmem.SetGarbageCollector(func() {})
}

func init() {
	TestInitGarabageCollector()
}
