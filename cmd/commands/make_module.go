package commands

import (
	"path/filepath"
	"time"

	"os"

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

	ui.RunStep("Generate full module "+name, func() {
		time.Sleep(5 * time.Second)
		modulePath := generator.GetModulePath()

		err := generator.GenerateFullModule(base, name, modulePath)
		if err != nil {
			panic(err)
		}
	})

	ui.RunStep("Register route "+name, func() {
		time.Sleep(5 * time.Second)
		err := generator.InjectRoute(name)
		if err != nil {
			panic(err)
		}
	})

	ui.Success("Module berhasil dibuat: " + name)
}
