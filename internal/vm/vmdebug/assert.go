//go:build debug

package vmdebug

import "fmt"

const (
	DebugAssert = true
)

type panicAssert struct{}

var _ Assert = (*panicAssert)(nil)

// Assertf implements AssertCond.
func (s *panicAssert) Assertf(condition bool, message string, args ...any) {
	if !condition {
		panic(fmt.Errorf(message, args...))
	}
}

func init() {
	gAssert = &panicAssert{}
}
