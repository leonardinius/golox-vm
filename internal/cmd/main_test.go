package cmd

import "testing"

func TestXxx(t *testing.T) {
	rc := Main("/Users/leo/src/golox-vm/testdata/benchmark/fib.lox")
	if rc != 0 {
		t.Errorf("Main() returned %d, expected 0", rc)
	}
}
