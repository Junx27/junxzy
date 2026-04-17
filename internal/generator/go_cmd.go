package generator

import (
	"os/exec"
)

func RunGoModTidy() error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}
