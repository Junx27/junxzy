package file

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateDirs(t *testing.T) {
	tmp := t.TempDir()

	err := CreateDirs(tmp, []string{"alpha", filepath.Join("beta", "gamma")})
	if err != nil {
		t.Fatalf("CreateDirs returned error: %v", err)
	}

	for _, path := range []string{
		filepath.Join(tmp, "alpha"),
		filepath.Join(tmp, "beta", "gamma"),
	} {
		if info, err := os.Stat(path); err != nil || !info.IsDir() {
			t.Fatalf("expected directory %q to exist, err=%v", path, err)
		}
	}
}