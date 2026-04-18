package commands

import (
	"path/filepath"
	"time"

	"os"

	"github.com/Junx27/junxzy/internal/file"
	"github.com/Junx27/junxzy/internal/generator"
	"github.com/Junx27/junxzy/internal/ui"
)

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

	modulePath, err := generator.EnsureGoMod(name)
	if err != nil {
		ui.StopError(err.Error())
		return
	}

	dirs := []string{"model", "repository", "service", "handler", "route"}

	ui.RunStep("Creating module structure "+name, func() {
		time.Sleep(3 * time.Second)
		if err := file.CreateDirs(base, dirs); err != nil {
			panic(err)
		}
	})

	ui.RunStep("Generate full module "+name, func() {
		time.Sleep(3 * time.Second)
		err := generator.GenerateFullModule(base, name, modulePath)
		if err != nil {
			panic(err)
		}
	})

	ui.RunStep("Registering route "+name, func() {
		time.Sleep(3 * time.Second)
		err := generator.InjectRoute(name)
		if err != nil {
			panic(err)
		}
	})

	ui.RunStep("Injecting main.go", func() {
		time.Sleep(3 * time.Second)
		err := generator.InjectMain()
		if err != nil {
			panic(err)
		}
	})

	ui.RunStep("Running go mod tidy", func() {
		time.Sleep(3 * time.Second)
		err := generator.RunGoModTidy()
		if err != nil {
			panic(err)
		}
	})

	ui.Success("Module created successfully: " + name)
}
