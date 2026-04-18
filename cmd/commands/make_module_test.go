package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMakeModuleCommandName(t *testing.T) {
	cmd := MakeModuleCommand{}

	if got := cmd.Name(); got != "make:module" {
		t.Fatalf("expected command name %q, got %q", "make:module", got)
	}
}

func TestMakeModuleCommandExecuteWithoutArgsPrintsError(t *testing.T) {
	cmd := MakeModuleCommand{}

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
	if !strings.Contains(output, "Module name is required") {
		t.Fatalf("expected output to contain missing module message, got %q", output)
	}
}

func TestMakeModuleCommandExecuteWhenModuleAlreadyExists(t *testing.T) {
	tmp := t.TempDir()
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd failed: %v", err)
	}

	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("Chdir failed: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalWd)
	})

	if err := os.MkdirAll(filepath.Join("modules", "demo"), 0o755); err != nil {
		t.Fatalf("failed to create existing module dir: %v", err)
	}

	cmd := MakeModuleCommand{}

	originalStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	defer func() {
		os.Stdout = originalStdout
	}()

	os.Stdout = w
	cmd.Execute([]string{"demo"})
	_ = w.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("failed to read command output: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Module already exists: demo") {
		t.Fatalf("expected output to contain existing module message, got %q", output)
	}
}