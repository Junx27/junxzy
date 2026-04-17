package generator

import (
	"fmt"
	"os"
	"strings"
)

func InjectMain() error {
	filePath := "main.go"

	modulePath := GetModulePath()
	if modulePath == "" {
		return fmt.Errorf("module path tidak ditemukan")
	}

	// ✅ kalau belum ada → buat dari nol
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		template := `package main

import (
	"github.com/gin-gonic/gin"
	"` + modulePath + `/router"
)

func main() {
	r := gin.Default()

	router.RegisterRoutes(r)

	r.Run(":8080")
}
`
		return os.WriteFile(filePath, []byte(template), os.ModePerm)
	}

	// ✅ kalau sudah ada → inject
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	content := string(data)

	importLine := `"` + modulePath + `/router"`

	// inject import
	if !strings.Contains(content, importLine) {
		content = strings.Replace(
			content,
			"import (",
			"import (\n\t"+importLine,
			1,
		)
	}

	// inject register
	if !strings.Contains(content, "router.RegisterRoutes") {
		content = strings.Replace(
			content,
			"r := gin.Default()",
			"r := gin.Default()\n\n\trouter.RegisterRoutes(r)",
			1,
		)
	}

	return os.WriteFile(filePath, []byte(content), os.ModePerm)
}
