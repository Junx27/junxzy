package commands

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
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

func TestMakeModuleCommandExecuteSuccess(t *testing.T) {
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

	originalSleep := makeModuleSleep
	originalEnsureGoMod := ensureGoModFn
	originalCreateDirs := createDirsFn
	originalGenerateFullModule := generateFullModuleFn
	originalInjectRoute := injectRouteFn
	originalInjectMain := injectMainFn
	originalRunGoModTidy := runGoModTidyFn
	t.Cleanup(func() {
		makeModuleSleep = originalSleep
		ensureGoModFn = originalEnsureGoMod
		createDirsFn = originalCreateDirs
		generateFullModuleFn = originalGenerateFullModule
		injectRouteFn = originalInjectRoute
		injectMainFn = originalInjectMain
		runGoModTidyFn = originalRunGoModTidy
	})

	makeModuleSleep = func(duration time.Duration) {}
	called := map[string]bool{}
	ensureGoModFn = func(projectName string) (string, error) {
		called["ensureGoMod"] = true
		return "example.com/demo", nil
	}
	createDirsFn = func(base string, dirs []string) error {
		called["createDirs"] = true
		return nil
	}
	generateFullModuleFn = func(base, moduleName, modulePath string) error {
		called["generateFullModule"] = true
		return nil
	}
	injectRouteFn = func(moduleName string) error {
		called["injectRoute"] = true
		return nil
	}
	injectMainFn = func() error {
		called["injectMain"] = true
		return nil
	}
	runGoModTidyFn = func() error {
		called["runGoModTidy"] = true
		return nil
	}

	cmd := MakeModuleCommand{}
	cmd.Execute([]string{"demo"})

	for _, name := range []string{"ensureGoMod", "createDirs", "generateFullModule", "injectRoute", "injectMain", "runGoModTidy"} {
		if !called[name] {
			t.Fatalf("expected %s to be called", name)
		}
	}
}

func TestMakeModuleCommandExecuteEnsureGoModError(t *testing.T) {
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

	originalEnsureGoMod := ensureGoModFn
	t.Cleanup(func() {
		ensureGoModFn = originalEnsureGoMod
	})

	ensureGoModFn = func(projectName string) (string, error) {
		return "", errors.New("init go.mod failed")
	}

	cmd := MakeModuleCommand{}
	output := captureMakeModuleOutput(t, func() {
		cmd.Execute([]string{"demo"})
	})

	if !strings.Contains(output, "init go.mod failed") {
		t.Fatalf("expected EnsureGoMod error output, got %q", output)
	}
}

func TestMakeModuleCommandExecuteCreateDirsError(t *testing.T) {
	assertMakeModuleStepError(t, func() {
		createDirsFn = func(base string, dirs []string) error {
			return errors.New("create dirs failed")
		}
	})
}

func TestMakeModuleCommandExecuteGenerateFullModuleError(t *testing.T) {
	assertMakeModuleStepError(t, func() {
		generateFullModuleFn = func(base, moduleName, modulePath string) error {
			return errors.New("generate failed")
		}
	})
}

func TestMakeModuleCommandExecuteInjectRouteError(t *testing.T) {
	assertMakeModuleStepError(t, func() {
		injectRouteFn = func(moduleName string) error {
			return errors.New("route failed")
		}
	})
}

func TestMakeModuleCommandExecuteInjectMainError(t *testing.T) {
	assertMakeModuleStepError(t, func() {
		injectMainFn = func() error {
			return errors.New("main failed")
		}
	})
}

func TestMakeModuleCommandExecuteRunGoModTidyError(t *testing.T) {
	assertMakeModuleStepError(t, func() {
		runGoModTidyFn = func() error {
			return errors.New("tidy failed")
		}
	})
}

func assertMakeModuleStepError(t *testing.T, configure func()) {
	t.Helper()

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

	originalSleep := makeModuleSleep
	originalEnsureGoMod := ensureGoModFn
	originalCreateDirs := createDirsFn
	originalGenerateFullModule := generateFullModuleFn
	originalInjectRoute := injectRouteFn
	originalInjectMain := injectMainFn
	originalRunGoModTidy := runGoModTidyFn
	t.Cleanup(func() {
		makeModuleSleep = originalSleep
		ensureGoModFn = originalEnsureGoMod
		createDirsFn = originalCreateDirs
		generateFullModuleFn = originalGenerateFullModule
		injectRouteFn = originalInjectRoute
		injectMainFn = originalInjectMain
		runGoModTidyFn = originalRunGoModTidy
	})

	makeModuleSleep = func(duration time.Duration) {}
	ensureGoModFn = func(projectName string) (string, error) {
		return "example.com/demo", nil
	}
	createDirsFn = func(base string, dirs []string) error { return nil }
	generateFullModuleFn = func(base, moduleName, modulePath string) error { return nil }
	injectRouteFn = func(moduleName string) error { return nil }
	injectMainFn = func() error { return nil }
	runGoModTidyFn = func() error { return nil }
	configure()

	cmd := MakeModuleCommand{}
	cmd.Execute([]string{"demo"})
}

func captureMakeModuleOutput(t *testing.T, fn func()) string {
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