//go:build debug

package vmmem

import "fmt"

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
