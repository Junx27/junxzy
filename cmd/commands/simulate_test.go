package commands

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"
)

func TestSimulateCommandName(t *testing.T) {
	cmd := SimulateCommand{}

	if got := cmd.Name(); got != "simulate" {
		t.Fatalf("expected command name %q, got %q", "simulate", got)
	}
}

func TestSimulateCommandExecutePrintsSteps(t *testing.T) {
	originalSleep := simulateSleep
	simulateSleep = func(duration time.Duration) {}
	t.Cleanup(func() {
		simulateSleep = originalSleep
	})

	cmd := SimulateCommand{}

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
	for _, expected := range []string{
		"Creating project folder",
		"Generate services",
		"Setup gateway",
		"Generate docker-compose",
		"Finalizing project",
		"Simulation completed",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected output to contain %q, got %q", expected, output)
		}
	}
}