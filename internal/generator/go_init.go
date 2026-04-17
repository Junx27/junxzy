package generator

import (
	"fmt"
	"os"
	"os/exec"
)

func EnsureGoMod(projectName string) (string, error) {
	// cek apakah go.mod sudah ada
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {

		// fallback module name (simple)
		moduleName := projectName

		// jalankan: go mod init
		cmd := exec.Command("go", "mod", "init", moduleName)
		err := cmd.Run()
		if err != nil {
			return "", fmt.Errorf("gagal init go.mod: %v", err)
		}

		return moduleName, nil
	}

	// kalau sudah ada → ambil dari go.mod
	return GetModulePath(), nil
}
