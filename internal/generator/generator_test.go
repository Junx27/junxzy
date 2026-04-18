package generator

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func withTempDir(t *testing.T) string {
	t.Helper()

	tmp := t.TempDir()
	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd failed: %v", err)
	}

	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("Chdir failed: %v", err)
	}

	t.Cleanup(func() {
		_ = os.Chdir(original)
	})

	return tmp
}

func TestGetModulePath(t *testing.T) {
	tmp := withTempDir(t)

	content := []byte("module example.com/demo\n\ngo 1.25.4\n")
	if err := os.WriteFile(filepath.Join(tmp, "go.mod"), content, 0o644); err != nil {
		t.Fatalf("failed to write go.mod: %v", err)
	}

	if got := GetModulePath(); got != "example.com/demo" {
		t.Fatalf("expected module path %q, got %q", "example.com/demo", got)
	}
}

func TestEnsureGoModReturnsExistingModulePath(t *testing.T) {
	tmp := withTempDir(t)

	if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module example.com/existing\n\ngo 1.25.4\n"), 0o644); err != nil {
		t.Fatalf("failed to write go.mod: %v", err)
	}

	got, err := EnsureGoMod("fallback")
	if err != nil {
		t.Fatalf("EnsureGoMod returned error: %v", err)
	}
	if got != "example.com/existing" {
		t.Fatalf("expected existing module path, got %q", got)
	}
}

func TestEnsureGoModInitializesMissingModule(t *testing.T) {
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go binary not available")
	}

	tmp := withTempDir(t)

	got, err := EnsureGoMod("sample-app")
	if err != nil {
		t.Fatalf("EnsureGoMod returned error: %v", err)
	}
	if got != "sample-app" {
		t.Fatalf("expected new module path %q, got %q", "sample-app", got)
	}

	data, err := os.ReadFile(filepath.Join(tmp, "go.mod"))
	if err != nil {
		t.Fatalf("expected go.mod to exist: %v", err)
	}
	if !strings.Contains(string(data), "module sample-app") {
		t.Fatalf("expected go.mod to contain module declaration")
	}
}

func TestGenerateFullModuleWritesTemplateFiles(t *testing.T) {
	tmp := t.TempDir()
	base := filepath.Join(tmp, "modules", "user")

	if err := GenerateFullModule(base, "user", "example.com/demo"); err != nil {
		t.Fatalf("GenerateFullModule returned error: %v", err)
	}

	checks := map[string][]string{
		filepath.Join(base, "model", "user.go"): {
			"package model",
			"type User struct",
		},
		filepath.Join(base, "repository", "user_repository.go"): {
			"package repository",
			"example.com/demo/modules/user/model",
		},
		filepath.Join(base, "service", "user_service.go"): {
			"package service",
			"example.com/demo/modules/user/repository",
		},
		filepath.Join(base, "handler", "user_handler.go"): {
			"package handler",
			"type UserHandler struct",
		},
		filepath.Join(base, "route", "user_route.go"): {
			"package route",
			"RegisterUserRoutes",
		},
	}

	for path, expectedSnippets := range checks {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("expected file %q to exist: %v", path, err)
		}

		content := string(data)
		for _, snippet := range expectedSnippets {
			if !strings.Contains(content, snippet) {
				t.Fatalf("expected %q to contain %q", path, snippet)
			}
		}
	}
}

func TestInjectRouteUpdatesRouterFile(t *testing.T) {
	tmp := withTempDir(t)

	if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module example.com/demo\n\ngo 1.25.4\n"), 0o644); err != nil {
		t.Fatalf("failed to write go.mod: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(tmp, "router"), 0o755); err != nil {
		t.Fatalf("failed to create router dir: %v", err)
	}

	template := `package router

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

}
`
	if err := os.WriteFile(filepath.Join(tmp, "router", "router.go"), []byte(template), 0o644); err != nil {
		t.Fatalf("failed to write router template: %v", err)
	}

	if err := InjectRoute("user"); err != nil {
		t.Fatalf("InjectRoute returned error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmp, "router", "router.go"))
	if err != nil {
		t.Fatalf("failed to read router.go: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, `"example.com/demo/modules/user/route"`) {
		t.Fatalf("expected injected import to be present")
	}
	if !strings.Contains(content, "route.RegisterUserRoutes(r)") {
		t.Fatalf("expected injected register call to be present")
	}
}

func TestInjectMainCreatesMainFile(t *testing.T) {
	tmp := withTempDir(t)

	if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module example.com/demo\n\ngo 1.25.4\n"), 0o644); err != nil {
		t.Fatalf("failed to write go.mod: %v", err)
	}

	if err := InjectMain(); err != nil {
		t.Fatalf("InjectMain returned error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmp, "main.go"))
	if err != nil {
		t.Fatalf("expected main.go to exist: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, `"example.com/demo/router"`) {
		t.Fatalf("expected router import to be present")
	}
	if !strings.Contains(content, "router.RegisterRoutes(r)") {
		t.Fatalf("expected router registration to be present")
	}
}

func TestRunGoModTidySucceedsInEmptyModule(t *testing.T) {
	tmp := withTempDir(t)

	if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module example.com/demo\n\ngo 1.25.4\n"), 0o644); err != nil {
		t.Fatalf("failed to write go.mod: %v", err)
	}

	if err := RunGoModTidy(); err != nil {
		t.Fatalf("RunGoModTidy returned error: %v", err)
	}
}