package cli

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/fatih/color"
)

func TestRootCommandConfiguration(t *testing.T) {
	if !rootCmd.CompletionOptions.DisableDefaultCmd {
		t.Fatalf("expected default completion command to be disabled")
	}

	if rootCmd.Use != "junxzy" {
		t.Fatalf("expected root command use to be %q, got %q", "junxzy", rootCmd.Use)
	}

	if rootCmd.Short != "Happy coding with Junxzy CLI! 🚀" {
		t.Fatalf("unexpected root command short description: %q", rootCmd.Short)
	}
}

func TestPrintBannerWritesExpectedText(t *testing.T) {
	originalNoColor := color.NoColor
	color.NoColor = true
	t.Cleanup(func() {
		color.NoColor = originalNoColor
	})

	originalStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	defer func() {
		os.Stdout = originalStdout
	}()

	os.Stdout = w
	printBanner()
	_ = w.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("failed to read banner output: %v", err)
	}

	output := buf.String()
	for _, expected := range []string{
		"Welcome to Junxzy CLI!",
		"Happy coding with Golang",
		"Type 'help' to see all commands.",
		"Type 'exit' to exit.",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected banner output to contain %q, got %q", expected, output)
		}
	}
}

func TestExecuteRunsRootCommand(t *testing.T) {
	originalRunRootCommand := runRootCommand
	originalExitProcess := exitProcess
	t.Cleanup(func() {
		runRootCommand = originalRunRootCommand
		exitProcess = originalExitProcess
	})

	called := false
	runRootCommand = func() error {
		called = true
		return nil
	}
	exitProcess = func(code int) {
		t.Fatalf("did not expect exit with code %d", code)
	}

	Execute()

	if !called {
		t.Fatalf("expected Execute to call root command")
	}
}

func TestExecuteExitsOnError(t *testing.T) {
	originalRunRootCommand := runRootCommand
	originalExitProcess := exitProcess
	t.Cleanup(func() {
		runRootCommand = originalRunRootCommand
		exitProcess = originalExitProcess
	})

	runRootCommand = func() error {
		return errors.New("boom")
	}
	exitCalled := false
	exitProcess = func(code int) {
		exitCalled = true
		if code != 1 {
			t.Fatalf("expected exit code 1, got %d", code)
		}
	}

	originalStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	defer func() {
		os.Stderr = originalStderr
	}()

	os.Stderr = w
	Execute()
	_ = w.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("failed to read stdout: %v", err)
	}

	if !exitCalled {
		t.Fatalf("expected exit to be called")
	}
	if !strings.Contains(buf.String(), "boom") {
		t.Fatalf("expected error output to contain boom, got %q", buf.String())
	}
}

func TestExecuteRunsRealRootCommand(t *testing.T) {
	originalStdin := os.Stdin
	originalStdout := os.Stdout
	originalNoColor := color.NoColor
	defer func() {
		os.Stdin = originalStdin
		os.Stdout = originalStdout
		color.NoColor = originalNoColor
		rootCmd.SetArgs([]string{})
	}()

	stdinReader, stdinWriter, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdin pipe: %v", err)
	}
	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdout pipe: %v", err)
	}

	color.NoColor = true
	rootCmd.SetArgs([]string{})
	os.Stdin = stdinReader
	os.Stdout = stdoutWriter
	_ = stdinWriter.Close()

	Execute()
	_ = stdoutWriter.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(stdoutReader); err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	if !strings.Contains(buf.String(), "Welcome to Junxzy CLI!") {
		t.Fatalf("expected banner output, got %q", buf.String())
	}
}

func TestStartREPLSkipsBlankInput(t *testing.T) {
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

	if _, err := stdinWriter.WriteString("\nhelp\n"); err != nil {
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
	if strings.Count(output, "junxzy") < 2 {
		t.Fatalf("expected REPL prompt to appear twice, got %q", output)
	}
	if !strings.Contains(output, "List command") {
		t.Fatalf("expected help output in REPL result, got %q", output)
	}
}
