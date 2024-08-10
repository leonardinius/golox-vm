//go:build debug

package vmdebug

import "fmt"

type panicAssert struct{}

var _ asserts = (*panicAssert)(nil)

// Assertf implements AssertCond.
func (s *panicAssert) Assertf(condition bool, message string, args ...any) {
	if !condition {
		panic(fmt.Errorf(message, args...))
	}
}
