package vm

import (
	"github.com/leonardinius/goloxvm/internal/vm/hashtable"
	"github.com/leonardinius/goloxvm/internal/vm/vmdebug"
	"github.com/leonardinius/goloxvm/internal/vm/vmvalue"
	"github.com/leonardinius/goloxvm/internal/vmcompiler"
)

func GC() {
	vmdebug.Printf("--  >> mark roots\n")
	markRoots()
	vmdebug.Printf("--  >> trace references\n")
	traceReferences()
	vmdebug.Printf("--  >> remove unused intern strings\n")
	tableRemoveWhiteInternStrings()
	vmdebug.Printf("--  >> sweep\n")
	sweep()
}

func markRoots() {
	vmdebug.Printf("--  >> | stack\n")
	for i := range GlobalVM.StackTop {
		vmvalue.MarkValue(StackAt(i))
	}

	vmdebug.Printf("--  >> | call frames\n")
	for i := range GlobalVM.FrameCount {
		vmvalue.MarkObject(GlobalVM.Frames[i].Closure)
	}

	vmdebug.Printf("--  >> | upvalues\n")
	for upvalue := GlobalVM.OpenUpvalues; upvalue != nil; upvalue = upvalue.Next {
		vmvalue.MarkObject(upvalue)
	}

	vmdebug.Printf("--  >> | globals\n")
	hashtable.MarkGlobals()

	vmdebug.Printf("--  >> | compiler roots\n")
	vmcompiler.MarkCompilerRoots()
}

func traceReferences() {
	vmvalue.GCTraceReferences()
}

func tableRemoveWhiteInternStrings() {
	hashtable.RemoveWhiteInternStrings()
}

func sweep() {
	vmvalue.GCSweep()
}
