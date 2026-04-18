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

func TestGetModulePathWithoutGoMod(t *testing.T) {
	tmp := withTempDir(t)

	if got := GetModulePath(); got != "" {
		t.Fatalf("expected empty module path when go.mod is missing in %q, got %q", tmp, got)
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

func TestEnsureGoModReturnsErrorWhenGoBinaryMissing(t *testing.T) {
	tmp := withTempDir(t)
	t.Setenv("PATH", "")

	if _, err := EnsureGoMod("fallback"); err == nil {
		t.Fatalf("expected EnsureGoMod to return error in %q", tmp)
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

func TestRunGoModTidyReturnsErrorWhenGoBinaryMissing(t *testing.T) {
	tmp := withTempDir(t)
	t.Setenv("PATH", "")

	if err := RunGoModTidy(); err == nil {
		t.Fatalf("expected RunGoModTidy to return error in %q", tmp)
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

func TestGenerateFullModuleReturnsErrorAtModelStage(t *testing.T) {
	assertGenerateFullModuleStageError(t, "model")
}

func TestGenerateFullModuleReturnsErrorAtRepositoryStage(t *testing.T) {
	assertGenerateFullModuleStageError(t, "repository")
}

func TestGenerateFullModuleReturnsErrorAtServiceStage(t *testing.T) {
	assertGenerateFullModuleStageError(t, "service")
}

func TestGenerateFullModuleReturnsErrorAtHandlerStage(t *testing.T) {
	assertGenerateFullModuleStageError(t, "handler")
}

func TestGenerateFullModuleReturnsErrorAtRouteStage(t *testing.T) {
	assertGenerateFullModuleStageError(t, "route")
}

func assertGenerateFullModuleStageError(t *testing.T, stage string) {
	t.Helper()

	tmp := t.TempDir()
	base := filepath.Join(tmp, "modules", "user")
	conflict := filepath.Join(base, stage)
	if err := os.MkdirAll(filepath.Dir(conflict), 0o755); err != nil {
		t.Fatalf("failed to create parent directory: %v", err)
	}
	if err := os.WriteFile(conflict, []byte("file"), 0o644); err != nil {
		t.Fatalf("failed to create conflict file: %v", err)
	}

	if err := GenerateFullModule(base, "user", "example.com/demo"); err == nil {
		t.Fatalf("expected GenerateFullModule to return error for stage %s", stage)
	}
}

func TestWriteReturnsErrorOnMkdirAllFailure(t *testing.T) {
	tmp := t.TempDir()
	conflict := filepath.Join(tmp, "conflict")
	if err := os.WriteFile(conflict, []byte("file"), 0o644); err != nil {
		t.Fatalf("failed to create conflict file: %v", err)
	}

	if err := write(filepath.Join(conflict, "child.go"), "content"); err == nil {
		t.Fatalf("expected write to return error when mkdir fails")
	}
}

func TestWriteReturnsErrorOnWriteFileFailure(t *testing.T) {
	tmp := t.TempDir()
	base := filepath.Join(tmp, "base")
	if err := os.MkdirAll(base, 0o755); err != nil {
		t.Fatalf("failed to create base dir: %v", err)
	}
	filePath := filepath.Join(base, "child.go")
	if err := os.WriteFile(filePath, []byte("content"), 0o444); err != nil {
		t.Fatalf("failed to create read-only file: %v", err)
	}

	if err := write(filePath, "new content"); err == nil {
		t.Fatalf("expected write to return error when write fails")
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

func TestInjectRouteCreatesRouterFileWhenMissing(t *testing.T) {
	tmp := withTempDir(t)
	if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module example.com/demo\n\ngo 1.25.4\n"), 0o644); err != nil {
		t.Fatalf("failed to write go.mod: %v", err)
	}

	if err := InjectRoute("user"); err != nil {
		t.Fatalf("InjectRoute returned error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmp, "router", "router.go"))
	if err != nil {
		t.Fatalf("expected router.go to exist: %v", err)
	}
	if !strings.Contains(string(data), "route.RegisterUserRoutes(r)") {
		t.Fatalf("expected generated router to contain registration call")
	}
}

func TestInjectRouteReturnsErrorWithoutGoMod(t *testing.T) {
	tmp := withTempDir(t)
	if err := InjectRoute("user"); err == nil {
		t.Fatalf("expected InjectRoute to return error in %q", tmp)
	}
}

func TestInjectRouteReturnsErrorWhenRouterPathIsDirectory(t *testing.T) {
	tmp := withTempDir(t)
	if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module example.com/demo\n\ngo 1.25.4\n"), 0o644); err != nil {
		t.Fatalf("failed to write go.mod: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(tmp, "router", "router.go"), 0o755); err != nil {
		t.Fatalf("failed to create router dir conflict: %v", err)
	}

	if err := InjectRoute("user"); err == nil {
		t.Fatalf("expected InjectRoute to return error when router path is a directory")
	}
}

func TestInjectRouteReturnsErrorWhenMkdirFails(t *testing.T) {
	tmp := withTempDir(t)
	if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module example.com/demo\n\ngo 1.25.4\n"), 0o644); err != nil {
		t.Fatalf("failed to write go.mod: %v", err)
	}
	if err := os.Chmod(tmp, 0o555); err != nil {
		t.Fatalf("failed to make temp dir read-only: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chmod(tmp, 0o755)
	})

	if err := InjectRoute("user"); err == nil {
		t.Fatalf("expected InjectRoute to return error when mkdir fails")
	}
}

func TestInjectRouteReturnsErrorWhenTemplateWriteFails(t *testing.T) {
	tmp := withTempDir(t)
	if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module example.com/demo\n\ngo 1.25.4\n"), 0o644); err != nil {
		t.Fatalf("failed to write go.mod: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(tmp, "router"), 0o555); err != nil {
		t.Fatalf("failed to create router dir: %v", err)
	}

	if err := InjectRoute("user"); err == nil {
		t.Fatalf("expected InjectRoute to return error when template write fails")
	}
}

func TestInjectRouteKeepsExistingGeneratedContent(t *testing.T) {
	tmp := withTempDir(t)
	if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module example.com/demo\n\ngo 1.25.4\n"), 0o644); err != nil {
		t.Fatalf("failed to write go.mod: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(tmp, "router"), 0o755); err != nil {
		t.Fatalf("failed to create router dir: %v", err)
	}
	content := `package router

import (
	"github.com/gin-gonic/gin"
	"example.com/demo/modules/user/route"
)

func RegisterRoutes(r *gin.Engine) {
	route.RegisterUserRoutes(r)
}
`
	if err := os.WriteFile(filepath.Join(tmp, "router", "router.go"), []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write router template: %v", err)
	}

	if err := InjectRoute("user"); err != nil {
		t.Fatalf("InjectRoute returned error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmp, "router", "router.go"))
	if err != nil {
		t.Fatalf("failed to read router.go: %v", err)
	}
	if !strings.Contains(string(data), "example.com/demo/modules/user/route") {
		t.Fatalf("expected router import to remain present")
	}
}

func TestInjectRouteReturnsErrorWhenWriteFails(t *testing.T) {
	tmp := withTempDir(t)
	if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module example.com/demo\n\ngo 1.25.4\n"), 0o644); err != nil {
		t.Fatalf("failed to write go.mod: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(tmp, "router"), 0o755); err != nil {
		t.Fatalf("failed to create router dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "router", "router.go"), []byte(`package router

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

}
`), 0o444); err != nil {
		t.Fatalf("failed to write router file: %v", err)
	}

	if err := InjectRoute("user"); err == nil {
		t.Fatalf("expected InjectRoute to return error when write fails")
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

func TestInjectMainUpdatesExistingMainFile(t *testing.T) {
	tmp := withTempDir(t)
	if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module example.com/demo\n\ngo 1.25.4\n"), 0o644); err != nil {
		t.Fatalf("failed to write go.mod: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "main.go"), []byte(`package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
}
`), 0o644); err != nil {
		t.Fatalf("failed to write main.go: %v", err)
	}

	if err := InjectMain(); err != nil {
		t.Fatalf("InjectMain returned error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmp, "main.go"))
	if err != nil {
		t.Fatalf("failed to read main.go: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "example.com/demo/router") {
		t.Fatalf("expected injected router import to be present")
	}
}

func TestInjectMainReturnsErrorWithoutGoMod(t *testing.T) {
	tmp := withTempDir(t)
	if err := InjectMain(); err == nil {
		t.Fatalf("expected InjectMain to return error in %q", tmp)
	}
}

func TestInjectMainReturnsErrorWhenMainPathIsDirectory(t *testing.T) {
	tmp := withTempDir(t)
	if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module example.com/demo\n\ngo 1.25.4\n"), 0o644); err != nil {
		t.Fatalf("failed to write go.mod: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(tmp, "main.go"), 0o755); err != nil {
		t.Fatalf("failed to create main.go dir conflict: %v", err)
	}

	if err := InjectMain(); err == nil {
		t.Fatalf("expected InjectMain to return error when main.go is a directory")
	}
}

func TestInjectMainReturnsErrorWhenWriteFails(t *testing.T) {
	tmp := withTempDir(t)
	if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module example.com/demo\n\ngo 1.25.4\n"), 0o644); err != nil {
		t.Fatalf("failed to write go.mod: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "main.go"), []byte(`package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
}
`), 0o444); err != nil {
		t.Fatalf("failed to write main.go: %v", err)
	}

	if err := InjectMain(); err == nil {
		t.Fatalf("expected InjectMain to return error when write fails")
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