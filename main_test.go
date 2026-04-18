package main

import "testing"

func TestMainCallsExecute(t *testing.T) {
	originalExecute := execute
	called := false
	execute = func() {
		called = true
	}
	t.Cleanup(func() {
		execute = originalExecute
	})

	main()

	if !called {
		t.Fatalf("expected main to call execute")
	}
}
