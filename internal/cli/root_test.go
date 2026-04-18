package cli

import (
	"bytes"
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