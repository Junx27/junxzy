package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestStartREPLExecutesKnownCommand(t *testing.T) {
	originalStdin := os.Stdin
	originalStdout := os.Stdout
	originalRegistry := commandRegistry
	defer func() {
		os.Stdin = originalStdin
		os.Stdout = originalStdout
		commandRegistry = originalRegistry
	}()

	commandRegistry = map[string]Command{}
	registerCommands()

	stdinReader, stdinWriter, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdin pipe: %v", err)
	}
	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdout pipe: %v", err)
	}

	os.Stdin = stdinReader
	os.Stdout = stdoutWriter

	if _, err := stdinWriter.WriteString("help\n"); err != nil {
		t.Fatalf("failed to write input: %v", err)
	}
	_ = stdinWriter.Close()

	startREPL()
	_ = stdoutWriter.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(stdoutReader); err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "List command") {
		t.Fatalf("expected help output in REPL result, got %q", output)
	}
	if !strings.Contains(output, "junxzy") {
		t.Fatalf("expected REPL prompt in output, got %q", output)
	}
}

func TestStartREPLPrintsUnknownCommandMessage(t *testing.T) {
	originalStdin := os.Stdin
	originalStdout := os.Stdout
	originalRegistry := commandRegistry
	defer func() {
		os.Stdin = originalStdin
		os.Stdout = originalStdout
		commandRegistry = originalRegistry
	}()

	commandRegistry = map[string]Command{}
	registerCommands()

	stdinReader, stdinWriter, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdin pipe: %v", err)
	}
	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdout pipe: %v", err)
	}

	os.Stdin = stdinReader
	os.Stdout = stdoutWriter

	if _, err := stdinWriter.WriteString("does-not-exist\n"); err != nil {
		t.Fatalf("failed to write input: %v", err)
	}
	_ = stdinWriter.Close()

	startREPL()
	_ = stdoutWriter.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(stdoutReader); err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Command tidak dikenal") {
		t.Fatalf("expected unknown command message in REPL output, got %q", output)
	}
	if !strings.Contains(output, "does-not-exist") {
		t.Fatalf("expected unknown command name in output, got %q", output)
	}
}