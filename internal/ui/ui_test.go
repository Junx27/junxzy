package ui

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/fatih/color"
)

func captureStdoutUI(t *testing.T, fn func()) string {
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

func TestSuccessPrintsMessage(t *testing.T) {
	originalNoColor := color.NoColor
	color.NoColor = true
	t.Cleanup(func() {
		color.NoColor = originalNoColor
	})

	output := captureStdoutUI(t, func() {
		Success("berhasil")
	})

	if !strings.Contains(output, "berhasil") {
		t.Fatalf("expected success output to contain message, got %q", output)
	}
}

func TestStartInitializesSpinner(t *testing.T) {
	Start("memuat data")
	t.Cleanup(func() {
		s = nil
	})

	if s == nil {
		t.Fatalf("expected spinner to be initialized")
	}
}

func TestStopSuccessPrintsMessage(t *testing.T) {
	originalNoColor := color.NoColor
	color.NoColor = true
	t.Cleanup(func() {
		color.NoColor = originalNoColor
		s = nil
	})

	s = nil
	output := captureStdoutUI(t, func() {
		StopSuccess("selesai")
	})

	if !strings.Contains(output, "✔ selesai") {
		t.Fatalf("expected success stop output to contain completion mark, got %q", output)
	}
}

func TestStopErrorPrintsMessage(t *testing.T) {
	originalNoColor := color.NoColor
	color.NoColor = true
	t.Cleanup(func() {
		color.NoColor = originalNoColor
		s = nil
	})

	s = nil
	output := captureStdoutUI(t, func() {
		StopError("gagal")
	})

	if !strings.Contains(output, "✖ gagal") {
		t.Fatalf("expected error stop output to contain error mark, got %q", output)
	}
}