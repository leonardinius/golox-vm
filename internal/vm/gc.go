package vm

import (
	"github.com/leonardinius/goloxvm/internal/vm/vmvalue"
	"github.com/leonardinius/goloxvm/internal/vmcompiler"
)

func GC() {
	markRoots()
	traceReferences()
	tableRemoveWhiteInternStrings()
	sweep()
}

func markRoots() {
	for i := range GlobalVM.StackTop {
		vmvalue.MarkValue(StackAt(i))
	}

	for i := range GlobalVM.FrameCount {
		vmvalue.MarkObject(GlobalVM.Frames[i].Closure)
	}

	for upvalue := GlobalVM.OpenUpvalues; upvalue != nil; upvalue = upvalue.Next {
		vmvalue.MarkObject(upvalue)
	}

	vmvalue.MarkGlobals()

	vmcompiler.MarkCompilerRoots()

	vmvalue.MarkObject(GlobalVM.InitString)
}

func traceReferences() {
	vmvalue.GCTraceReferences()
}

func tableRemoveWhiteInternStrings() {
	vmvalue.RemoveWhiteInternStrings()
}

func sweep() {
	vmvalue.GCSweep()
}
