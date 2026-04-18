package commands

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateDirCreatesTargetFolder(t *testing.T) {
	tmp := t.TempDir()
	target := filepath.Join(tmp, "project")

	if err := createDir(target); err != nil {
		t.Fatalf("createDir returned error: %v", err)
	}

	if info, err := os.Stat(target); err != nil || !info.IsDir() {
		t.Fatalf("expected directory %q to exist, err=%v", target, err)
	}
}

func TestCreateServiceCreatesStructureAndMainFile(t *testing.T) {
	tmp := t.TempDir()

	if err := createService(tmp, "user"); err != nil {
		t.Fatalf("createService returned error: %v", err)
	}

	for _, dir := range []string{"handler", "service", "repository"} {
		path := filepath.Join(tmp, "services", "user", dir)
		if info, err := os.Stat(path); err != nil || !info.IsDir() {
			t.Fatalf("expected directory %q to exist, err=%v", path, err)
		}
	}

	mainFile := filepath.Join(tmp, "services", "user", "main.go")
	data, err := os.ReadFile(mainFile)
	if err != nil {
		t.Fatalf("expected main.go to exist: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "package main") {
		t.Fatalf("expected main.go to contain package declaration")
	}
	if !strings.Contains(content, "user service running") {
		t.Fatalf("expected main.go to contain service message")
	}
}

func TestCreateDockerComposeAndReadme(t *testing.T) {
	tmp := t.TempDir()

	if err := createDockerCompose(tmp); err != nil {
		t.Fatalf("createDockerCompose returned error: %v", err)
	}
	if err := createReadme(tmp, "demo"); err != nil {
		t.Fatalf("createReadme returned error: %v", err)
	}

	dockerComposePath := filepath.Join(tmp, "docker-compose.yml")
	dockerCompose, err := os.ReadFile(dockerComposePath)
	if err != nil {
		t.Fatalf("expected docker-compose.yml to exist: %v", err)
	}
	if !strings.Contains(string(dockerCompose), "user:") || !strings.Contains(string(dockerCompose), "8001:8000") {
		t.Fatalf("unexpected docker-compose.yml content")
	}

	readmePath := filepath.Join(tmp, "README.md")
	readme, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("expected README.md to exist: %v", err)
	}
	if !strings.Contains(string(readme), "# demo") {
		t.Fatalf("unexpected README.md content")
	}
}

func TestInitCommandName(t *testing.T) {
	cmd := InitCommand{}

	if got := cmd.Name(); got != "init" {
		t.Fatalf("expected command name %q, got %q", "init", got)
	}
}

func TestInitCommandExecuteWithoutArgsPrintsError(t *testing.T) {
	cmd := InitCommand{}

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
		t.Fatalf("failed to read output: %v", err)
	}

	if !strings.Contains(buf.String(), "Nama project wajib diisi") {
		t.Fatalf("expected missing project message, got %q", buf.String())
	}
}

func TestInitCommandExecuteWhenProjectAlreadyExists(t *testing.T) {
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

	if err := os.MkdirAll("demo", 0o755); err != nil {
		t.Fatalf("failed to create project dir: %v", err)
	}

	cmd := InitCommand{}
	output := captureInitOutput(t, func() {
		cmd.Execute([]string{"demo"})
	})

	if !strings.Contains(output, "Project sudah ada: demo") {
		t.Fatalf("expected existing project message, got %q", output)
	}
}

func TestInitCommandExecuteSuccess(t *testing.T) {
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

	cmd := InitCommand{}
	cmd.Execute([]string{"demo"})

	for _, path := range []string{
		filepath.Join("demo", "services", "user", "main.go"),
		filepath.Join("demo", "services", "auth", "main.go"),
		filepath.Join("demo", "services", "product", "main.go"),
		filepath.Join("demo", "gateway"),
		filepath.Join("demo", "docker-compose.yml"),
		filepath.Join("demo", "README.md"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected %q to exist: %v", path, err)
		}
	}
}

func TestInitCommandExecuteCreateDirError(t *testing.T) {
	assertInitStepError(t, func() {
		initCreateDirFn = func(path string) error {
			return errors.New("create dir failed")
		}
	}, "create dir failed")
}

func TestInitCommandExecuteCreateServiceError(t *testing.T) {
	assertInitStepError(t, func() {
		initCreateServiceFn = func(base, name string) error {
			return errors.New("create service failed")
		}
	}, "create service failed")
}

func TestInitCommandExecuteGatewayError(t *testing.T) {
	assertInitStepError(t, func() {
		initCreateGatewayFn = func(base string) error {
			return errors.New("gateway failed")
		}
	}, "gateway failed")
}

func TestInitCommandExecuteDockerComposeError(t *testing.T) {
	assertInitStepError(t, func() {
		initCreateDockerComposeFn = func(base string) error {
			return errors.New("compose failed")
		}
	}, "compose failed")
}

func TestInitCommandExecuteReadmeError(t *testing.T) {
	assertInitStepError(t, func() {
		initCreateReadmeFn = func(base, name string) error {
			return errors.New("readme failed")
		}
	}, "readme failed")
}

func TestCreateDirReturnsError(t *testing.T) {
	tmp := t.TempDir()
	conflict := filepath.Join(tmp, "conflict")
	if err := os.WriteFile(conflict, []byte("file"), 0o644); err != nil {
		t.Fatalf("failed to create conflict file: %v", err)
	}

	if err := createDir(filepath.Join(conflict, "child")); err == nil {
		t.Fatalf("expected createDir to return error")
	}
}

func TestCreateServiceReturnsError(t *testing.T) {
	tmp := t.TempDir()
	baseServices := filepath.Join(tmp, "services")
	if err := os.WriteFile(baseServices, []byte("file"), 0o644); err != nil {
		t.Fatalf("failed to create conflict file: %v", err)
	}

	if err := createService(tmp, "user"); err == nil {
		t.Fatalf("expected createService to return error")
	}
}

func TestCreateDockerComposeReturnsError(t *testing.T) {
	tmp := t.TempDir()
	conflict := filepath.Join(tmp, "docker-compose.yml")
	if err := os.MkdirAll(conflict, 0o755); err != nil {
		t.Fatalf("failed to create conflict directory: %v", err)
	}

	if err := createDockerCompose(tmp); err == nil {
		t.Fatalf("expected createDockerCompose to return error")
	}
}

func TestCreateReadmeReturnsError(t *testing.T) {
	tmp := t.TempDir()
	conflict := filepath.Join(tmp, "README.md")
	if err := os.MkdirAll(conflict, 0o755); err != nil {
		t.Fatalf("failed to create conflict directory: %v", err)
	}

	if err := createReadme(tmp, "demo"); err == nil {
		t.Fatalf("expected createReadme to return error")
	}
}

func captureInitOutput(t *testing.T, fn func()) string {
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

func assertInitStepError(t *testing.T, configure func(), expected string) {
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

	originalCreateDir := initCreateDirFn
	originalCreateService := initCreateServiceFn
	originalCreateGateway := initCreateGatewayFn
	originalCreateDockerCompose := initCreateDockerComposeFn
	originalCreateReadme := initCreateReadmeFn
	t.Cleanup(func() {
		initCreateDirFn = originalCreateDir
		initCreateServiceFn = originalCreateService
		initCreateGatewayFn = originalCreateGateway
		initCreateDockerComposeFn = originalCreateDockerCompose
		initCreateReadmeFn = originalCreateReadme
	})

	initCreateDirFn = func(path string) error { return nil }
	initCreateServiceFn = func(base, name string) error { return nil }
	initCreateGatewayFn = func(base string) error { return nil }
	initCreateDockerComposeFn = func(base string) error { return nil }
	initCreateReadmeFn = func(base, name string) error { return nil }
	configure()

	cmd := InitCommand{}
	output := captureInitOutput(t, func() {
		cmd.Execute([]string{"demo"})
	})

	if !strings.Contains(output, expected) {
		t.Fatalf("expected error output, got %q", output)
	}
}