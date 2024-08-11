//go:build debug

package vmmem

import "fmt"

const DebugGC = true

func debugPrintln(message string) {
	fmt.Println(message)
}

func debugPrintlf(message string, args ...any) {
	fmt.Printf(message, args...)
	fmt.Println()
}

func debugStressGC() {
	CollectGarbage()
}

func DebugStressGC() {
	debugStressGC()
}
