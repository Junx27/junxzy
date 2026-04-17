package commands

import (
	"fmt"
	"os"
	"path/filepath"
)

type MakeModuleCommand struct{}

func (m MakeModuleCommand) Name() string {
	return "make:module"
}

func (m MakeModuleCommand) Execute(args []string) {
	if len(args) < 1 {
		fmt.Println("Nama module wajib diisi")
		return
	}

	name := args[0]
	base := filepath.Join("modules", name)

	os.MkdirAll(filepath.Join(base, "handler"), os.ModePerm)
	os.MkdirAll(filepath.Join(base, "service"), os.ModePerm)
	os.MkdirAll(filepath.Join(base, "repository"), os.ModePerm)

	fmt.Println("Module berhasil dibuat:", name)
}
