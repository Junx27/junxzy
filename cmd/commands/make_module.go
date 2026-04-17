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
		ui.StopError("Nama module wajib diisi")
		return
	}

	name := args[0]
	base := filepath.Join("modules", name)

	if _, err := os.Stat(base); !os.IsNotExist(err) {
		ui.StopError("Module sudah ada: " + name)
		return
	}

	modulePath, err := generator.EnsureGoMod(name)
	if err != nil {
		ui.StopError(err.Error())
		return
	}

	dirs := []string{"model", "repository", "service", "handler", "route"}

	ui.RunStep("Membuat struktur module "+name, func() {
		time.Sleep(3 * time.Second)
		file.CreateDirs(base, dirs)
	})

	ui.RunStep("Generate full module "+name, func() {
		time.Sleep(3 * time.Second)
		err := generator.GenerateFullModule(base, name, modulePath)
		if err != nil {
			panic(err)
		}
	})

	ui.RunStep("Register route "+name, func() {
		time.Sleep(3 * time.Second)
		err := generator.InjectRoute(name)
		if err != nil {
			panic(err)
		}
	})

	ui.RunStep("Inject main.go", func() {
		time.Sleep(3 * time.Second)
		err := generator.InjectMain()
		if err != nil {
			panic(err)
		}
	})

	ui.RunStep("Run go mod tidy", func() {
		time.Sleep(3 * time.Second)
		err := generator.RunGoModTidy()
		if err != nil {
			panic(err)
		}
	})

	ui.Success("Module berhasil dibuat: " + name)
}
