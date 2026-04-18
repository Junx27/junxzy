package commands

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestHelloCommandName(t *testing.T) {
	cmd := HelloCommand{}

	if got := cmd.Name(); got != "hello" {
		t.Fatalf("expected command name %q, got %q", "hello", got)
	}
}

func TestHelloCommandExecuteWithoutArgsPrintsUsage(t *testing.T) {
	cmd := HelloCommand{}

	output := captureHelloOutput(t, func() {
		cmd.Execute(nil)
	})

	if !strings.Contains(output, "Usage: hello <name>") {
		t.Fatalf("expected usage message, got %q", output)
	}
}

func TestHelloCommandExecuteWithNamePrintsGreeting(t *testing.T) {
	cmd := HelloCommand{}

	output := captureHelloOutput(t, func() {
		cmd.Execute([]string{"Junx"})
	})

	if !strings.Contains(output, "Hello, Junx!") {
		t.Fatalf("expected greeting message, got %q", output)
	}
}

func captureHelloOutput(t *testing.T, fn func()) string {
	t.Helper()

	originalStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	defer func() {
		os.Stdout = originalStdout
	}()

	os.Stdout = w
	fn()
	_ = w.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	return buf.String()
}
