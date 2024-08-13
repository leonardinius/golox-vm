package tests

import (
	"github.com/leonardinius/goloxvm/internal/vm"
	"github.com/leonardinius/goloxvm/internal/vm/vmmem"
	"github.com/leonardinius/goloxvm/internal/vm/vmvalue"
)

func TestInitGarabageCollector() {
	vmmem.SetGarbageCollector(vm.GC)
	vmmem.SetGarbageCollectorRetain(func(v uint64) { vm.Push(vmvalue.NanBoxedAsValue(v)) })
	vmmem.SetGarbageCollectorRelease(func() { _ = vm.Pop() })
}

func init() {
	TestInitGarabageCollector()
}
