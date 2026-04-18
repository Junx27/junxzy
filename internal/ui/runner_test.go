package ui

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/fatih/color"
)

func captureStdout(t *testing.T, fn func()) string {
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
		t.Fatalf("failed to read stdout: %v", err)
	}

	return buf.String()
}

func TestRunStepSuccess(t *testing.T) {
	originalNoColor := color.NoColor
	color.NoColor = true
	t.Cleanup(func() {
		color.NoColor = originalNoColor
		s = nil
	})

	output := captureStdout(t, func() {
		RunStep("mengerjakan tugas", func() {})
	})

	if !strings.Contains(output, "mengerjakan tugas selesai") {
		t.Fatalf("expected success output to contain completion message, got %q", output)
	}
	if strings.Contains(output, "Error:") {
		t.Fatalf("did not expect error output on success path, got %q", output)
	}
}

func TestRunStepRecoversFromPanic(t *testing.T) {
	originalNoColor := color.NoColor
	color.NoColor = true
	t.Cleanup(func() {
		color.NoColor = originalNoColor
		s = nil
	})

	output := captureStdout(t, func() {
		RunStep("mengerjakan tugas", func() {
			panic("boom")
		})
	})

	if !strings.Contains(output, "Error: boom") {
		t.Fatalf("expected panic output to contain recovered error, got %q", output)
	}
}
