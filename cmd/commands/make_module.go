package commands

import (
	"path/filepath"

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

	dirs := []string{"handler", "service", "repository"}

	ui.RunStep("Membuat struktur module "+name, func() {
		file.CreateDirs(base, dirs)
	})

	ui.RunStep("Generate CRUD "+name, func() {
		err := generator.CreateCRUD(base, name)
		if err != nil {
			ui.StopError("Gagal generate CRUD: " + err.Error())
			return
		}
	})

	ui.Success("Module berhasil dibuat: " + name)
}
