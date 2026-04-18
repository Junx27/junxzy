package commands

import (
	"path/filepath"
	"time"

	"os"

	"github.com/Junx27/junxzy/internal/file"
	"github.com/Junx27/junxzy/internal/generator"
	"github.com/Junx27/junxzy/internal/ui"
)

var makeModuleSleep = time.Sleep
var ensureGoModFn = generator.EnsureGoMod
var createDirsFn = file.CreateDirs
var generateFullModuleFn = generator.GenerateFullModule
var injectRouteFn = generator.InjectRoute
var injectMainFn = generator.InjectMain
var runGoModTidyFn = generator.RunGoModTidy

type MakeModuleCommand struct{}

func (m MakeModuleCommand) Name() string {
	return "make:module"
}

func (m MakeModuleCommand) Execute(args []string) {
	if len(args) < 1 {
		ui.StopError("Module name is required")
		return
	}

	name := args[0]
	base := filepath.Join("modules", name)

	if _, err := os.Stat(base); !os.IsNotExist(err) {
		ui.StopError("Module already exists: " + name)
		return
	}

	modulePath, err := ensureGoModFn(name)
	if err != nil {
		ui.StopError(err.Error())
		return
	}

	dirs := []string{"model", "repository", "service", "handler", "route"}

	ui.RunStep("Creating module structure "+name, func() {
		makeModuleSleep(3 * time.Second)
		if err := createDirsFn(base, dirs); err != nil {
			panic(err)
		}
	})

	ui.RunStep("Generate full module "+name, func() {
		makeModuleSleep(3 * time.Second)
		err := generateFullModuleFn(base, name, modulePath)
		if err != nil {
			panic(err)
		}
	})

	ui.RunStep("Registering route "+name, func() {
		makeModuleSleep(3 * time.Second)
		err := injectRouteFn(name)
		if err != nil {
			panic(err)
		}
	})

	ui.RunStep("Injecting main.go", func() {
		makeModuleSleep(3 * time.Second)
		err := injectMainFn()
		if err != nil {
			panic(err)
		}
	})

	ui.RunStep("Running go mod tidy", func() {
		makeModuleSleep(3 * time.Second)
		err := runGoModTidyFn()
		if err != nil {
			panic(err)
		}
	})

	ui.Success("Module created successfully: " + name)
}
