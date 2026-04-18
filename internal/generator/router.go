package generator

import (
	"fmt"
	"os"
	"strings"
)

func InjectRoute(moduleName string) error {
	filePath := "router/router.go"

	// ✅ kalau router belum ada → buat dulu
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		template := `package router

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

}
`
		if err := os.MkdirAll("router", os.ModePerm); err != nil {
			return err
		}
		if err := os.WriteFile(filePath, []byte(template), os.ModePerm); err != nil {
			return err
		}
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	content := string(data)
	lower := strings.ToLower(moduleName)

	modulePath := GetModulePath()
	if modulePath == "" {
		return fmt.Errorf("module path tidak ditemukan")
	}

	importLine := fmt.Sprintf(`"%s/modules/%s/route"`, modulePath, lower)
	registerLine := fmt.Sprintf("\troute.Register%sRoutes(r)", capitalize(moduleName))

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
	if !strings.Contains(content, registerLine) {
		content = strings.Replace(
			content,
			"func RegisterRoutes(r *gin.Engine) {",
			"func RegisterRoutes(r *gin.Engine) {\n"+registerLine,
			1,
		)
	}

	if err := os.WriteFile(filePath, []byte(content), os.ModePerm); err != nil {
		return err
	}

	return nil
}
