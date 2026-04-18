package commands

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestExitCommandName(t *testing.T) {
	cmd := ExitCommand{}

	if got := cmd.Name(); got != "exit" {
		t.Fatalf("expected command name %q, got %q", "exit", got)
	}
}

func TestExitCommandExecuteInSubprocess(t *testing.T) {
	if os.Getenv("BE_EXIT_COMMAND") != "1" {
		return
	}

	cmd := ExitCommand{}
	cmd.Execute(nil)
}

func TestExitCommandExecutePrintsAndExitsZero(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=TestExitCommandExecuteInSubprocess")
	cmd.Env = append(os.Environ(), "BE_EXIT_COMMAND=1")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("subprocess failed: %v, output=%s", err, string(output))
	}

	if !strings.Contains(string(output), "Bye!") {
		t.Fatalf("expected subprocess output to contain farewell message, got %q", string(output))
	}
}