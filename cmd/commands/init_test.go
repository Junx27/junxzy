package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateDirCreatesTargetFolder(t *testing.T) {
	tmp := t.TempDir()
	target := filepath.Join(tmp, "project")

	createDir(target)

	if info, err := os.Stat(target); err != nil || !info.IsDir() {
		t.Fatalf("expected directory %q to exist, err=%v", target, err)
	}
}

func TestCreateServiceCreatesStructureAndMainFile(t *testing.T) {
	tmp := t.TempDir()

	createService(tmp, "user")

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

	createDockerCompose(tmp)
	createReadme(tmp, "demo")

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