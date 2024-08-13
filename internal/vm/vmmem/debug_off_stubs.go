//go:build !debug

package vmmem

const DebugGC = false

func debugPrintln(message string) {
}

func debugPrintlf(message string, args ...any) {
}

func debugStressGC() {}
