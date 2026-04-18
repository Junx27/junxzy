package commands

import (
	"bytes"
	"os"
	"testing"
)

func TestClearCommandName(t *testing.T) {
	cmd := ClearCommand{}

	if got := cmd.Name(); got != "clear" {
		t.Fatalf("expected command name %q, got %q", "clear", got)
	}
}

func TestClearCommandExecutePrintsClearSequence(t *testing.T) {
	cmd := ClearCommand{}

	originalStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	defer func() {
		os.Stdout = originalStdout
	}()

	os.Stdout = w
	cmd.Execute(nil)
	_ = w.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("failed to read command output: %v", err)
	}

	if got := buf.String(); got != "\033[H\033[2J" {
		t.Fatalf("unexpected clear output: %q", got)
	}
}