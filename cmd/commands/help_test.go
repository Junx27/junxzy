package commands

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestHelpCommandName(t *testing.T) {
	cmd := HelpCommand{}

	if got := cmd.Name(); got != "help" {
		t.Fatalf("expected command name %q, got %q", "help", got)
	}
}

func TestHelpCommandExecutePrintsCommands(t *testing.T) {
	cmd := HelpCommand{}

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

	output := buf.String()
	for _, expected := range []string{"Daftar command", "-- help", "-- clear", "-- exit", "-- make:module [name]"} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected output to contain %q, got %q", expected, output)
		}
	}
}