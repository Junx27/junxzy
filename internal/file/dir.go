package file

import (
	"os"
	"path/filepath"
)

func CreateDirs(base string, dirs []string) error {
	for _, d := range dirs {
		path := filepath.Join(base, d)
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}
